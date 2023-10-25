package telebot

func (b *Bot) getChanStopClient() chan struct{} {
	b.stopClientMux.RLock()
	defer b.stopClientMux.RUnlock()

	return b.stopClient
}

func (b *Bot) closeChanStopClient() {
	b.stopClientMux.Lock()
	defer b.stopClientMux.Unlock()
	
	close(b.stopClient)
}

func (b *Bot) createNewChanStopClient() {
	b.stopClientMux.Lock()
	defer b.stopClientMux.Unlock()

	b.stopClient = make(chan struct{})
}

func (b *Bot) destroyChanStopClient() {
	b.stopClientMux.Lock()
	defer b.stopClientMux.Unlock()

	b.stopClient = nil
}
