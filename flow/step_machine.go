package flow

import (
	"errors"
	"fmt"
)

type Machine interface {
	// Back move backward by step.
	Back(state *State) error

	// Next move forward by step.
	Next(state *State) error

	// ToStep Move to the step.
	ToStep(step int, state *State) error

	// Ok run [Flow.okFn].
	Ok(state *State) error

	// Fail sets failed and run [Flow.errFn].
	Fail(state *State, err error) error

	// ActiveStep returns the current step.
	ActiveStep() int
}

type machine struct {
	flow       *Flow
	activeStep int

	failed bool
}

func (m *machine) Back(state *State) error {
	if m.activeStep-1 <= 0 {
		return errors.New("already first step")
	}

	m.activeStep -= 1

	return m.run(state)
}

func (m *machine) Next(state *State) error {
	if m.activeStep+1 >= len(m.flow.steps) {
		return errors.New("already last step")
	}

	m.activeStep += 1

	return m.run(state)
}

func (m *machine) ToStep(step int, state *State) error {
	if step < 0 {
		return errors.New("step cannot be less than zero")
	}

	if step > len(m.flow.steps) {
		return errors.New("step cannot be greater than steps count")
	}

	m.activeStep = step

	return m.run(state)
}

func (m *machine) Ok(state *State) error {
	if m.failed {
		return errors.New("flow was already failed")
	}

	if m.flow.okFn != nil {
		return m.flow.okFn(state)
	}

	return nil
}

func (m *machine) Fail(state *State, err error) error {
	m.failed = true

	if m.flow.errFn != nil {
		return m.flow.errFn(state, err)
	}

	return nil
}

func (m *machine) ActiveStep() int {
	return m.activeStep
}

func (m *machine) run(state *State) error {
	if m.failed {
		return errors.New("flow was already failed")
	}

	if len(m.flow.steps) < m.activeStep {
		return errors.New(fmt.Sprintf("step isn't defined (%d)", m.activeStep))
	}

	step := m.flow.steps[m.activeStep]
	if step.beginFn != nil {
		return step.beginFn(state)
	}

	return nil
}

// NewMachine constructor [Machine].
func NewMachine(flow *Flow) Machine {
	return &machine{
		flow:       flow,
		activeStep: 0,
	}
}
