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
	Flow(endpoint interface{}, factory *Factory) error
}

// describes the state to the [SimpleBus.states] value
type flowState struct {
	// telegram bot endpoint
	endpoint interface{}
	// User's flow
	flow     *Flow
	state    State
	machine  Machine
	metaData *MetaData
}

// SimpleBus implementation for the [Bus] contract
type SimpleBus struct {
	bot *telebot.Bot
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

// Calls the meta functions for the step [validators/assign/etc].
func (b *SimpleBus) handleMetaForStep(st flowState, c telebot.Context, step Step) error {
	// Call validators if they are defined
	validators := step.validators
	if len(validators) > 0 {
		for _, validator := range validators {
			err := validator.Validate(st.state)
			// Fill metadata information on error
			if err != nil {
				st.metaData.FailureStage = MetaDataFailureStageValidation
				st.metaData.FailedError = err

				return err
			}
		}
	}

	// Call [assign]
	if step.assign != nil {
		if err := step.assign(st.state); err != nil {
			// Fill step metadata
			st.metaData.FailureStage = MetaDataFailureStageAssign
			st.metaData.FailedError = err

			return err
		}
	}

	return nil
}

// Call the [catch] function for the [flow] and remove the flow from the state.
func (b *SimpleBus) handleCatch(st flowState, c telebot.Context) error {
	// Remove flow on any error occurring within flow logic.
	// We need to call the [Fail] function because, typically,
	// that handler should send something to the user like [Try again].
	b.removeState(c.Sender().ID)

	if st.flow.catch == nil {
		return nil
	}

	return st.flow.catch(st.state, st.metaData)
}

func (b *SimpleBus) Handle(c telebot.Context) error {
	stV, exists := b.states.Load(c.Sender().ID)
	if !exists {
		return UserDoesNotHaveActiveFlow
	}

	st := stV.(flowState)
	// Update context for the state
	st.state.Add(StateContextKey, c)
	activeStep := st.machine.ActiveStep()
	// Get active step
	step := st.flow.steps[activeStep]
	// Begin filling metadata for the current step.
	st.metaData.LastActiveStep = StepMetaData{Step: st.machine.ActiveStep()}
	defer func() {
		// Update the flowState for the user only if [failedError] is nil.
		// Otherwise, if the flow failed and the [catch] function was called,
		// we don't need to update the flow because it no longer exists.
		if st.metaData.FailureStage == MetaDataFailureStageNone {
			b.states.Store(c.Sender().ID, st)
		}
	}()

	if err := b.handleMetaForStep(st, c, step); err != nil {
		st.metaData.FailedStep = &st.metaData.LastActiveStep

		return b.handleCatch(st, c)
	}

	// Since it is the last step, call the [then] handler.
	if len(st.flow.steps) <= activeStep+1 {
		// Call on each step handler if it is defined
		if st.flow.onEachStep != nil {
			st.flow.onEachStep(st.state, st.metaData.LastActiveStep)
		}

		if st.flow.then == nil {
			b.removeState(c.Sender().ID)

			return nil
		}

		// If an error is returned, we need to call [catch] for the flow.
		err := st.flow.then(st.state)
		if err != nil {
			// Fill step metadata
			st.metaData.FailureStage = MetaDataFailureStageThen
			st.metaData.FailedError = err

			return b.handleCatch(st, c)
		}

		return err
	}

	// Sometimes, the user may navigate through steps within handlers.
	// If this occurs, we don't need to call the [next] function because navigating
	// through the machine already triggers it.
	if activeStep == st.machine.ActiveStep() {
		// Process to the next step
		err := st.machine.Next(st.state)
		if err != nil {
			// Fill step metadata
			st.metaData.FailureStage = MetaDataFailureStageBegin
			st.metaData.FailedError = err

			return b.handleCatch(st, c)
		}

		// Call on each step handler if it is defined
		if st.flow.onEachStep != nil {
			st.flow.onEachStep(st.state, st.metaData.LastActiveStep)
		}
	}

	return nil
}

// Remove the [state] from the [SimpleBus.states].
func (b *SimpleBus) removeState(userID int64) {
	b.states.Delete(userID)
}

func (b *SimpleBus) Flow(endpoint interface{}, factory *Factory) error {
	if len(factory.flow.steps) == 0 {
		return NoStepsDefined
	}

	b.bot.Handle(endpoint, func(c telebot.Context) error {
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
			endpoint: endpoint,
			flow:     factory.flow,
			state:    state,
			machine:  machine,
			metaData: &MetaData{
				Endpoint: endpoint,
				// Sets the first step as the last active step.
				LastActiveStep: StepMetaData{Step: 0},
			},
		}
		b.states.Store(c.Sender().ID, st)

		// Call the machine for the start the first step
		return machine.ToStep(0, st.state)
	})

	return nil
}

// Removes flows that have been active for longer than [flowSessionIsAvailableFor] time.
func (b *SimpleBus) removeIdleFlows() {
	// @TODO: Provide an API for clients.
	// For example, a developer may want to notify a user that their session has expired.
}

func NewBus(bot *telebot.Bot, flowSessionIsAvailableFor time.Duration) Bus {
	bus := &SimpleBus{
		bot:                       bot,
		flowSessionIsAvailableFor: flowSessionIsAvailableFor,
	}

	// @TODO: do we need to create an API for users to interact with this?
	go bus.removeIdleFlows()

	return bus
}
