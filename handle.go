package telebot

import (
	"regexp"
	"strings"
)

// Handle lets you set the handler for some command name or
// one of the supported endpoints.
//
// See Endpoint.
func (b *Bot) Handle(endpoint, handler interface{}) {
	if cmd, ok := endpoint.(string); ok {
		b.handlers[cmd] = handler

	} else if end, ok := endpoint.(Endpoint); ok {
		b.handlers[string(end)] = handler

	} else {
		panic("Handle() only supports patterns and endpoints")
	}
}

var cmdRx = regexp.MustCompile(`^\/(\w+)(@(\w+))?`)

func (b *Bot) handleMessages(messages chan Message) {
	for m := range messages {
		// Text message
		if m.Text != "" {
			match := cmdRx.FindAllStringSubmatch(m.Text, -1)

			// Command found
			if match != nil {
				if b.handleCommand(&m, match[0][1], match[0][3]) {
					continue
				}
			}

			// Feeding it to OnMessage if one exists.
			if handler, ok := b.handlers[string(OnMessage)]; ok {
				if handler, ok := handler.(func(*Message)); ok {
					go handler(&m)
					continue
				}
			}
		}
	}
}

func (b *Bot) handleCommand(m *Message, cmdName, cmdBot string) bool {
	// Group-syntax: "/cmd@bot"
	if cmdBot != "" && !strings.EqualFold(b.Me.Username, cmdBot) {
		return false
	}

	if handler, ok := b.handlers[cmdName]; ok {
		if handler, ok := handler.(func(*Message)); ok {
			go handler(m)
			return true
		}
	}

	return false
}

func (b *Bot) handleQueries(queries chan Query) {
	for _ = range queries {

	}
}

func (b *Bot) handleCallbacks(callbacks chan Callback) {
	for _ = range callbacks {

	}
}
