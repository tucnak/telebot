package telebot

type testPoller struct {
	updates chan Update
}

func (p *testPoller) Poll(b *Bot, updates chan Update, stop chan struct{}) {

	for {
		select {
		case upd := <-p.updates:
			updates <- upd
		case <-stop:
			close(stop)
			return
		default:
		}
	}
}
