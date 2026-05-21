package telebot

import (
	"sync"
	"time"
)

// OnState converts a state name into an internal endpoint key.
//
// Example:
//
//	bot.Handle(OnState("name"), NameHandler)
func OnState(state string) string {
	return "\r" + state
}

// StateStorage provides methods for managing user states
// and temporary state-related data.
type StateStorage interface {

	// SetState updates the user's current state.
	SetState(userID int64, state string)

	// GetState returns the user's current state.
	GetState(userID int64) (state string, ok bool)

	// SetData stores additional data associated with the user's state.
	SetData(userID int64, key string, value interface{})

	// GetData returns state-related data by key.
	GetData(userID int64, key string) interface{}

	// Delete removes the user's state and all associated data.
	Delete(userID int64)
}

// Cleaner is implemented by StateStorage backends
// that support cleanup of inactive user states.
type Cleaner interface {

	// CleanInactive removes user states that have not been
	// updated within the specified period.
	CleanInactive(period time.Duration, stop chan struct{})
}

// StartCleanup starts inactive state cleanup
// if the storage implements Cleaner.
func StartCleanup(storage StateStorage, period time.Duration, stop chan struct{}) {
	if cleaner, ok := storage.(Cleaner); ok {
		cleaner.CleanInactive(period, stop)
	}
}

// StateContext stores the current user state
// and its associated temporary data.
type stateContext struct {
	State      string
	LastUpdate time.Time
	Data       map[string]interface{}
}

func newStateContext(state string) *stateContext {
	return &stateContext{
		State:      state,
		LastUpdate: time.Now(),
		Data:       make(map[string]interface{}),
	}
}

// MemoryStorage stores user states in memory.
type MemoryStorage struct {
	mu    sync.RWMutex
	users map[int64]*stateContext // int64 - userID
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users: make(map[int64]*stateContext),
	}
}

func (s *MemoryStorage) SetState(userID int64, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.users[userID]
	if !ok {
		s.users[userID] = newStateContext(state)
		return
	}

	ctx.State = state
	ctx.LastUpdate = time.Now()
}

func (s *MemoryStorage) GetState(userID int64) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.users[userID]
	if !ok {
		return "", false
	}

	return ctx.State, true
}

func (s *MemoryStorage) SetData(userID int64, key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.users[userID]
	if !ok {
		return
	}

	ctx.Data[key] = value
}

func (s *MemoryStorage) GetData(userID int64, key string) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.users[userID]
	if !ok {
		return nil
	}

	return ctx.Data[key]
}

func (s *MemoryStorage) Delete(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.users, userID)
}

func (s *MemoryStorage) CleanInactive(period time.Duration, stop chan struct{}) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()

			for id, ctx := range s.users {
				if time.Since(ctx.LastUpdate) > period {
					delete(s.users, id)
				}
			}

			s.mu.Unlock()
		case <-stop:
			return
		}
	}
}
