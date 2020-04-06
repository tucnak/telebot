package telebot

type testPoller struct {
	Message string
}

func (p *testPoller) Poll(b *Bot, updates chan Update, stop chan struct{}) {
	updates <- Update{Message: &Message{Text: p.Message}}

	for {
		select {
		case <-stop:
			close(stop)
			return
		default:
		}
	}
}
