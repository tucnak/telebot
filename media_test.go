package telebot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlbumSetCaption(t *testing.T) {
	tests := []struct {
		name  string
		media Inputtable
	}{
		{
			name:  "photo",
			media: &Photo{Caption: "wrong_caption"},
		},
		{
			name:  "animation",
			media: &Animation{Caption: "wrong_caption"},
		},
		{
			name:  "video",
			media: &Video{Caption: "wrong_caption"},
		},
		{
			name:  "audio",
			media: &Audio{Caption: "wrong_caption"},
		},
		{
			name:  "document",
			media: &Document{Caption: "wrong_caption"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a Album
			a = append(a, tt.media)
			a = append(a, &Photo{Caption: "random_caption"})
			a.SetCaption("correct_caption")
			assert.Equal(t, "correct_caption", a[0].InputMedia().Caption)
			assert.Equal(t, "random_caption", a[1].InputMedia().Caption)
		})
	}
}
