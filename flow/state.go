package flow

import (
	"gopkg.in/telebot.v3"
)

// StateHandler is a common handler for global processes, such as finally, fail, step, and so on.
type StateHandler func(state *State) error

// State defines the user's state and persists common elements, such as the bot instance, flow handler instance, etc.
type State struct {
	// Instance for the user's flow.
	Machine Machine
	// Received message from a user.
	Context telebot.Context

	// @TODO: We can provide a full history of every step,
	// including contexts, validation results, and so on.
	// contextHistory map[string][]StepHistory

	// @TODO: make it concurrently safe?
	// User state represents any custom data for a user.
	// It's simply a container that you can use within steps.
	// For instance, you may use it if you want to populate a struct at each step
	// and then use this data after the flow has successfully completed.
	userState map[interface{}]interface{}
}

// Set initializes the [userState] field.
func (s *State) Set(userState map[interface{}]interface{}) {
	s.userState = userState
}

// Get returns data corresponding to the provided key, along with a boolean value indicating if the key exists.
func (s *State) Get(key interface{}) (interface{}, bool) {
	value, exists := s.userState[key]

	return value, exists
}

// Put adds data to the storage.
// Caution: it does not check if the key already exists. It will overwrite any existing data associated with the key.
func (s *State) Put(key interface{}, value interface{}) {
	s.userState[key] = value
}

// Exists returns a boolean representing whether the key exists in the map.
func (s *State) Exists(key interface{}) bool {
	_, exists := s.userState[key]

	return exists
}

func NewState(machine Machine, c telebot.Context, userState map[interface{}]interface{}) *State {
	return &State{
		Machine:   machine,
		Context:   c,
		userState: userState,
	}
}
