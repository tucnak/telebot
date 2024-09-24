package flow

// This package contains factories for describing flows.
// Factories generate a flow object in a simple manner.

// Factory for creating a [Flow] object.
type Factory struct {
	flow *Flow
	// Represents any user state with [State.userState].
	userState map[interface{}]interface{}
}

// AddState adds a state to the [Factory.userState]
func (f *Factory) AddState(key interface{}, value interface{}) *Factory {
	f.userState[key] = value

	return f
}

// WithState sets a value for [Factory.userState]
func (f *Factory) WithState(userState map[interface{}]interface{}) *Factory {
	f.userState = userState

	return f
}

// Next adds a step to the [Flow.steps]
func (f *Factory) Next(step *StepFactory) *Factory {
	f.flow.steps = append(f.flow.steps, *step.step)

	return f
}

// OnEachStep sets a handler for the [Flow.onEachStep] event.
func (f *Factory) OnEachStep(handler OnEachStepHandler) *Factory {
	f.flow.onEachStep = handler

	return f
}

// Then sets a handler for the [Flow.then] event.
func (f *Factory) Then(handler StateHandler) *Factory {
	f.flow.then = handler

	return f
}

// Catch sets a handler for the [Flow.catch] event.
func (f *Factory) Catch(handler FailHandler) *Factory {
	f.flow.catch = handler

	return f
}

// New start describing the flow.
func New() *Factory {
	return &Factory{
		flow:      &Flow{},
		userState: make(map[interface{}]interface{}),
	}
}

// StepFactory for creating a [Step] object.
type StepFactory struct {
	step *Step
}

// Validate sets values for the [Step.validators].
func (f *StepFactory) Validate(validators ...StepValidator) *StepFactory {
	f.step.validators = validators

	return f
}

// Assign sets a value for the [Step.assign].
func (f *StepFactory) Assign(assign StateHandler) *StepFactory {
	f.step.assign = assign

	return f
}

// NewStep initiates the description of a step for the flow.
func NewStep(handler StateHandler) *StepFactory {
	return &StepFactory{step: &Step{handler: handler}}
}
