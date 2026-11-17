package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStories(t *testing.T) {
	b, srv := newFakeBot(t)
	story := json.RawMessage(`{"chat":{"id":7,"type":"private"},"id":11}`)
	srv.Reply("postStory", story)
	srv.Reply("editStory", story)
	srv.Reply("repostStory", story)

	posted, err := b.PostStory("bc", StoryParams{
		Content:      &InputStoryContentVideo{Video: FromReader(strings.NewReader("mp4")), Duration: 4.5, IsAnimation: true},
		ActivePeriod: 86400,
		Caption:      "*hi*",
		ParseMode:    ModeMarkdownV2,
		Areas: []StoryArea{
			{Position: AreaPosition{X: 50, Y: 50, Width: 10, Height: 10}, Type: StoryAreaType{Type: StoryAreaLink, URL: "https://example.com"}},
			{Type: StoryAreaType{Type: StoryAreaSuggestedReaction, ReactionType: &Reaction{Type: ReactionTypeEmoji, Emoji: "👍"}, IsDark: true}},
		},
		PostToChatPage: true,
	})
	require.NoError(t, err)
	assert.Equal(t, 11, posted.ID)
	assert.Equal(t, int64(7), posted.Poster.ID)

	req := srv.Last()
	assert.Equal(t, "postStory", req.Method)
	assert.Equal(t, "mp4", string(req.Files["story"].Data))
	assert.JSONEq(t, `{"type":"video","video":"attach://story","duration":4.5,"is_animation":true}`, req.Params["content"])
	assert.Equal(t, "86400", req.Params["active_period"])
	assert.Equal(t, "MarkdownV2", req.Params["parse_mode"])
	assert.Equal(t, "true", req.Params["post_to_chat_page"])
	assert.NotContains(t, req.Params, "protect_content")
	var areas []map[string]interface{}
	req.JSON(t, "areas", &areas)
	assert.Equal(t, map[string]interface{}{"type": "link", "url": "https://example.com"}, areas[0]["type"])
	assert.Equal(t, map[string]interface{}{"type": "suggested_reaction", "reaction_type": map[string]interface{}{"type": "emoji", "emoji": "👍"}, "is_dark": true}, areas[1]["type"])

	_, err = b.EditStory("bc", 11, StoryParams{Content: &InputStoryContentPhoto{Photo: FromReader(strings.NewReader("jpg"))}, Caption: "new"})
	require.NoError(t, err)
	req = srv.Last()
	assert.Equal(t, "editStory", req.Method)
	assert.Equal(t, "11", req.Params["story_id"])
	assert.JSONEq(t, `{"type":"photo","photo":"attach://story"}`, req.Params["content"])
	assert.NotContains(t, req.Params, "active_period")

	_, err = b.RepostStory("bc2", posted, StoryParams{ActivePeriod: 6 * 3600, Protected: true})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"business_connection_id": "bc2", "from_chat_id": "7", "from_story_id": "11",
		"active_period": "21600", "protect_content": "true"}, srv.Last().Params)

	require.NoError(t, b.DeleteStory("bc", 11))
	assert.Equal(t, map[string]string{"business_connection_id": "bc", "story_id": "11"}, srv.Last().Params)

	_, err = b.PostStory("bc", StoryParams{})
	assert.Error(t, err)
}

func TestStoryAreaDecode(t *testing.T) {
	var a StoryArea
	require.NoError(t, json.Unmarshal([]byte(`{"position":{"x_percentage":1,"y_percentage":2,"width_percentage":3,"height_percentage":4,"rotation_angle":5,"corner_radius_percentage":6},
		"type":{"type":"location","latitude":1.5,"longitude":2.5,"address":{"country_code":"DE","city":"Berlin"}}}`), &a))
	assert.Equal(t, StoryAreaLocation, a.Type.Type)
	assert.Equal(t, "Berlin", a.Type.Address.City)
	assert.Equal(t, 6.0, a.Position.CornerRadius)
}

func TestStoryAreaRequiredZeroValues(t *testing.T) {
	// Telegram rejects an area without the fields its type requires.
	data, err := json.Marshal([]StoryAreaType{
		{Type: StoryAreaWeather, Emoji: "❄"},
		{Type: StoryAreaLocation},
		{Type: StoryAreaLink, URL: "https://example.com"},
	})
	require.NoError(t, err)
	assert.JSONEq(t, `[
		{"type":"weather","temperature":0,"emoji":"❄","background_color":0},
		{"type":"location","latitude":0,"longitude":0},
		{"type":"link","url":"https://example.com"}
	]`, string(data))
}
