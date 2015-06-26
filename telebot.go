// Package telebot provides a handy wrapper for interactions
// with Telegram bots.
package telebot

// Attempts to construct a Bot with `token` given.
func Create(token string) (Bot, error) {
	user, err := api_getMe(token)
	if err != nil {
		return Bot{}, err
	}

	return Bot{
		Token:    token,
		Identity: user,
	}, nil
}
