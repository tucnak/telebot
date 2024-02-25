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

// Success sets a handler for the [Flow.Success] event.
func (f *Factory) Success(handler StateHandler) *Factory {
	f.flow.success = handler

	return f
}

// Fail sets a handler for the [Flow.Fail] event.
func (f *Factory) Fail(handler FailHandler) *Factory {
	f.flow.fail = handler

	return f
}

// Step adds a step to the [Flow.Steps]
func (f *Factory) Step(step *StepFactory) *Factory {
	f.flow.steps = append(f.flow.steps, *step.step)

	return f
}

// UseValidatorErrorsAsUserResponse sets a value for the [Flow.useValidatorErrorsAsUserResponse].
func (f *Factory) UseValidatorErrorsAsUserResponse(value bool) *Factory {
	f.flow.useValidatorErrorsAsUserResponse = value

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

// Begin sets a handler for the [Step.begin] event.
func (f *StepFactory) Begin(handler StateHandler) *StepFactory {
	f.step.begin = handler

	return f
}

// Name sets a value for the [Step.name].
func (f *StepFactory) Name(name int) *StepFactory {
	f.step.name = name

	return f
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

// Success sets a value for the [Step.success].
func (f *StepFactory) Success(success StateHandler) *StepFactory {
	f.step.success = success

	return f
}

// NewStep initiates the description of a step for the flow.
func NewStep() *StepFactory {
	return &StepFactory{step: &Step{}}
}
