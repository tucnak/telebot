package flow

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/flow/storage"
	"strconv"
	"time"
)

var (
	NoStepsDefined = errors.New("flow: no steps defined")
	IsExpired      = errors.New("flow: expired")
)

// Bus defines the methods for managing user flows.
type Bus interface {
	// Flow creates a new flow builder. The timeout defines how long the entire flow can run before expiring.
	// If timeout is zero, the flow never expires.
	Flow(name string, timeout time.Duration) *Flow

	// Start initiates or resumes a flow for the user.
	// Only one flow can be active per user at a time; if an existing flow is active, it resumes from its current step.
	//
	// Example:
	//	bot.Handle("/signup", func(c tele.Context) error {
	//		f := bus.Flow("signup", flow.DefaultTimeout)
	//
	//		f.OnText(func(state *flow.State) error {
	//			// ..
	//		}).Assign(func(state *flow.State) error {
	//			// ..
	//		}).Ok(func(state *flow.State) error {
	//			// ..
	//		})
	//
	//		return bus.Start(ctx,c, f)
	//	})
	Start(ctx context.Context, c telebot.Context, flow *Flow) error

	// Handle processes incoming events and advances the flow state if applicable.
	// Register this as a handler for relevant telebot events, e.g.:
	//   bot.Handle(tele.OnText, bus.Handle)
	//   bot.Handle(tele.OnCallback, bus.Handle)
	Handle(telebot.Context) error
}

// state wraps the flow and its current state for storage.
type state struct {
	flow  *Flow
	state *State
}

// bus implementation for the [Bus] contract.
type bus struct {
	// Stores the active user flows.
	// userID -> state
	//
	// NOTE: Please check if your [storage.Storage] is secure to concurrency.
	states storage.Storage
}

func (b *bus) Flow(name string, timeout time.Duration) *Flow {
	return &Flow{
		name:    name,
		steps:   make([]*Step, 0),
		timeout: timeout,
	}
}

func (b *bus) Handle(c telebot.Context) error {
	uid := strconv.FormatInt(c.Sender().ID, 10)

	stValue, err := b.states.Get(uid)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil
		}

		return fmt.Errorf("flow: get stored value: %w", err)
	}

	if stValue == nil {
		return nil
	}

	st := stValue.(*state)
	if st.flow.IsExpired() {
		if err = b.states.Delete(uid); err != nil {
			return fmt.Errorf("flow: failed to remove expired state: %w", err)
		}

		return IsExpired
	}

	st.state.Context = c
	step := st.flow.steps[st.state.Machine.ActiveStep()]

	// checks if the event type matches the step's expected event.
	event := b.getEventType(c)
	if event != step.on || event == "" {
		return nil
	}

	// call [validatorFns] event if it's defined.
	for _, fn := range step.validatorsFns {
		handler := fn(st.state)
		if handler != nil {
			if err = handler(c); err != nil {
				return fmt.Errorf("flow: failed to execute validator func: %w", err)
			}
		}
	}

	// call [assignFn] event if it's defined.
	if step.assignFn != nil {
		if err = step.assignFn(st.state); err != nil {
			return fmt.Errorf("flow: failed to assign step: %w", err)
		}
	}

	// call [okFn] event if it's defined.
	if step.okFn != nil {
		if err := step.okFn(st.state); err != nil {
			return fmt.Errorf("flow: failed to ok step: %w", err)
		}
	}

	// it was the last step.
	if len(st.flow.steps) <= st.state.Machine.ActiveStep()+1 {
		if err = b.states.Delete(uid); err != nil {
			return fmt.Errorf("flow: failed to remove completed state: %w", err)
		}

		// always call the [ok] handler if the flow is finished.
		return st.state.Machine.Ok(st.state)
	}

	// process to the next step
	err = st.state.Machine.Next(st.state)
	if err != nil {
		// in case of any error, we delete the current status from the user.
		if delErr := b.states.Delete(uid); delErr != nil {
			return fmt.Errorf("flow: failed to remove state on error: %w (original err: %v)", delErr, err)
		}

		// always call the [error] if err !=.
		return st.state.Machine.Fail(st.state, err)
	}

	return nil
}

func (b *bus) Start(ctx context.Context, c telebot.Context, flow *Flow) error {
	if len(flow.steps) == 0 {
		return NoStepsDefined
	}

	uid := strconv.FormatInt(c.Sender().ID, 10)

	// if the user already has a flow, we need to recall the active step.
	stV, _ := b.states.Get(uid)
	if stV != nil {
		st := stV.(*state)
		st.state.Context = c
		if st.flow.IsExpired() {
			if err := b.states.Delete(uid); err != nil {
				return fmt.Errorf("flow: failed to remove expired state: %w", err)
			}
			return IsExpired
		}

		return st.state.Machine.ToStep(st.state.Machine.ActiveStep(), st.state)
	}

	machine := NewMachine(flow)

	// register flow for the user.
	st := state{
		flow:  flow,
		state: NewState(machine, c),
	}
	// set start time when actually starting.
	flow.start = time.Now()

	if err := b.states.Set(ctx, strconv.FormatInt(c.Sender().ID, 10), &st, flow.timeout); err != nil {
		return fmt.Errorf("flow: failed to set state: %w", err)
	}

	// call the machine for the start the first step.
	return machine.ToStep(0, st.state)
}

// getEventType determines the telebot event type from the context.
func (b *bus) getEventType(c telebot.Context) string {
	if c.Message().Text != "" {
		return telebot.OnText
	}
	if c.Callback() != nil {
		return telebot.OnCallback
	}
	if c.Message().Audio != nil {
		return telebot.OnAudio
	}
	if c.Message().Document != nil {
		return telebot.OnDocument
	}
	if c.Message().Video != nil {
		return telebot.OnVideo
	}
	if c.Message().Voice != nil {
		return telebot.OnVoice
	}
	if c.Message().Contact != nil {
		return telebot.OnContact
	}
	if c.Message().Location != nil {
		return telebot.OnLocation
	}

	return ""
}

// NewBus constructor [Bus].
//
// NOTE: Please check if your [storage.Storage] is secure to concurrency.
func NewBus(states storage.Storage) Bus {
	return &bus{
		states: states,
	}
}
