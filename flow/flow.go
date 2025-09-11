package flow

import (
	tele "gopkg.in/telebot.v4"
	"sync"
	"time"
)

// Flow is a builder for defining a sequence of steps in a user interaction flow.
type Flow struct {
	name string

	steps []*Step

	okFn  StateHandler
	errFn ErrHandler

	timeout time.Duration
	start   time.Time

	mu sync.RWMutex
}

// OnText adds a step that triggers on text messages.
func (f *Flow) OnText(beginFn StateHandler) *Step {
	return f.addStep(tele.OnText, beginFn)
}

// OnCallback adds a step that triggers on callback queries.
func (f *Flow) OnCallback(beginFn StateHandler) *Step {
	return f.addStep(tele.OnCallback, beginFn)
}

// OnPhoto adds a step that triggers on photo messages.
func (f *Flow) OnPhoto(beginFn StateHandler) *Step {
	return f.addStep(tele.OnPhoto, beginFn)
}

// OnAudio adds a step that triggers on audio messages.
func (f *Flow) OnAudio(beginFn StateHandler) *Step {
	return f.addStep(tele.OnAudio, beginFn)
}

// OnDocument adds a step that triggers on document messages.
func (f *Flow) OnDocument(beginFn StateHandler) *Step {
	return f.addStep(tele.OnDocument, beginFn)
}

// OnVideo adds a step that triggers on video messages.
func (f *Flow) OnVideo(beginFn StateHandler) *Step {
	return f.addStep(tele.OnVideo, beginFn)
}

// OnVoice adds a step that triggers on voice messages.
func (f *Flow) OnVoice(beginFn StateHandler) *Step {
	return f.addStep(tele.OnVoice, beginFn)
}

// OnContact adds a step that triggers on contact messages.
func (f *Flow) OnContact(beginFn StateHandler) *Step {
	return f.addStep(tele.OnContact, beginFn)
}

// OnLocation adds a step that triggers on location messages.
func (f *Flow) OnLocation(beginFn StateHandler) *Step {
	return f.addStep(tele.OnLocation, beginFn)
}

func (f *Flow) addStep(on string, beginFn StateHandler) *Step {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// simple autoincrement.
	id := len(f.steps)

	step := &Step{
		id:      id,
		on:      on,
		beginFn: beginFn,
	}
	f.steps = append(f.steps, step)

	return step
}

// Ok sets the handler called when the entire flow completes successfully.
func (f *Flow) Ok(fn StateHandler) *Flow {
	f.okFn = fn
	return f
}

// Err sets the handler called when the flow fails.
func (f *Flow) Err(fn ErrHandler) *Flow {
	f.errFn = fn
	return f
}

// IsExpired checks if the flow has timed out since starting.
func (f *Flow) IsExpired() bool {
	if f.timeout == 0 {
		return false
	}
	return time.Since(f.start) > f.timeout
}
