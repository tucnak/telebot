package telebot

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStickerSet(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for StickerSet methods test")
	}

	input := []InputSticker{
		{
			File:     FromURL("https://placehold.co/512/000000/FFFFFF/png"),
			Emojis:   []string{"🤖"},
			Keywords: []string{"telebot", "robot", "bot"},
		},
		{
			File:     FromURL("https://placehold.co/512/000000/999999/png"),
			Emojis:   []string{"🤖"},
			Keywords: []string{"telebot", "robot", "bot"},
		},
	}

	original := &StickerSet{
		Name:   fmt.Sprintf("telebot_%d_by_%s", time.Now().Unix(), b.Me.Username),
		Type:   StickerRegular,
		Format: StickerStatic,
		Title:  "Telebot Stickers",
		Input:  input[:1],
	}

	// 1
	err := b.CreateStickerSet(user, original)
	require.NoError(t, err)
	// 2
	err = b.AddStickerToSet(user, original.Name, input[1])
	require.NoError(t, err)

	original.Thumbnail = &Photo{File: thumb}
	err = b.SetStickerSetThumb(user, original)
	require.NoError(t, err)

	set, err := b.StickerSet(original.Name)
	require.NoError(t, err)
	require.Equal(t, original.Name, set.Name)
	require.Equal(t, len(input), len(set.Stickers))

	_, err = b.Send(user, &set.Stickers[0])
	require.NoError(t, err)

	_, err = b.Send(user, &set.Stickers[1])
	require.NoError(t, err)
}
