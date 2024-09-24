package flow

import (
	"errors"
	"fmt"
)

// Machine describes the contract for the flow handling
type Machine interface {
	// Back move backward by step
	Back(state State) error
	// Next move forward by step
	Next(state State) error
	// ToStep Move to the step
	ToStep(step int, state State) error
	// ActiveStep returns the current step
	ActiveStep() int
}

// SimpleMachine implements the [Machine] contract
type SimpleMachine struct {
	// User defined flow
	flow *Flow
	// Active step for the user
	activeStep int
}

func (m *SimpleMachine) Back(state State) error {
	if m.activeStep <= 0 {
		return errors.New("already first step")
	}

	m.activeStep -= 1

	return m.run(state)
}

func (m *SimpleMachine) Next(state State) error {
	if m.activeStep+1 >= len(m.flow.steps) {
		return errors.New("already last step")
	}

	m.activeStep += 1

	return m.run(state)
}

func (m *SimpleMachine) ToStep(step int, state State) error {
	if step < 0 {
		return errors.New("step cannot be less than zero")
	}

	if step > len(m.flow.steps) {
		return errors.New("step cannot be greater than steps count")
	}

	m.activeStep = step

	return m.run(state)
}

func (m *SimpleMachine) ActiveStep() int {
	return m.activeStep
}

// Run the current step (this function should be called by [Back]/[Next]/[ToStep] functions).
func (m *SimpleMachine) run(state State) error {
	if len(m.flow.steps) < m.activeStep {
		return errors.New(fmt.Sprintf("step isn't defined (%d)", m.activeStep))
	}

	step := m.flow.steps[m.activeStep]
	if step.handler != nil {
		return step.handler(state)
	}

	return nil
}

func NewMachine(flow *Flow) Machine {
	return &SimpleMachine{
		flow:       flow,
		activeStep: 0,
	}
}
