package telebot

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Stories are posted, edited and deleted on behalf of a managed business
// account (Bot API 9.0); the bot needs the can_manage_stories business right.

// InputStoryContent is the content of a story to post: an
// *InputStoryContentPhoto or an *InputStoryContentVideo.
//
// Bot API: InputStoryContent.
type InputStoryContent interface {
	storyContent(files map[string]File) (richWire, error)
}

// InputStoryContentPhoto is a 1080x1920 photo of up to 10 MB to post as a
// story. It is always uploaded.
//
// Bot API: InputStoryContentPhoto.
type InputStoryContentPhoto struct {
	Photo File
}

// InputStoryContentVideo is a 720x1280 streamable H.265 video of up to 60
// seconds and 30 MB to post as a story. It is always uploaded.
//
// Bot API: InputStoryContentVideo.
type InputStoryContentVideo struct {
	Video File

	// (Optional) Precise duration and the time of the cover frame, in
	// seconds.
	Duration            float64
	CoverFrameTimestamp float64

	// (Optional) The video has no sound.
	IsAnimation bool
}

func (c *InputStoryContentPhoto) storyContent(files map[string]File) (richWire, error) {
	repr := c.Photo.process("story", files)
	if repr == "" {
		return nil, errors.New("telebot: story photo does not exist")
	}
	return richWire{"type": "photo", "photo": repr}, nil
}

func (c *InputStoryContentVideo) storyContent(files map[string]File) (richWire, error) {
	repr := c.Video.process("story", files)
	if repr == "" {
		return nil, errors.New("telebot: story video does not exist")
	}
	w := richWire{"type": "video", "video": repr}
	if c.Duration != 0 {
		w["duration"] = c.Duration
	}
	if c.CoverFrameTimestamp != 0 {
		w["cover_frame_timestamp"] = c.CoverFrameTimestamp
	}
	w.set("is_animation", c.IsAnimation)
	return w, nil
}

// StoryArea is a clickable area on a story.
//
// Bot API: StoryArea.
type StoryArea struct {
	Position AreaPosition  `json:"position"`
	Type     StoryAreaType `json:"type"`
}

// Story area types.
const (
	StoryAreaLocation          = "location"
	StoryAreaSuggestedReaction = "suggested_reaction"
	StoryAreaLink              = "link"
	StoryAreaWeather           = "weather"
	StoryAreaUniqueGift        = "unique_gift"
)

// StoryAreaType is what a story area points to: a location, a suggested
// reaction, a link, weather, or a unique gift. Only the fields of Type are
// used.
//
// Bot API: StoryAreaType, StoryAreaTypeLocation, StoryAreaTypeSuggestedReaction, StoryAreaTypeLink, StoryAreaTypeWeather, StoryAreaTypeUniqueGift.
type StoryAreaType struct {
	Type string `json:"type"`

	// "location": coordinates and an optional address.
	Latitude  float64          `json:"latitude,omitempty"`
	Longitude float64          `json:"longitude,omitempty"`
	Address   *LocationAddress `json:"address,omitempty"`

	// "suggested_reaction": the reaction and its look.
	ReactionType *Reaction `json:"reaction_type,omitempty"`
	IsDark       bool      `json:"is_dark,omitempty"`
	IsFlipped    bool      `json:"is_flipped,omitempty"`

	// "link": an HTTP or tg:// URL.
	URL string `json:"url,omitempty"`

	// "weather": degrees Celsius, an emoji and an ARGB background color.
	Temperature     float64 `json:"temperature,omitempty"`
	Emoji           string  `json:"emoji,omitempty"`
	BackgroundColor int     `json:"background_color,omitempty"`

	// "unique_gift": the unique name of the gift.
	Name string `json:"name,omitempty"`
}

// MarshalJSON encodes the area with the fields its Type requires even when
// they are zero, such as a temperature of 0 °C or a point on the equator.
func (a StoryAreaType) MarshalJSON() ([]byte, error) {
	type alias StoryAreaType
	switch a.Type {
	case StoryAreaLocation:
		return json.Marshal(struct {
			alias
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{alias(a), a.Latitude, a.Longitude})
	case StoryAreaWeather:
		return json.Marshal(struct {
			alias
			Temperature     float64 `json:"temperature"`
			BackgroundColor int     `json:"background_color"`
		}{alias(a), a.Temperature, a.BackgroundColor})
	}
	return json.Marshal(alias(a))
}

// LocationAddress is the physical address of a location.
//
// Bot API: LocationAddress.
type LocationAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	Street      string `json:"street,omitempty"`
}

// StoryParams describes a story to post, edit or repost; each method uses
// the fields the Bot API takes for it.
type StoryParams struct {
	// Content of the story; for PostStory and EditStory.
	Content InputStoryContent

	// Seconds before the story is archived: 6*3600, 12*3600, 86400 or
	// 2*86400; for PostStory and RepostStory.
	ActivePeriod int

	// (Optional) Caption, 0-2048 characters, with its parse mode or
	// entities, and clickable areas; for PostStory and EditStory.
	Caption   string
	ParseMode ParseMode
	Entities  []MessageEntity
	Areas     []StoryArea

	// (Optional) Keep the story on the chat page after it expires, and
	// protect it from forwarding and screenshots; for PostStory and
	// RepostStory.
	PostToChatPage bool
	Protected      bool
}

func (p StoryParams) embed(params map[string]string, files map[string]File) error {
	if p.Content != nil {
		content, err := p.Content.storyContent(files)
		if err != nil {
			return err
		}
		data, _ := json.Marshal(content)
		params["content"] = string(data)
	}
	if p.Caption != "" {
		params["caption"] = p.Caption
	}
	if len(p.Entities) > 0 {
		data, _ := json.Marshal(p.Entities)
		params["caption_entities"] = string(data)
	} else if p.ParseMode != ModeDefault {
		params["parse_mode"] = p.ParseMode
	}
	if len(p.Areas) > 0 {
		data, _ := json.Marshal(p.Areas)
		params["areas"] = string(data)
	}
	return nil
}

func (b *Bot) story(method string, files map[string]File, params map[string]string) (*Story, error) {
	data, err := b.sendFiles(method, files, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *Story
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

func embedStoryFlags(params map[string]string, p StoryParams) {
	params["active_period"] = strconv.Itoa(p.ActivePeriod)
	if p.PostToChatPage {
		params["post_to_chat_page"] = "true"
	}
	if p.Protected {
		params["protect_content"] = "true"
	}
}

// PostStory posts a story on behalf of the business account (Bot API 9.0).
//
// Bot API: postStory.
func (b *Bot) PostStory(connectionID string, story StoryParams) (*Story, error) {
	if story.Content == nil {
		return nil, errors.New("telebot: story has no content")
	}
	params := map[string]string{"business_connection_id": connectionID}
	files := make(map[string]File)
	if err := story.embed(params, files); err != nil {
		return nil, err
	}
	embedStoryFlags(params, story)
	return b.story("postStory", files, params)
}

// EditStory edits a story the bot posted on behalf of the business account
// (Bot API 9.0).
//
// Bot API: editStory.
func (b *Bot) EditStory(connectionID string, storyID int, story StoryParams) (*Story, error) {
	if story.Content == nil {
		return nil, errors.New("telebot: story has no content")
	}
	params := map[string]string{
		"business_connection_id": connectionID,
		"story_id":               strconv.Itoa(storyID),
	}
	files := make(map[string]File)
	if err := story.embed(params, files); err != nil {
		return nil, err
	}
	return b.story("editStory", files, params)
}

// DeleteStory deletes a story the bot posted on behalf of the business
// account (Bot API 9.0).
//
// Bot API: deleteStory.
func (b *Bot) DeleteStory(connectionID string, storyID int) error {
	_, err := b.Raw("deleteStory", map[string]string{
		"business_connection_id": connectionID,
		"story_id":               strconv.Itoa(storyID),
	})
	return err
}

// RepostStory reposts from, a story of another business account managed by
// the bot, on behalf of the business account (Bot API 9.3).
//
// Bot API: repostStory.
func (b *Bot) RepostStory(connectionID string, from *Story, story StoryParams) (*Story, error) {
	if from == nil || from.Poster == nil {
		return nil, errors.New("telebot: story to repost has no chat")
	}
	params := map[string]string{
		"business_connection_id": connectionID,
		"from_chat_id":           from.Poster.Recipient(),
		"from_story_id":          strconv.Itoa(from.ID),
	}
	embedStoryFlags(params, story)
	return b.story("repostStory", nil, params)
}
