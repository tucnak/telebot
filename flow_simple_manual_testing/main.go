package main

import (
	"errors"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/flow"
	"log"
	"reflect"
	"strconv"
	"time"
)

// Validators: I'm certain that we need to implement some basic validators for users.

type NonEmptyValidator struct{}

func (NonEmptyValidator) Validate(state flow.State) error {
	if len(state.Read(flow.StateContextKey).(tele.Context).Message().Text) < 2 {
		return errors.New("message is required")
	}

	return nil
}

type BadValidator struct{}

func (BadValidator) Validate(state flow.State) error {
	return errors.New("test")
	//return state.Machine.Fail(state)
}

var (
	nonEmptyValidator = NonEmptyValidator{}
	badValidator      = BadValidator{}
)

// TextAssigner I'm certain that we need to implement some basic assigner for users.
func TextAssigner(value interface{}) flow.StateHandler {
	vai := reflect.ValueOf(value)
	vai = vai.Elem()

	return func(state flow.State) error {
		text := state.Read(flow.StateContextKey).(tele.Context).Text()
		if len(text) == 0 {
			return nil
		}

		switch vai.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(text, 10, 32)
			if err != nil {
				return err
			}

			vai.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(text, 10, 32)
			if err != nil {
				return err
			}

			vai.SetUint(n)
		// ...floating-point and complex cases omitted for brevity...
		case reflect.Bool:
			n, err := strconv.ParseBool(text)
			if err != nil {
				return err
			}

			vai.SetBool(n)
		case reflect.String:
			vai.SetString(text)
		//case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		//	return v.Type().String() + " 0x" +
		//		strconv.FormatUint(uint64(v.Pointer()), 16)
		default: // reflect.Array, reflect.Struct, reflect.Interface
			vai.SetString(text)
		case reflect.Invalid:
			return errors.New("invalid type")
		}

		return nil
	}
}

func main() {
	pref := tele.Settings{
		Token:  "",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	sendUserMessage := func(message string) func(flow.State) error {
		return func(state flow.State) error {
			return state.Read(flow.StateContextKey).(tele.Context).Reply(message)
		}
	}
	stepCompletedLogging := func(state flow.State, step *flow.Step) error {
		log.Println("Step completed")

		return nil
	}

	// Configure flow bus
	flowBus := flow.NewBus(5 * time.Minute)
	// Handle any text by flow bus
	b.Handle(tele.OnText, flowBus.Handle)
	// First flow
	var email string
	b.Handle("/start", flowBus.Flow(
		flow.New().
			Next(
				flow.NewStep(sendUserMessage("Enter email:")).
					Validate(nonEmptyValidator).
					Assign(TextAssigner(&email)).
					Then(stepCompletedLogging),
			).
			Next(
				flow.NewStep(sendUserMessage("Enter password:")).
					Validate(nonEmptyValidator).
					Assign(TextAssigner(&email)),
			).
			Next(
				flow.NewStep(sendUserMessage("Third step:")).
					Then(func(state flow.State, step *flow.Step) error {
						//return state.Read(flow.StateMachineKey).(flow.Machine).ToStep(0, state)

						return nil
					}),
			).
			Then(func(state flow.State) error {
				log.Println("Steps are completed!")

				return state.Read(flow.StateContextKey).(tele.Context).Reply("Done")
			}).
			Catch(func(state flow.State, err error) error {
				log.Println("FAILED: ", err)

				return nil
			}),
	))

	// Flow using state storage
	type user struct {
		email    string
		password string
	}
	userStorageKey := "user"
	b.Handle("/start2", flowBus.Flow(
		flow.New().
			AddState(userStorageKey, &user{}).
			Next(
				flow.NewStep(sendUserMessage("Enter email:")).
					Validate(nonEmptyValidator).
					Assign(func(state flow.State) error {
						u := state.Read(userStorageKey).(*user)
						u.email = state.Read(flow.StateContextKey).(tele.Context).Message().Text

						return nil
					}),
			).
			Next(
				flow.NewStep(sendUserMessage("Enter password:")).
					Validate(nonEmptyValidator).
					Assign(func(state flow.State) error {
						u := state.Read(userStorageKey).(*user)
						u.password = state.Read(flow.StateContextKey).(tele.Context).Message().Text

						return nil
					}).
					Then(func(state flow.State, step *flow.Step) error {
						log.Println("Second step successfully passed!")

						return nil
					}),
			).
			Next(
				flow.NewStep(func(state flow.State) error {
					return errors.New("should be passed to the [Catch]")
				}),
			).
			Then(func(state flow.State) error {
				log.Println("Steps are completed!. User: ", state.Read(userStorageKey))

				return state.Read(flow.StateContextKey).(tele.Context).Reply("Done")
			}).
			Catch(func(state flow.State, err error) error {
				log.Println("FAILED: ", err)

				return nil
			}),
	))

	b.Start()
}
