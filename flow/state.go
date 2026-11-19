package flow

import (
	tele "gopkg.in/telebot.v4"
)

// StateHandler is a function that handles a state transition or action.
type StateHandler func(state *State) error

// ValidateHandler is a function for validation data from [tele.Context].
type ValidateHandler func(state *State) tele.HandlerFunc

// ErrHandler handles flow errors.
type ErrHandler func(*State, error) error

// State holds the current [Machine] and [tele.Context] for a flow step.
type State struct {
	Machine Machine
	Context tele.Context
}

// NewState constructor [State].
func NewState(machine Machine, c tele.Context) *State {
	return &State{
		Machine: machine,
		Context: c,
	}
}
