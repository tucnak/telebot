package telebot

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPoller struct {
	updates chan Update
	done    chan struct{}
}

func newTestPoller() *testPoller {
	return &testPoller{
		updates: make(chan Update, 1),
		done:    make(chan struct{}, 1),
	}
}

func (p *testPoller) Poll(b *Bot, updates chan Update, stop chan struct{}) {
	for {
		select {
		case upd := <-p.updates:
			updates <- upd
		case <-stop:
			return
		default:
		}
	}
}

func TestMiddlewarePoller(t *testing.T) {
	tp := newTestPoller()
	var ids []int

	pref := defaultSettings()
	pref.Offline = true

	b, err := NewBot(pref)
	if err != nil {
		t.Fatal(err)
	}

	b.Poller = NewMiddlewarePoller(tp, func(u *Update) bool {
		if u.ID > 0 {
			ids = append(ids, u.ID)
			return true
		}

		tp.done <- struct{}{}
		return false
	})

	go func() {
		tp.updates <- Update{ID: 1}
		tp.updates <- Update{ID: 2}
		tp.updates <- Update{ID: 0}
	}()

	go b.Start(context.Background())
	<-tp.done
	b.Stop()

	assert.Contains(t, ids, 1)
	assert.Contains(t, ids, 2)
}
