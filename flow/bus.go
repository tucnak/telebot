package flow

import (
	"errors"
	"gopkg.in/telebot.v3"
	"sync"
	"time"
)

var (
	UserDoesNotHaveActiveFlow = errors.New("flow: user does not have active flow")
	NoStepsDefined            = errors.New("flow: no steps defined")
)

// The Bus handles user actions, such as [telebot.OnText, telebot.OnMedia, etc].
// Flow processing involves handling any user prompt after the user has begun the flow process.
// We offer this contract to give developers more control, avoiding reliance on obscure mechanisms.
type Bus interface {
	// UserInFlow returns true if the user is currently engaged in flow processing.
	//
	// Example:
	//  bot.Handle("/start", func(c telebot.Context) error {
	//    if flowBus.UserInFlow(c.Sender().ID) { // Reply with an error message. }
	//
	//    return c.Reply("Hello!")
	//  })
	UserInFlow(userID int64) bool
	// UserContinueFlow initiates or continues the flow process for a user if one is already in progress.
	//
	// Example:
	//  bot.Handle("/start", func(c telebot.Context) error {
	//    if flowBus.UserInFlow(c.Sender().ID) { flowBus.UserContinueFlow(c.Sender().ID) }
	//
	//    return c.Reply("Hello!")
	//  })
	UserContinueFlow(userID int64, c telebot.Context) error
	// UserContinueFlowOrCustom calls [UserContinueFlow] if the flow process for a user is in progress.
	// Otherwise, it calls a custom function.
	// For instance, you may need to call this function to define a custom handler for any action required by the flow.
	//
	// Example:
	//
	//  bot.Handle(telebot.OnText, flowBus.ProcessUserToFlowOrCustom(func (c telebot.Context) error {
	//     // Called only if the user hasn't begun the flow.
	//
	//     return nil
	//   }))
	UserContinueFlowOrCustom(telebot.HandlerFunc) telebot.HandlerFunc
	// Handle implements any message handler.
	// This function checks if the user is continuing work on their active flow and processes it if so.
	//
	// Example:
	//  bot.Handle(telebot.OnText, flowBus.Handle)
	Handle(telebot.Context) error

	// Flow initiates flow configuration
	Flow(factory *Factory) telebot.HandlerFunc
}

// describes the state to the [SimpleBus.states] value
type flowState struct {
	// User's flow
	flow    *Flow
	state   State
	machine Machine
}

// SimpleBus implementation for the [Bus] contract
type SimpleBus struct {
	// Stores the active user flows by their IDs.
	// Key - user id (int64)
	// Value - the [state] instance
	states sync.Map
	// We don't need to keep active flows indefinitely.
	// This setting defines the maximum lifespan for each flow.
	// Background process will remove flows that have been alive longer than the defined duration.
	// @TODO: Provide a callback handler for every deletion process.
	flowSessionIsAvailableFor time.Duration
}

func (b *SimpleBus) UserInFlow(userID int64) bool {
	_, exists := b.states.Load(userID)

	return exists
}

func (b *SimpleBus) UserContinueFlow(userID int64, c telebot.Context) error {
	//flow, exists := b.states.Load(userID)
	_, exists := b.states.Load(userID)
	if !exists {
		return UserDoesNotHaveActiveFlow
	}

	// @TODO: call machine
	return nil
}

func (b *SimpleBus) UserContinueFlowOrCustom(handler telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if b.UserInFlow(c.Sender().ID) {
			return b.UserContinueFlow(c.Sender().ID, c)
		}

		return handler(c)
	}
}

func (b *SimpleBus) Handle(c telebot.Context) error {
	stV, exists := b.states.Load(c.Sender().ID)
	if !exists {
		return UserDoesNotHaveActiveFlow
	}

	st := stV.(flowState)
	// Update context for the state
	st.state.Add(StateContextKey, c)
	// Get active step
	step := st.flow.steps[st.machine.ActiveStep()]
	// Call validators if they are defined
	validators := step.validators
	if len(validators) > 0 {
		for _, validator := range validators {
			err := validator.Validate(st.state)
			if err != nil {
				if st.flow.UseValidatorErrorsAsUserResponse {
					return c.Reply(err.Error())
				} else {
					return err
				}
			}
		}
	}

	// Call [assign]
	if step.assign != nil {
		if err := step.assign(st.state); err != nil {
			return err
		}
	}

	activeStep := st.machine.ActiveStep()
	// Call [then] event if it's defined
	if step.then != nil {
		if err := step.then(st.state, &step); err != nil {
			return err
		}
	}

	// It was the last step. Call the [then] handler
	if len(st.flow.steps) <= st.machine.ActiveStep()+1 {
		b.removeState(c.Sender().ID)

		if st.flow.then == nil {
			return nil
		}

		return st.flow.then(st.state)
	}

	// Sometimes, the user may navigate through steps within handlers.
	// If this occurs, we don't need to call the [next] function because navigating
	// through the machine already triggers it.
	if activeStep == st.machine.ActiveStep() {
		// Process to the next step
		err := st.machine.Next(st.state)
		if err != nil {
			// Remove flow on any error occurring within flow logic.
			// We need to call the [Fail] function because, typically,
			// that handler should send something to the user like [Try again].
			b.removeState(c.Sender().ID)

			if st.flow.catch == nil {
				return nil
			}

			return st.flow.catch(st.state, err)
		}
	}

	return nil
}

// Remove [state] from the [SimpleBus.states]
func (b *SimpleBus) removeState(userID int64) {
	b.states.Delete(userID)
}

func (b *SimpleBus) Flow(factory *Factory) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if len(factory.flow.steps) == 0 {
			return NoStepsDefined
		}

		// If the user already has a flow, we need to recall the active step.
		stV, exists := b.states.Load(c.Sender().ID)
		if exists {
			st := stV.(flowState)
			// Update context
			st.state.Add(StateContextKey, c)

			return st.machine.ToStep(st.machine.ActiveStep(), st.state)
		}

		machine := NewMachine(factory.flow)
		state := NewRuntimeState(factory.userState).
			Add(StateContextKey, c).
			Add(StateMachineKey, machine)
		// Register flow for the user
		st := flowState{
			flow:    factory.flow,
			state:   state,
			machine: machine,
		}
		b.states.Store(c.Sender().ID, st)
		// Call the machine for the start the first step
		return machine.ToStep(0, st.state)
	}
}

// Removes flows that have been active for longer than [flowSessionIsAvailableFor] time.
func (b *SimpleBus) removeIdleFlows() {
	// @TODO: Provide an API for clients.
	// For example, a developer may want to notify a user that their session has expired.
}

func NewBus(flowSessionIsAvailableFor time.Duration) Bus {
	bus := &SimpleBus{
		flowSessionIsAvailableFor: flowSessionIsAvailableFor,
	}

	// @TODO: do we need to create an API for users to interact with this?
	go bus.removeIdleFlows()

	return bus
}
