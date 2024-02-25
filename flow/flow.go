package flow

type FailHandler func(State, error) error

// Flow describes a process from beginning to end. It retains all defined steps, the user's final handler, and more.
// Additionally, it offers a straightforward interface to access internal storage for marshaling and saving elsewhere.
type Flow struct {
	// User's defined steps
	steps []Step
	// Calls after successfully passing full flow
	then StateHandler
	// Calls on any error (@TODO: update the comment)
	catch FailHandler

	// User options

	// Determines whether we need to send errors from a validator to the user as a response.
	// If true, errors from a validator are responded, otherwise, no response is sent.
	UseValidatorErrorsAsUserResponse bool
}
