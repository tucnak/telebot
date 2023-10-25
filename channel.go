package telebot

func (b *Bot) getChanStopClient() chan struct{} {
	return b.stopClient
}

func (b *Bot) createNewChanStopClient() {
	b.stopClient = make(chan struct{})
}

func (b *Bot) closeChanStopClient() {
	close(b.stopClient)
}

func (b *Bot) destroyChanStopClient() {
	b.stopClient = nil
}
