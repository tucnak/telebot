package telebot

import (
	"time"

	"github.com/pkg/errors"
)

// Poller is a provider of Updates.
//
// All pollers must implement Poll(), which accepts bot
// pointer and subscription channel and start polling
// asynchronously straight away.
type Poller interface {
	// Poll is supposed to take the bot object
	// subscription channel and start polling
	// for Updates immediately.
	Poll(b *Bot, dest chan Update)
}

// LongPoller is a classic LongPoller with timeout.
type LongPoller struct {
	Timeout time.Duration
}

// Poll does long polling.
func (p *LongPoller) Poll(b *Bot, dest chan Update) {
	var latestUpd int

	for {
		updates, err := b.getUpdates(latestUpd+1, p.Timeout)

		if err != nil {
			b.debug(errors.Wrap(err, "getUpdates() failed"))
			continue
		}

		for _, update := range updates {
			latestUpd = update.ID
			dest <- update
		}
	}
}
