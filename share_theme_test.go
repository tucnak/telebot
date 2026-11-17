package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatShareTheme(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChat", json.RawMessage(`{"id":7,"type":"private","share_theme":{"type":"emoji",
		"emojis":[{"type":"emoji","emoji":"🎉"},{"type":"custom_emoji","custom_emoji_id":"42"}],"background_backdrop_id":3}}`))
	srv.Reply("getChat", json.RawMessage(`{"id":7,"type":"private","share_theme":{"type":"unique_gift","gift":`+uniqueGiftJSON+`}}`))

	chat, err := b.ChatByID(7)
	require.NoError(t, err)
	assert.Equal(t, &ChatShareTheme{
		Type:                 ShareThemeEmoji,
		Emojis:               []EmojiType{{Type: "emoji", Emoji: "🎉"}, {Type: "custom_emoji", CustomEmojiID: "42"}},
		BackgroundBackdropID: 3,
	}, chat.ShareTheme)

	info, err := b.ChatFullInfo(ChatID(7))
	require.NoError(t, err)
	assert.Equal(t, ShareThemeUniqueGift, info.ShareTheme.Type)
	assert.Equal(t, "Cake-7", info.ShareTheme.Gift.Name)
}

func TestShareThemeBackdrops(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChatShareThemeBackdrops", json.RawMessage(`{"backdrops":[{"id":1,"name":"Sky",
		"colors":{"center_color":1,"edge_color":2,"symbol_color":3,"text_color":4}}]}`))

	backdrops, err := b.ShareThemeBackdrops()
	require.NoError(t, err)
	assert.Equal(t, []ChatShareThemeBackdrop{{ID: 1, Name: "Sky",
		Colors: UniqueGiftBackdropColors{CenterColor: 1, EdgeColor: 2, SymbolColor: 3, TextColor: 4}}}, backdrops)
	assert.Equal(t, "getChatShareThemeBackdrops", srv.Last().Method)
	assert.Empty(t, srv.Last().Params)
}
