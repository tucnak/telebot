package telebot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRaceStopClientChannel(t *testing.T) {
	t.Parallel()

	api, err := NewBot(Settings{Offline: true})
	require.NoError(t, err)

	syncChan := make(chan struct{})

	go func() {
		syncChan <- struct{}{}

		_, err = api.Raw("setMyCommands", CommandParams{
			Commands: []Command{
				{
					Text:        "/test",
					Description: "test",
				},
			},
			LanguageCode: "en",
		})
		require.EqualError(t, err, "telegram: Not Found (404)")

		close(syncChan)
	}()

	<-syncChan

	time.Sleep(time.Second)

	// act and assert
	go api.Start()
	defer func() {
		_, err = api.Close()
		require.NoError(t, err)
	}()

	<-syncChan
}
