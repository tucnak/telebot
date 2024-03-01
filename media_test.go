package telebot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlbumSetCaption(t *testing.T) {
	var a Album
	a = append(a, &Photo{Caption: "wrong_caption"})
	a = append(a, &Photo{Caption: "t"})
	a.SetCaption("correct_caption")
	assert.Equal(t, "correct_caption", a[0].InputMedia().Caption)
	assert.Equal(t, "t", a[1].InputMedia().Caption)

	a = Album{}
	a = append(a, &Animation{Caption: "wrong_caption"})
	a = append(a, &Photo{Caption: "t"})
	a.SetCaption("correct_caption")
	assert.Equal(t, "correct_caption", a[0].InputMedia().Caption)
	assert.Equal(t, "t", a[1].InputMedia().Caption)

	a = Album{}
	a = append(a, &Audio{Caption: "wrong_caption"})
	a = append(a, &Photo{Caption: "t"})
	a.SetCaption("correct_caption")
	assert.Equal(t, "correct_caption", a[0].InputMedia().Caption)
	assert.Equal(t, "t", a[1].InputMedia().Caption)

	a = Album{}
	a = append(a, &Document{Caption: "wrong_caption"})
	a = append(a, &Photo{Caption: "t"})
	a.SetCaption("correct_caption")
	assert.Equal(t, "correct_caption", a[0].InputMedia().Caption)
	assert.Equal(t, "t", a[1].InputMedia().Caption)

	a = Album{}
	a = append(a, &Video{Caption: "wrong_caption"})
	a = append(a, &Photo{Caption: "t"})
	a.SetCaption("correct_caption")
	assert.Equal(t, "correct_caption", a[0].InputMedia().Caption)
	assert.Equal(t, "t", a[1].InputMedia().Caption)
}
