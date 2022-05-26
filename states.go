package telebot

var States = StatesNative{}

type StatesNative struct {
	AutoStateStorage State
}

// 	auto fill State mashine
// 	if you are using Auto() -> do not use hand-writed (or iota) State
func (s *StatesNative) Auto() State {
	s.AutoStateStorage++
	return s.AutoStateStorage

}

// 	SetState - analog for Context.SetState()
// 	Use it for set State
func (StatesNative) SetState(c Context, state State) {
	userId := c.Sender().ID
	storage := c.Bot().statesStorage
	storage[userId] = state
}

// 	SetState - analog for Context.NextState()
// 	Use it for set current State + 1. Know, that your import numerate doesn't make sense (or use hand-writed State)
func (StatesNative) NextState(c Context) {
	userId := c.Sender().ID
	storage := c.Bot().statesStorage
	storage[userId]++
}

// 	SetState - analog for Context.NextState()
// 	Use it for set current State + 1. Know, that your import numerate doesn't make sense (or use hand-writed State)
func (StatesNative) FinishState(c Context) {

	userId := c.Sender().ID
	storage := c.Bot().statesStorage
	storage[userId] = 0
}

// 	return State(0) - state at the function not required
//	use as b.Handler(endpoint, h, States.Zero())
func (StatesNative) Zero() State {
	return State(0)
}
