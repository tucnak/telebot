package telebot

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Context = (*nativeContext)(nil)

func TestBotBeginFlow(t *testing.T) {
	b := &Bot{}
	h := func(c Context) error { return nil }
	flow := b.BeginFlow("final", h)

	assert.NotNil(t, flow)
	assert.Equal(t, "final", flow.current)
	assert.True(t, reflect.ValueOf(flow.begin).Pointer() == reflect.ValueOf(h).Pointer())
}

func TestFlowTransite(t *testing.T) {
	flow := &Flow{
		current:     "start",
		steps:       map[string]HandlerFunc{"start": nil, "step1": nil},
		transitions: make(map[string]map[string]TransitionFunc),
	}

	tests := []struct {
		name        string
		from        string
		to          string
		shouldPanic bool
		setup       func()
	}{
		{
			name:        "Panic if current step does not exist",
			from:        "unknown",
			to:          "step1",
			shouldPanic: true,
		},
		{
			name:        "Panic if next step does not exist",
			from:        "start",
			to:          "unknown",
			shouldPanic: true,
		},
		{
			name:        "Panic if trying to transite to the same step",
			from:        "start",
			to:          "start",
			shouldPanic: true,
		},
		{
			name:        "Successful transition registration",
			from:        "start",
			to:          "step1",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					flow.Transite(tt.from, tt.to, func(c Context) bool { return true })
				})
			} else {
				flow.Transite(tt.from, tt.to, func(c Context) bool { return true })
				assert.Contains(t, flow.transitions[tt.from], tt.to)
			}
		})
	}
}

func TestFlowForward(t *testing.T) {
	tests := []struct {
		name          string
		flow          *Flow
		expectAdvance bool
	}{
		{
			name: "Forward success",
			flow: &Flow{
				current: "start",
				transitions: map[string]map[string]TransitionFunc{
					"start": {
						"next": func(c Context) bool { return true },
					},
				},
			},
			expectAdvance: true,
		},
		{
			name: "No forward when transition condition is false",
			flow: &Flow{
				current: "start",
				transitions: map[string]map[string]TransitionFunc{
					"start": {
						"next": func(c Context) bool { return false },
					},
				},
			},
			expectAdvance: false,
		},
		{
			name: "No forward when no transitions exist",
			flow: &Flow{
				current:     "start",
				transitions: map[string]map[string]TransitionFunc{},
			},
			expectAdvance: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(nativeContext)
			result := tt.flow.Forward(c)
			assert.Equal(t, tt.expectAdvance, result)
		})
	}
}

func TestFlowIsLast(t *testing.T) {
	tests := []struct {
		name     string
		flow     *Flow
		expected bool
	}{
		{
			name: "IsLast true when no transitions",
			flow: &Flow{
				current:     "final",
				transitions: map[string]map[string]TransitionFunc{},
			},
			expected: true,
		},
		{
			name: "IsLast false when transitions exist",
			flow: &Flow{
				current: "step1",
				transitions: map[string]map[string]TransitionFunc{
					"step1": {
						"step2": func(c Context) bool { return true },
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.flow.IsLast()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlowHandle(t *testing.T) {
	t.Run("Basic handler without middleware", func(t *testing.T) {
		flow := &Flow{
			steps:      make(map[string]HandlerFunc),
			processors: make(map[string]map[string]HandlerFunc),
			current:    "step",
		}

		handlerCalled := false
		handler := func(c Context) error {
			handlerCalled = true
			return nil
		}

		flow.Handle("step", handler)
		flow.OnUpdate("update_type", "step", handler)

		stepHandler := flow.steps["step"]
		require.NotNil(t, stepHandler)
		_ = stepHandler(nil)
		assert.True(t, handlerCalled)

		handlerCalled = false
		updateHandler := flow.ProcessUpdate("update_type")
		require.NotNil(t, updateHandler)
		_ = updateHandler(nil)
		assert.True(t, handlerCalled)

		assert.Nil(t, flow.ProcessUpdate("unknown"))
	})

	t.Run("Handler with middleware", func(t *testing.T) {
		flow := &Flow{
			steps:       make(map[string]HandlerFunc),
			middlewares: []MiddlewareFunc{},
			current:     "step",
		}

		handlerCalled := false
		middlewareCalled := false

		handler := func(c Context) error {
			handlerCalled = true
			return nil
		}

		middleware := func(next HandlerFunc) HandlerFunc {
			return func(c Context) error {
				middlewareCalled = true
				return next(c)
			}
		}

		// Додаємо middleware до Flow
		flow.middlewares = append(flow.middlewares, middleware)

		flow.Handle("step", handler)

		stepHandler := flow.steps["step"]
		require.NotNil(t, stepHandler)

		_ = stepHandler(nil)

		assert.True(t, middlewareCalled, "middleware should be called")
		assert.True(t, handlerCalled, "handler should be called")
	})
}

func TestFlowContains(t *testing.T) {
	flow := &Flow{
		steps:      map[string]HandlerFunc{"a": nil},
		processors: map[string]map[string]HandlerFunc{"b": {"x": nil}},
	}

	assert.True(t, flow.Contains("a"))
	assert.True(t, flow.Contains("b"))
	assert.False(t, flow.Contains("c"))
}

func TestFlowManagerRegister(t *testing.T) {
	fm := &FlowManager{}
	flow := &Flow{current: "test"}

	fm.Register(flow)
	assert.True(t, fm.Contains("test"))

	flow2 := &Flow{}
	assert.Panics(t, func() { fm.RegisterAt("", flow2) })
	fm.RegisterAt("custom", flow2)
	assert.True(t, fm.Contains("custom"))
}

func TestCloneFlow(t *testing.T) {
	f := &Flow{current: "step"}
	clone := cloneFlow(f)
	assert.Equal(t, f.current, clone.current)
	assert.NotSame(t, f, clone)

	assert.Nil(t, cloneFlow(nil))
}
