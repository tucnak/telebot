package telebot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Bot struct {
	Token string

	// A "user" behind bot instance
	Identity User
}

func Create(token string) (Bot, error) {
	request := "https://api.telegram.org/bot" + token + "/getMe"

	resp, err := http.Get(request)
	if err != nil {
		return Bot{}, err
	}

	defer resp.Body.Close()
	me_json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Bot{}, err
	}

	var bot_info struct {
		Ok          bool
		Result      User
		Description string
	}

	err = json.Unmarshal(me_json, &bot_info)
	if err != nil {
		return Bot{}, err
	}

	if bot_info.Ok {
		return Bot{
			Token:    token,
			Identity: bot_info.Result,
		}, nil
	} else {
		return Bot{}, AuthError{bot_info.Description}
	}
}
