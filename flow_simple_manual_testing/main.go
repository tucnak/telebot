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

func (NonEmptyValidator) Validate(state *flow.State) error {
	if len(state.Context.Message().Text) < 2 {
		return errors.New("message is required")
	}

	return nil
}

type BadValidator struct{}

func (BadValidator) Validate(state *flow.State) error {
	return nil
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

	return func(state *flow.State) error {
		text := state.Context.Text()
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
		Token:  "5931155624:AAGLxTOnMt2O3UYLGpZSZAacxBVJONO1UP4",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	sendUserMessage := func(message string) func(state *flow.State) error {
		return func(state *flow.State) error {
			return state.Context.Reply(message)
		}
	}

	// Configure flow bus
	flowBus := flow.NewBus(b, 5*time.Minute)
	// Handle any text by flow bus
	b.Handle(tele.OnText, flowBus.Handle)
	// First flow
	var email string
	b.Handle("/start", flowBus.Flow(
		flow.New().
			Step(
				flow.NewStep().
					Begin(sendUserMessage("Enter email:")).
					Validate(nonEmptyValidator).
					Assign(TextAssigner(&email)),
			).
			Step(
				flow.NewStep().
					Begin(sendUserMessage("Enter password:")).
					Validate(badValidator).
					Success(func(state *flow.State) error {
						log.Println("Second step successfully passed!")

						return nil
					}),
			).
			Step(
				flow.NewStep().
					Begin(func(state *flow.State) error {
						return nil
						//return state.Machine.Back(state)
						//return state.Machine.ToStep(0, state)
						//return state.Machine.Next(state)
					}),
			).
			UseValidatorErrorsAsUserResponse(true).
			Fail(func(state *flow.State, err error) error {
				log.Println("Something get wrong: ", err)

				return nil
			}).
			Success(func(state *flow.State) error {
				log.Println(email)

				return state.Context.Reply("You have successfully completed all the steps!")
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
			Step(
				flow.NewStep().
					Begin(sendUserMessage("Enter email:")).
					Validate(nonEmptyValidator).
					Assign(func(state *flow.State) error {
						value, _ := state.Get(userStorageKey)
						userValue := value.(*user)
						userValue.email = state.Context.Message().Text

						return nil
					}),
			).
			Step(
				flow.NewStep().
					Begin(sendUserMessage("Enter password:")).
					Validate(badValidator).
					Assign(func(state *flow.State) error {
						value, _ := state.Get(userStorageKey)
						userValue := value.(*user)
						userValue.password = state.Context.Message().Text

						return nil
					}).
					Success(func(state *flow.State) error {
						log.Println("Second step successfully passed!")

						return nil
					}),
			).
			Step(
				flow.NewStep().
					Begin(func(state *flow.State) error {
						return state.Machine.Fail(state, errors.New("should be passed to the [Fail]"))
						//return state.Machine.Back(state)
						//return state.Machine.ToStep(0, state)
						//return state.Machine.Next(state)
					}),
			).
			UseValidatorErrorsAsUserResponse(true).
			Fail(func(state *flow.State, err error) error {
				log.Println("Something get wrong: ", err)

				return nil
			}).
			Success(func(state *flow.State) error {
				log.Println(state.Get(userStorageKey))

				return state.Context.Reply("You have successfully completed all the steps!")
			}),
	))

	b.Start()
}
