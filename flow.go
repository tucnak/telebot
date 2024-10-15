package telebot

import (
	"errors"
	"fmt"
	"sync"
)

// FlowFunc represents a flow processing function,
// which compile flow like handler function.
type FlowFunc func(m ...MiddlewareFunc) *Flow

// Flow represents the flow of steps and transitions in a bot's conversation.
type Flow struct {
	*Bot

	current string // Current step in the flow.

	steps       map[string]interface{}               // Registered steps in the flow.
	processors  map[string]map[string]HandlerFunc    // Handlers for updates specific to a step.
	transitions map[string]map[string]TransitionFunc // Transition functions between steps.

	middlewares []MiddlewareFunc
}

// Begin initializes a new flow from a specified step,
// applying optional middleware.
//
// Transite function is necessary to move between handlers of the flow.
//
// Example:
//
//	b.
//	     Begin("lang_choose").
//	     Handle("lang_choose", b.OnLangChoose).
//	     Handle("lang_chosen", b.OnLangChosen).
//	     Transite("lang_choose", "lang_chosen", func(c tele.Context, u tele.Update) bool { return u.Callback != nil }).
//
// With OnUpdate Example:
//
//	b.
//	     Begin("lang_choose").
//	     Handle("lang_choose", b.OnLangChoose).
//	     Handle("lang_chosen", b.OnLangChosen).
//	     OnUpdate(tele.OnCallback, "lang_chosen", func(c tele.Context) error { return nil }).
//	     Transite("lang_choose", "lang_chosen", func(c tele.Context, u tele.Update) bool { return u.Callback != nil }).
func (b *Bot) Begin(step string, m ...MiddlewareFunc) *Flow {
	return &Flow{
		Bot:         b,
		steps:       make(map[string]interface{}),
		current:     step,
		middlewares: appendMiddleware(b.group.middleware, m),
	}
}

// TransitionFunc defines a function to determine whether a transition to the next step is possible.
type TransitionFunc func(c Context, u Update) bool

// Forward moves the flow to the next step if a transition function returns true.
func (f *Flow) Forward(c Context, u Update) {
	possibleTransitions := f.transitions[f.current]
	for nextStep, transition := range possibleTransitions {
		if transition(c, u) {
			f.current = nextStep
			return
		}
	}
}

// IsLast checks if the current step is the last one in the flow.
func (f *Flow) IsLast() bool {
	return f.transitions[f.current] == nil
}

// Enter executes the handler or flow for the current step, applying middleware if necessary.
func (f *Flow) Enter(m ...MiddlewareFunc) func(c Context) error {
	return func(c Context) error {
		step, exists := f.steps[f.current]
		if !exists {
			return fmt.Errorf("step %s not found in flow", f.current)
		}

		switch h := step.(type) {
		case HandlerFunc:
			if f.middlewares != nil && len(f.middlewares) > 0 {
				step = func(c Context) error {
					return applyMiddleware(h, f.middlewares...)(c)
				}
			}

			step = applyMiddleware(h, m...)
			return h(c)

		case *Flow:
			return h.Enter(m...)(c)
		default:
			return errors.New("telebot: unknown step type")
		}
	}
}

// OnUpdate registers a handler for a specific update at a specific step, with optional middleware.
func (f *Flow) OnUpdate(update string, step string, handler HandlerFunc, m ...MiddlewareFunc) *Flow {
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

	return f
}

// Handle registers a handler or a sub-flow for a specific step, with optional middleware.
func (f *Flow) Handle(step string, handler interface{}, m ...MiddlewareFunc) *Flow {
	if len(f.middlewares) > 0 {
		m = append(f.middlewares, m...)
	}

	switch h := handler.(type) {
	case func(c Context) error:
		f.steps[step] = applyMiddleware(h, m...)
	case HandlerFunc:
		f.steps[step] = applyMiddleware(h, m...)
	case func(m ...MiddlewareFunc) *Flow:
		newFlow := h(appendMiddleware(f.middlewares, m)...)

		f.steps[step] = newFlow
		f.flowManager.Register(step, newFlow)
	case FlowFunc:
		newFlow := h(appendMiddleware(f.middlewares, m)...)

		f.steps[step] = newFlow
		f.flowManager.Register(step, newFlow)
	default:
		panic("telebot: invalid handler")
	}

	return f
}

// Transite registers a transition function between two steps.
func (f *Flow) Transite(step, next string, t TransitionFunc) *Flow {
	if f.transitions == nil {
		f.transitions = make(map[string]map[string]TransitionFunc)
	}

	if f.transitions[step] == nil {
		f.transitions[step] = make(map[string]TransitionFunc)
	}

	f.transitions[step][next] = t
	return f
}

// FlowManager manages multiple flows and stores active ones.
// Helps in navigation among users' flows.
type FlowManager struct {
	flows map[string]*Flow
	store map[string]*Flow // Active user flows.

	middlewares []MiddlewareFunc

	mu sync.Mutex
}

// IsUsed checks if there are any active flows.
func (fm *FlowManager) IsUsed() bool {
	return len(fm.flows) > 0
}

// Exists checks if a flow exists by name.
func (fm *FlowManager) Exists(name string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	_, exists := fm.flows[name]
	return exists
}

// Get retrieves a flow by name.
func (fm *FlowManager) Get(name string) *Flow {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	return fm.flows[name]
}

// Register registers a new flow by its endpoint.
func (fm *FlowManager) Register(endpoint string, f *Flow) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.flows == nil {
		fm.flows = make(map[string]*Flow)
	}

	if f.current != endpoint {
		fm.flows[f.current] = f
	}

	fm.flows[endpoint] = f
}

// Close removes the flow associated with a user if it's completed.
func (fm *FlowManager) Close(u Recipient) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if flow, exists := fm.store[u.Recipient()]; exists && flow.IsLast() {
		delete(fm.store, u.Recipient())
		return true
	}

	return false
}

// MakeProcessing returns the handler function for the current step and update.
func (fm *FlowManager) MakeProcessing(c Context, u Update) HandlerFunc {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	user := c.Recipient()
	f, exists := fm.store[user.Recipient()]
	if !exists {
		return nil
	}

	if f.processors[f.current] == nil {
		return nil
	}

	if p, exists := f.processors[f.current][u.String()]; exists {
		return p
	}

	if u.Message != nil {
		if p, exists := f.processors[f.current][u.Message.String()]; exists {
			return p
		}
	}

	return nil
}

// MakeTransition advances the flow to the next step based on the update.
func (fm *FlowManager) MakeTransition(c Context, u Update) *Flow {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	user := c.Recipient()
	f, exists := fm.store[user.Recipient()]
	if !exists {
		return nil
	}

	f.Forward(c, u)
	return f
}

// Follow returns the flow associated with a user.
func (fm *FlowManager) Follow(user Recipient) *Flow {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	return fm.store[user.Recipient()]
}

// IsFollowed checks if a user is associated with an active flow.
func (fm *FlowManager) IsFollowed(user Recipient) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	_, exists := fm.store[user.Recipient()]
	return exists
}

// IsRegistred checks if a command is registered as a flow.
func (fm *FlowManager) IsRegistred(command string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	_, exists := fm.flows[command]
	return exists
}

// Start initializes and stores a flow for a user.
func (fm *FlowManager) Start(c Context, f *Flow) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	user := c.Recipient()
	if user == nil {
		return
	}

	if fm.store == nil {
		fm.store = make(map[string]*Flow)
	}

	fm.store[user.Recipient()] = cloneFlow(f)
}

// cloneFlow creates a shallow copy of the original flow.
func cloneFlow(original *Flow) *Flow {
	if original == nil {
		return nil
	}

	clone := *original
	return &clone
}
