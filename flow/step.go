package flow

// StepMetaData is an object that provides the user with information for different stages
type StepMetaData struct {
	// Ordinal number
	Step int
}

// StepThenHandler handler for the successfully completed step
type StepThenHandler func(State, *Step) error

// StepValidator ensures that each step can be validated before the flow process progresses to the next step.
// Typically, users require simple validators, such as ensuring text is not empty or a photo is uploaded.
// Therefore, having a variety of different validators and describing them in a single function is not ideal.
type StepValidator interface {
	// Validate is called after the user prompts anything.
	Validate(State) error
}

// Step describes a user's step within a flow.
type Step struct {
	// This is the user's custom function called at the beginning of a step.
	// There are no restrictions on logic; the handler is not required, and the user can even use an empty mock.
	// Therefore, you can do whatever you want: move backward or forward steps, validate previously saved prompts, and so on.
	handler StateHandler
	// Defined validators
	validators []StepValidator
	// Callback called after the validation process if successful.
	// It can, for example, assign the user's prompt to a variable.
	assign StateHandler
}
