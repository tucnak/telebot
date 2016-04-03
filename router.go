package telebot

import (
	"errors"
	"regexp"
)

type Router struct {
	routes map[string]Controller
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]Controller)}
}

func (r *Router) When(pattern string, controller Controller) {
	r.routes[pattern] = controller
}

func (r *Router) GetController(message *Message) (Controller, *map[string]string, error) {
	for k, v := range r.routes {
		reg, err := regexp.Compile(k)
		if err != nil {
			return nil, nil, err
		}

		if match := reg.FindStringSubmatch(message.Text); len(match) > 0 {
			args := make(map[string]string)

			for i, name := range reg.SubexpNames() {
				if i != 0 {
					args[name] = match[i]
				}
			}
			return v, &args, nil
		}
	}

	return nil, nil, errors.New("No controller found")
}
