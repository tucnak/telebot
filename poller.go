package telebot

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrCouldNotUpdate = errors.New("getUpdates() failed")
)

// Poller is a provider of Updates.
//
// All pollers must implement Poll(), which accepts bot
// pointer and subscription channel and start polling
// synchronously straight away.
type Poller interface {
	// Poll is supposed to take the bot object
	// subscription channel and start polling
	// for Updates immediately.
	//
	// Poller must listen for stop constantly and close
	// it as soon as it's done polling.
	Poll(b *Bot, updates chan Update, stop chan struct{})
}

// MiddlewarePoller is a special kind of poller that acts
// like a filter for updates. It could be used for spam
// handling, banning or whatever.
//
// For heavy middleware, use increased capacity.
//
type MiddlewarePoller struct {
	Capacity int // Default: 1
	Poller   Poller
	Filter   func(*Update) bool
}

// NewMiddlewarePoller wait for it... constructs a new middleware poller.
func NewMiddlewarePoller(original Poller, filter func(*Update) bool) *MiddlewarePoller {
	return &MiddlewarePoller{
		Poller: original,
		Filter: filter,
	}
}

// Poll sieves updates through middleware filter.
func (p *MiddlewarePoller) Poll(b *Bot, dest chan Update, stop chan struct{}) {
	cap := 1
	if p.Capacity > 1 {
		cap = p.Capacity
	}

	middle := make(chan Update, cap)
	stopPoller := make(chan struct{})

	go p.Poller.Poll(b, middle, stopPoller)

	for {
		select {
		// call to stop
		case <-stop:
			stopPoller <- struct{}{}

		// poller is done
		case <-stopPoller:
			close(stop)
			return

		case upd := <-middle:
			if p.Filter(&upd) {
				dest <- upd
			}
		}
	}
}

// LongPoller is a classic LongPoller with timeout.
type LongPoller struct {
	Timeout time.Duration

	LastUpdateID int
}

// Poll does long polling.
func (p *LongPoller) Poll(b *Bot, dest chan Update, stop chan struct{}) {
	go func(stop chan struct{}) {
		<-stop
		close(stop)
	}(stop)

	for {
		updates, err := b.getUpdates(p.LastUpdateID+1, p.Timeout)

		if err != nil {
			b.debug(ErrCouldNotUpdate)
			continue
		}

		for _, update := range updates {
			p.LastUpdateID = update.ID
			dest <- update
		}
	}
}
