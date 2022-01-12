package telebot

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
