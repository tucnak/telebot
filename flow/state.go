package flow

// Defines keys for the basic data that must be in the state
const (
	StateMachineKey = "machine"
	StateContextKey = "context"
)

// StateHandler is a common handler for global processes, such as finally, fail, step, and so on.
type StateHandler func(State) error

// State defines the user's state and persists common elements, such as the bot instance, flow handler instance, etc.
type State interface {
	// Set initializes the [userState] field.
	Set(userState map[interface{}]interface{})
	// Get returns data corresponding to the provided key, along with a boolean value indicating if the key exists.
	Get(key interface{}) (interface{}, bool)
	// Read returns data corresponding to the provided key.
	// You can especially use this for fast type assertion.
	//
	//  state.Read("machine").(flow.Machine)
	Read(key interface{}) interface{}
	// Add adds data to the storage.
	// Caution: it does not check if the key already exists. It will overwrite any existing data associated with the key.
	Add(key interface{}, value interface{}) State
	// Exists returns a boolean representing whether the key exists in the map.
	Exists(key interface{}) bool
}

type RuntimeState struct {
	// @TODO: make it concurrently safe?
	// User state represents any custom data for a user.
	// It's simply a container that you can use within steps.
	// For instance, you may use it if you want to populate a struct at each step
	// and then use this data after the flow has successfully completed.
	userState map[interface{}]interface{}
}

func (s *RuntimeState) Set(userState map[interface{}]interface{}) {
	s.userState = userState
}

func (s *RuntimeState) Get(key interface{}) (interface{}, bool) {
	value, exists := s.userState[key]

	return value, exists
}

func (s *RuntimeState) Read(key interface{}) interface{} {
	return s.userState[key]
}

func (s *RuntimeState) Add(key interface{}, value interface{}) State {
	s.userState[key] = value

	return s
}

func (s *RuntimeState) Exists(key interface{}) bool {
	_, exists := s.userState[key]

	return exists
}

func NewRuntimeState(userState map[interface{}]interface{}) State {
	return &RuntimeState{userState: userState}
}
