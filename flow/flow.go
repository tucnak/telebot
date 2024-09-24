package flow

// MetaDataFailureStage provides information on the case in which the step failed: begin/validation/assign, and so on.
type MetaDataFailureStage uint8

const (
	// MetaDataFailureStageNone indicates that the step was successfully passed.
	MetaDataFailureStageNone MetaDataFailureStage = iota
	// MetaDataFailureStageBegin means fail happened on the first stage (maybe some network problems for instance)
	MetaDataFailureStageBegin
	// MetaDataFailureStageValidation indicates that the user input prompted bad data
	// or something occurred during the validation process.
	MetaDataFailureStageValidation
	// MetaDataFailureStageAssign means that the assigner function returned an error
	// (which could be due to an internal problem, for instance).
	MetaDataFailureStageAssign
	MetaDataFailureStageThen
)

// MetaData is an object that provides the user with information for different stages
type MetaData struct {
	// Endpoint for the Telegram bot that is served by the flow.
	Endpoint interface{}
	// Provides the last active step.
	// For example, this is useful when the flow was terminated due to inactivity.
	LastActiveStep StepMetaData
	FailureStage   MetaDataFailureStage
	// Sometimes the flow can fail due to a step.
	// If this occurs, the data indicates the failed step.
	FailedStep *StepMetaData
	// Provides the error that caused the failure, if there was one.
	FailedError error
}

// OnEachStepHandler called after the step is executed.
// Please refer to the documentation for [Flow.onEachStep] below for more details.
type OnEachStepHandler func(State, StepMetaData)

type FailHandler func(State, *MetaData) error

// Flow describes a process from beginning to end. It retains all defined steps, the user's final handler, and more.
// Additionally, it offers a straightforward interface to access internal storage for marshaling and saving elsewhere.
type Flow struct {
	// User's defined steps
	steps []Step
	// Calls after successfully passing full flow
	then StateHandler
	// Calls on any error (@TODO: update the comment)
	catch FailHandler
	// This handler is called for each step and only on success.
	onEachStep OnEachStepHandler
}
