package telebot

type states struct {
	autoStateStorage State
}

var States states

// 	auto fill State mashine
// 	if you are using Auto() -> do not use hand-writed (or iota) State
func (s states) Auto() State {
	result := s.autoStateStorage
	s.autoStateStorage++
	return result

}

// 	SetState - analog for Context.SetState()
// 	Use it for set State
func (states) SetState(c Context, state State) {
	userId := c.Sender().ID
	storage := c.Bot().statesStorage
	storage[userId] = state
}

// 	SetState - analog for Context.NextState()
// 	Use it for set current State + 1. Know, that your import numerate doesn't make sense (or use hand-writed State)
func (states) NextState(c Context) {
	userId := c.Sender().ID
	storage := c.Bot().statesStorage
	storage[userId]++
}

// 	SetState - analog for Context.NextState()
// 	Use it for set current State + 1. Know, that your import numerate doesn't make sense (or use hand-writed State)
func (states) FinishState(c Context) {

	userId := c.Sender().ID
	storage := c.Bot().statesStorage
	storage[userId] = 0
}

// 	return State(0) - state at the function not required
//	use as b.Handler(endpoint, h, States.Zero())
func (states) Zero() State {
	return State(0)
}

func init() {
	States = states{1}
}
