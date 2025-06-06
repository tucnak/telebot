package telebot

import (
	"fmt"
	"sync"
)

const FlowBeginProcessor = "begin"

// Begin initializes a new flow from a specified step,
// applying optional middleware.
//
// Transite function is necessary to move between handlers of the flow.
//
// Example:
//
//	b.
//	     Begin("lang_choose", b.OnLangChoose).
//	     Handle("lang_chosen", b.OnLangChosen).
//	     Transite("lang_choose", "lang_chosen", func(c tele.Context, u tele.Update) bool { return u.Callback != nil }).
//
// With Handle Example:
//
//	b.
//	     Begin("lang_choose", b.OnLangChoose).
//	     Handle("lang_chosen", b.OnLangChosen).
//	     Handle(tele.OnCallback, "lang_chosen", func(c tele.Context) error { return nil }).
//	     Transite("lang_choose", "lang_chosen", func(c tele.Context, u tele.Update) bool { return u.Callback != nil }).
func (b *Bot) BeginFlow(h HandlerFunc) *Flow {
	return &Flow{
		begin:       func(c Context) error { return applyMiddleware(h, b.group.middleware...)(c) },
		middlewares: b.group.middleware,
	}
}

func (b *Bot) advanceFlow(c Context, endpoint string) (flow *Flow, skip bool) {
	u := c.Recipient()
	if u == nil {
		return
	}

	if f, exists := b.flowManager.flows[endpoint]; exists {
		b.flowManager.mu.Lock()
		defer b.flowManager.mu.Unlock()

		if b.flowManager.store == nil {
			b.flowManager.store = make(map[string]*Flow)
		}

		b.flowManager.store[u.Recipient()] = cloneFlow(f)
		return b.flowManager.store[u.Recipient()], false
	}

	if f, exists := b.flowManager.store[u.Recipient()]; exists {
		return f, false
	}

	b.flowManager.Close(u)
	return
}

// hasActiveFlow checks if a user is associated with an active flow.
func (b *Bot) hasActiveFlow(user Recipient) bool {
	b.flowManager.mu.Lock()
	defer b.flowManager.mu.Unlock()

	_, exists := b.flowManager.store[user.Recipient()]
	return exists
}

const FlowStateName = "flow.state"

type FlowState int

const (
	FlowContinue FlowState = iota
	FlowRepeat
	FlowEnd
)

// Flow represents the flow of steps and transitions in a bot's conversation.
type Flow struct {
	begin   HandlerFunc // handler to begin the flow
	current string      // Current step in the flow

	processors  map[string]map[string]HandlerFunc    // Handlers for specific updates to a step.
	transitions map[string]map[string]TransitionFunc // Transition functions between steps.

	middlewares []MiddlewareFunc
}

// TransitionFunc defines a function to determine whether a transition to the next step is possible.
type TransitionFunc func(c Context) bool

var NoConditionTransition = func(c Context) bool { return true }

func (f *Flow) Contains(endpoint interface{}) bool {
	end := extractEndpoint(endpoint)
	if end == "" {
		return false
	}

	if _, exists := f.processors[end]; exists {
		return true
	}

	return false
}

// Forward moves the flow to the next step if a transition function returns true.
func (f *Flow) Forward(c Context) bool {
	possibleTransitions := f.transitions[f.current]
	for nextStep, transition := range possibleTransitions {
		if transition(c) {
			f.current = nextStep
			return true
		}
	}

	if f.current == "" {
		f.current = FlowBeginProcessor
		return true
	}

	return false
}

// IsLast checks if the current step is the last one in the flow.
func (f *Flow) IsLast() bool {
	return f.transitions[f.current] == nil
}

// Handle registers a handler for a specific update at a specific step, with optional middleware.
func (f *Flow) Handle(update string, step string, handler HandlerFunc, m ...MiddlewareFunc) {
	if len(f.middlewares) > 0 {
		m = append(f.middlewares, m...)
	}

	if f.processors == nil {
		f.processors = make(map[string]map[string]HandlerFunc)
	}

	if f.processors[step] == nil {
		f.processors[step] = make(map[string]HandlerFunc)
	}

	f.processors[step][update] = applyMiddleware(handler, m...)
}

func (f *Flow) ProcessUpdate(update string) HandlerFunc {
	if f.processors[f.current] == nil {
		return nil
	}

	return f.processors[f.current][update]
}

// Transite registers a transition function between two steps.
func (f *Flow) Transite(step, next string, t TransitionFunc) {
	if f.transitions == nil {
		f.transitions = make(map[string]map[string]TransitionFunc)
	}

	if f.transitions[step] == nil {
		f.transitions[step] = make(map[string]TransitionFunc)
	}

	if !f.Contains(step) {
		panic(fmt.Sprintf("step %s not found in registry", step))
	}

	if !f.Contains(next) {
		panic(fmt.Sprintf("step %s not found in registry", next))
	}

	if next == f.current {
		panic(fmt.Sprintf("transition cannot be continue from %s", f.current))
	}

	f.transitions[step][next] = t
}

// FirstTransite registers a transition function to process from begin.
func (f *Flow) FirstTransite(next string, t TransitionFunc) {
	if f.transitions == nil {
		f.transitions = make(map[string]map[string]TransitionFunc)
	}

	if f.transitions[FlowBeginProcessor] == nil {
		f.transitions[FlowBeginProcessor] = make(map[string]TransitionFunc)
	}

	if !f.Contains(next) {
		panic(fmt.Sprintf("step %s not found in registry", next))
	}

	if next == f.current {
		panic(fmt.Sprintf("transition cannot be continue from %s", f.current))
	}

	f.transitions[FlowBeginProcessor][next] = t
}

// FlowManager manages multiple flows and stores active ones.
// Helps in navigation among users' flows.
type FlowManager struct {
	flows map[string]*Flow
	store map[string]*Flow // Active user flows.

	middlewares []MiddlewareFunc

	mu sync.Mutex
}

// Contains checks a flow is existed by its name, e.g.
func (fm *FlowManager) Contains(name string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	_, exists := fm.flows[name]
	return exists
}

// RegisterAt registers a new flow by its endpoint.
func (fm *FlowManager) RegisterAt(endpoint string, f *Flow) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.flows == nil {
		fm.flows = make(map[string]*Flow)
	}

	if endpoint == "" {
		panic("telebot: empty endpoint")
	}

	fm.flows[endpoint] = f
}

// Register registers a new flow.
func (fm *FlowManager) Register(f *Flow) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.flows == nil {
		fm.flows = make(map[string]*Flow)
	}

	fm.flows[f.current] = f
}

// Close removes the flow associated with a user if it's completed.
func (fm *FlowManager) Close(r Recipient) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, exists := fm.store[r.Recipient()]; exists {
		delete(fm.store, r.Recipient())
		return true
	}

	return false
}

// cloneFlow creates a shallow copy of the original flow.
func cloneFlow(original *Flow) *Flow {
	if original == nil {
		return nil
	}

	clone := *original
	return &clone
}
