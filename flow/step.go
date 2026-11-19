package flow

// Step defines users step in [Flow].
type Step struct {
	// Step id (default autoincrement). (Required)
	id int

	// TODO: Maybe add `name` field for detail logging?

	// Basic message handlers [tele.OnText, tele.OnCallback, tele.OnAudio, etc]. (Required)
	// https://github.com/tucnak/telebot/blob/v4/telebot.go
	on string

	// Callback called at the beginning of the step.
	// It is mainly used to encourage the user to take action.
	//
	// Example:
	//		flow.OnText(func(state *flow.State) error {
	//			return state.Context.Send("Enter username:")
	//		})
	beginFn StateHandler

	// Callback called after [beginFn].
	// It is intended for validation only.
	//
	// Example:
	//	f.OnText(func(state *flow.State) error {
	//			return state.Context.EditOrSend("Enter username:")
	//		}).Validations(func(state *flow.State) tele.HandlerFunc {
	//			data := state.Context.Message().Text
	//			if len(data) > 10 {
	//				return func(c tele.Context) error {
	//					return c.EditOrSend("username max length 10, please check data and send me")
	//				}
	//			}
	//			return nil
	//		})
	validatorsFns []ValidateHandler

	// Callback called after validation.
	// It can, for example, assign the user's result to a variable.
	//
	// Example:
	// 		var username string
	//
	//		f.OnText(func(state *flow.State) error {
	//			return state.Context.EditOrSend("Enter username:")
	//		}).Validations(func(state *flow.State) tele.HandlerFunc {
	//          // ...
	//		}).Assign(func(state *flow.State) error {
	//			username = state.Context.Message().Text
	//			return nil
	//		})
	assignFn StateHandler

	// Callback called after [assignFn].
	// Example:
	// 		var username string
	//
	//		f.OnText(func(state *flow.State) error {
	//			return state.Context.EditOrSend("Enter username:")
	//		}).Validations(func(state *flow.State) tele.HandlerFunc {
	//			// ...
	//		}).Assign(func(state *flow.State) error {
	//			username = state.Context.Message().Text
	//			return nil
	//		}).Ok(func(state *flow.State) error {
	//			log.Println("username = ", username)
	//			return nil
	//		})
	okFn StateHandler
}

// Validations adds validation handlers to the step.
// Each returns a handler to execute on failure (e.g., send error message); nil on success.
func (s *Step) Validations(fns ...ValidateHandler) *Step {
	s.validatorsFns = append(s.validatorsFns, fns...)
	return s
}

// Assign sets the handler to assign or process data after validation.
func (s *Step) Assign(fn StateHandler) *Step {
	s.assignFn = fn
	return s
}

// Ok sets the handler called after successful assignment in this step.
func (s *Step) Ok(fn StateHandler) *Step {
	s.okFn = fn
	return s
}
