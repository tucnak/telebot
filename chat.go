package telebot

import (
	"encoding/json"
	"strconv"
	"time"
)

// User object represents a Telegram user, bot.
//
// Bot API: User.
type User struct {
	ID int64 `json:"id"`

	FirstName           string   `json:"first_name"`
	LastName            string   `json:"last_name"`
	Username            string   `json:"username"`
	LanguageCode        string   `json:"language_code"`
	IsBot               bool     `json:"is_bot"`
	IsPremium           bool     `json:"is_premium"`
	AddedToMenu         bool     `json:"added_to_attachment_menu"`
	Usernames           []string `json:"active_usernames"`
	StatusCustomEmojiID string   `json:"emoji_status_custom_emoji_id"`

	// Returns only in getMe
	CanJoinGroups        bool `json:"can_join_groups"`
	CanReadMessages      bool `json:"can_read_all_group_messages"`
	SupportsInline       bool `json:"supports_inline_queries"`
	SupportsGuest        bool `json:"supports_guest_queries"`
	CanConnectToBusiness bool `json:"can_connect_to_business"`
	HasMainWebApp        bool `json:"has_main_web_app"`

	// Bot API 9.4: (Optional) True, if the bot allows users to create
	// and delete topics in private chats.
	AllowsUsersToCreateTopics bool `json:"allows_users_to_create_topics,omitempty"`

	// Bot API 9.3: (Optional) True, if the bot has forum topic mode enabled
	// in private chats.
	HasTopicsEnabled bool `json:"has_topics_enabled,omitempty"`

	// Bot API 9.6: (Optional) True, if the bot can create and control other
	// bots.
	CanManageBots bool `json:"can_manage_bots,omitempty"`

	// Bot API 10.1: (Optional) True, if the bot can be assigned to process
	// join request queries.
	SupportsJoinRequestQueries bool `json:"supports_join_request_queries,omitempty"`
}

// Recipient returns user ID (see Recipient interface).
func (u *User) Recipient() string {
	return strconv.FormatInt(u.ID, 10)
}

// Chat object represents a Telegram user, bot, group or a channel.
//
// Bot API: Chat, ChatFullInfo.
type Chat struct {
	ID int64 `json:"id"`

	// See ChatType and consts.
	Type ChatType `json:"type"`

	// Won't be there for ChatPrivate.
	Title string `json:"title"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`

	// Returns only in getChat
	Bio                   string        `json:"bio,omitempty"`
	Photo                 *ChatPhoto    `json:"photo,omitempty"`
	Description           string        `json:"description,omitempty"`
	InviteLink            string        `json:"invite_link,omitempty"`
	PinnedMessage         *Message      `json:"pinned_message,omitempty"`
	Permissions           *Rights       `json:"permissions,omitempty"`
	Reactions             []Reaction    `json:"available_reactions"`
	SlowMode              int           `json:"slow_mode_delay,omitempty"`
	StickerSet            string        `json:"sticker_set_name,omitempty"`
	CanSetStickerSet      bool          `json:"can_set_sticker_set,omitempty"`
	CanSendPaidMedia      bool          `json:"can_send_paid_media"`
	CustomEmojiStickerSet string        `json:"custom_emoji_sticker_set_name"`
	LinkedChatID          int64         `json:"linked_chat_id,omitempty"`
	ChatLocation          *ChatLocation `json:"location,omitempty"`
	Private               bool          `json:"has_private_forwards,omitempty"`
	Protected             bool          `json:"has_protected_content,omitempty"`
	NoVoiceAndVideo       bool          `json:"has_restricted_voice_and_video_messages"`
	HasHiddenMembers      bool          `json:"has_hidden_members,omitempty"`
	AggressiveAntiSpam    bool          `json:"has_aggressive_anti_spam_enabled,omitempty"`
	IsForum               bool          `json:"is_forum"`

	// Bot API 9.0: the gifts the user or chat accepts; returned only by
	// ChatByID.
	AcceptedGiftTypes *AcceptedGiftTypes `json:"accepted_gift_types,omitempty"`

	// Returned only by ChatByID: all active usernames (6.0), whether users
	// must join the supergroup to write (6.4) or be approved to join (6.4),
	// the user's rating (9.3), the gift-based color scheme of the chat
	// (9.3), and the Stars a user pays per message (9.0).
	ActiveUsernames      []string          `json:"active_usernames,omitempty"`
	JoinToSendMessages   bool              `json:"join_to_send_messages,omitempty"`
	JoinByRequest        bool              `json:"join_by_request,omitempty"`
	Rating               *UserRating       `json:"rating,omitempty"`
	UniqueGiftColors     *UniqueGiftColors `json:"unique_gift_colors,omitempty"`
	PaidMessageStarCount int               `json:"paid_message_star_count,omitempty"`

	// Bot API 9.2: the supergroup is the direct messages chat of a channel,
	// and ParentChat (returned only by ChatByID) is that channel.
	IsDirectMessages bool  `json:"is_direct_messages,omitempty"`
	ParentChat       *Chat `json:"parent_chat,omitempty"`

	StatusCustomEmojiID            string               `json:"emoji_status_custom_emoji_id"`
	EmojiExpirationUnixtime        int64                `json:"emoji_status_expiration_date"`
	BackgroundCustomEmojiID        string               `json:"background_custom_emoji_id"`
	AccentColorID                  int                  `json:"accent_color_id"`
	ProfileAccentColorID           int                  `json:"profile_accent_color_id"`
	ProfileBackgroundCustomEmojiID string               `json:"profile_background_custom_emoji_id"`
	HasVisibleHistory              bool                 `json:"has_visible_history"`
	UnrestrictBoosts               int                  `json:"unrestrict_boost_count"`
	MaxReactions                   int                  `json:"max_reaction_count"`
	Birthdate                      Birthdate            `json:"birthdate,omitempty"`
	PersonalChat                   *Chat                `json:"personal_chat,omitempty"`
	BusinessIntro                  BusinessIntro        `json:"business_intro,omitempty"`
	BusinessLocation               BusinessLocation     `json:"business_location,omitempty"`
	BusinessOpeningHours           BusinessOpeningHours `json:"business_opening_hours,omitempty"`

	// (Optional) The community the chat belongs to; returned only by
	// ChatByID (Bot API 10.2).
	Community *Community `json:"community,omitempty"`

	// (Optional) The TON wallet card of a private chat's user; returned only
	// by ChatByID (Bot API 10.4).
	TonWalletCard *TonWalletCard `json:"ton_wallet_card,omitempty"`

	// (Optional) The theme of the screen for sharing a link to the chat;
	// returned only by ChatByID (Bot API 10.4).
	ShareTheme *ChatShareTheme `json:"share_theme,omitempty"`

	// (Optional) The bot that processes join request queries in the chat;
	// returned only by ChatByID, to administrators (Bot API 10.1).
	GuardBot *User `json:"guard_bot,omitempty"`
}

// Recipient returns chat ID (see Recipient interface).
func (c *Chat) Recipient() string {
	return strconv.FormatInt(c.ID, 10)
}

// ChatType represents one of the possible chat types.
type ChatType string

const (
	ChatPrivate        ChatType = "private"
	ChatGroup          ChatType = "group"
	ChatSuperGroup     ChatType = "supergroup"
	ChatChannel        ChatType = "channel"
	ChatChannelPrivate ChatType = "privatechannel"
)

// ChatLocation represents a location to which a chat is connected.
//
// Bot API: ChatLocation.
type ChatLocation struct {
	Location Location `json:"location,omitempty"`
	Address  string   `json:"address,omitempty"`
}

// ChatPhoto object represents a chat photo.
//
// Bot API: ChatPhoto.
type ChatPhoto struct {
	// File identifiers of small (160x160) chat photo
	SmallFileID   string `json:"small_file_id"`
	SmallUniqueID string `json:"small_file_unique_id"`

	// File identifiers of big (640x640) chat photo
	BigFileID   string `json:"big_file_id"`
	BigUniqueID string `json:"big_file_unique_id"`
}

// ChatMember object represents information about a single chat member.
//
// Bot API: ChatMember, ChatMemberOwner, ChatMemberAdministrator, ChatMemberMember, ChatMemberRestricted, ChatMemberLeft, ChatMemberBanned.
type ChatMember struct {
	Rights

	User      *User        `json:"user"`
	Role      MemberStatus `json:"status"`
	Title     string       `json:"custom_title"`
	Anonymous bool         `json:"is_anonymous"`
	Member    bool         `json:"is_member,omitempty"`

	// Date when restrictions will be lifted for the user, unix time.
	//
	// If user is restricted for more than 366 days or less than
	// 30 seconds from the current time, they are considered to be
	// restricted forever.
	//
	// Use tele.Forever().
	//
	RestrictedUntil int64 `json:"until_date,omitempty"`

	JoinToSend    string `json:"join_to_send_messages"`
	JoinByRequest string `json:"join_by_request"`

	// Bot API 9.5: (Optional) Tag of the member; for regular and restricted members only.
	Tag string `json:"tag,omitempty"`
}

// MemberStatus is one's chat status.
type MemberStatus string

const (
	Creator       MemberStatus = "creator"
	Administrator MemberStatus = "administrator"
	Member        MemberStatus = "member"
	Restricted    MemberStatus = "restricted"
	Left          MemberStatus = "left"
	Kicked        MemberStatus = "kicked"
)

// ChatMemberUpdate object represents changes in the status of a chat member.
//
// Bot API: ChatMemberUpdated.
type ChatMemberUpdate struct {
	// Chat where the user belongs to.
	Chat *Chat `json:"chat"`

	// Sender which user the action was triggered.
	Sender *User `json:"from"`

	// Unixtime, use Date() to get time.Time.
	Unixtime int64 `json:"date"`

	// Previous information about the chat member.
	OldChatMember *ChatMember `json:"old_chat_member"`

	// New information about the chat member.
	NewChatMember *ChatMember `json:"new_chat_member"`

	// (Optional) InviteLink which was used by the user to
	// join the chat; for joining by invite link events only.
	InviteLink *ChatInviteLink `json:"invite_link"`

	// (Optional) True, if the user joined the chat after sending
	// a direct join request without using an invite link and being
	// approved by an administrator
	ViaJoinRequest bool `json:"via_join_request"`

	// (Optional) True, if the user joined the chat via a chat folder invite link.
	ViaFolderLink bool `json:"via_chat_folder_invite_link"`
}

// Time returns the moment of the change in local time.
func (c *ChatMemberUpdate) Time() time.Time {
	return time.Unix(c.Unixtime, 0)
}

// ChatID represents a chat or an user integer ID, which can be used
// as recipient in bot methods. It is very useful in cases where
// you have special group IDs, for example in your config, and don't
// want to wrap it into *tele.Chat every time you send messages.
//
// Example:
//
//	group := tele.ChatID(-100756389456)
//	b.Send(group, "Hello!")
//
//	type Config struct {
//		AdminGroup tele.ChatID `json:"admin_group"`
//	}
//	b.Send(conf.AdminGroup, "Hello!")
type ChatID int64

// Recipient returns chat ID (see Recipient interface).
func (i ChatID) Recipient() string {
	return strconv.FormatInt(int64(i), 10)
}

// ChatJoinRequest represents a join request sent to a chat.
//
// Bot API: ChatJoinRequest.
type ChatJoinRequest struct {
	// Chat to which the request was sent.
	Chat *Chat `json:"chat"`

	// Sender is the user that sent the join request.
	Sender *User `json:"from"`

	// UserChatID is an ID of a private chat with the user
	// who sent the join request. The bot can use this ID
	// for 5 minutes to send messages until the join request
	// is processed, assuming no other administrator contacted the user.
	UserChatID int64 `json:"user_chat_id"`

	// Unixtime, use ChatJoinRequest.Time() to get time.Time.
	Unixtime int64 `json:"date"`

	// Bio of the user, optional.
	Bio string `json:"bio"`

	// InviteLink is the chat invite link that was used by
	//the user to send the join request, optional.
	InviteLink *ChatInviteLink `json:"invite_link"`

	// (Optional) Identifier of the join request query, for a bot assigned to
	// process join requests. The bot must answer it with
	// AnswerJoinRequest, directly or after SendJoinRequestWebApp (Bot API 10.1).
	QueryID string `json:"query_id,omitempty"`
}

// ChatInviteLink object represents an invite for a chat.
//
// Bot API: ChatInviteLink.
type ChatInviteLink struct {
	// The invite link.
	InviteLink string `json:"invite_link"`

	// Invite link name.
	Name string `json:"name"`

	// The creator of the link.
	Creator *User `json:"creator"`

	// If the link is primary.
	IsPrimary bool `json:"is_primary"`

	// If the link is revoked.
	IsRevoked bool `json:"is_revoked"`

	// (Optional) Point in time when the link will expire,
	// use ExpireDate() to get time.Time.
	ExpireUnixtime int64 `json:"expire_date,omitempty"`

	// (Optional) Maximum number of users that can be members of
	// the chat simultaneously.
	MemberLimit int `json:"member_limit,omitempty"`

	// (Optional) True, if users joining the chat via the link need to
	// be approved by chat administrators. If True, member_limit can't be specified.
	JoinRequest bool `json:"creates_join_request"`

	// (Optional) Number of pending join requests created using this link.
	PendingCount int `json:"pending_join_request_count"`

	// Bot API 7.9: (Optional) For subscription links, the seconds between
	// payments and the Telegram Stars paid each period.
	SubscriptionPeriod int `json:"subscription_period,omitempty"`
	SubscriptionPrice  int `json:"subscription_price,omitempty"`
}

// Story is a story posted to a chat.
//
// Bot API: Story.
type Story struct {
	// Unique identifier for the story in the chat
	ID int `json:"id"`

	// Chat that posted the story
	Poster *Chat `json:"chat"`
}

// Birthdate is the birthdate of a user.
//
// Bot API: Birthdate.
type Birthdate struct {
	// Day of the user's birth; 1-31
	Day int `json:"day"`

	// Month of the user's birth; 1-12
	Month int `json:"month"`

	// (Optional) Year of the user's birth
	Year int `json:"year"`
}

// ChatFullInfo object contains full information about a chat.
//
// Bot API: ChatFullInfo.
type ChatFullInfo struct {
	ID int64 `json:"id"`

	// See ChatType and consts.
	Type ChatType `json:"type"`

	// Won't be there for ChatPrivate.
	Title string `json:"title"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`

	// (Optional) Bio of the chat
	Bio string `json:"bio,omitempty"`

	// (Optional) Chat photo
	Photo *ChatPhoto `json:"photo,omitempty"`

	// (Optional) Description of the chat
	Description string `json:"description,omitempty"`

	// (Optional) Invite link for the chat
	InviteLink string `json:"invite_link,omitempty"`

	// (Optional) Pinned message
	PinnedMessage *Message `json:"pinned_message,omitempty"`

	// (Optional) Default chat member permissions
	Permissions *Rights `json:"permissions,omitempty"`

	// (Optional) The minimum allowed delay between consecutive messages sent by each user
	SlowMode int `json:"slow_mode_delay,omitempty"`

	// (Optional) The time after which all messages sent to the chat will be automatically deleted
	MessageAutoDeleteTime int `json:"message_auto_delete_time,omitempty"`

	// (Optional) True, if new members of the chat can be anonymous
	HasAnonymousMembers bool `json:"has_anonymous_members,omitempty"`

	// (Optional) True, if reactions are available in the chat
	HasAggressiveAntiSpam bool `json:"has_aggressive_anti_spam_enabled,omitempty"`

	// (Optional) True, if anti-spam checks are enabled in the chat
	HasHiddenMembers bool `json:"has_hidden_members,omitempty"`

	// (Optional) True, if the chat has protected content
	HasProtectedContent bool `json:"has_protected_content,omitempty"`

	// (Optional) True, if non-administrators can only get the list of bots and administrators in the chat
	HasVisibleHistory bool `json:"has_visible_history,omitempty"`

	// (Optional) True, if messages from the chat can't be forwarded to other chats
	NoMessagesToOther bool `json:"has_restricted_voice_and_video_messages,omitempty"`

	// (Optional) For supergroups, name of group sticker set
	StickerSetName string `json:"sticker_set_name,omitempty"`

	// (Optional) True, if the bot can change the group sticker set
	CanSetStickerSet bool `json:"can_set_sticker_set,omitempty"`

	// (Optional) Custom emoji sticker set name
	CustomEmojiStickerSet string `json:"custom_emoji_sticker_set_name,omitempty"`

	// (Optional) Chat location
	Location *ChatLocation `json:"location,omitempty"`

	// (Optional) Identifier of the linked chat
	LinkedChatID int64 `json:"linked_chat_id,omitempty"`

	// Bot API 8.3: True, if the chat can be gifted with Telegram Stars gifts
	CanSendGift bool `json:"can_send_gift"`

	// Bot API 9.4: For private chats, the first audio added to the profile of the user.
	FirstProfileAudio *Audio `json:"first_profile_audio,omitempty"`

	// Bot API 9.0: the gifts the user or chat accepts.
	AcceptedGiftTypes *AcceptedGiftTypes `json:"accepted_gift_types,omitempty"`

	// Bot API 9.2: the supergroup is the direct messages chat of a channel,
	// ParentChat.
	IsDirectMessages bool  `json:"is_direct_messages,omitempty"`
	ParentChat       *Chat `json:"parent_chat,omitempty"`

	// Bot API 10.2: (Optional) The community the chat belongs to.
	Community *Community `json:"community,omitempty"`

	// Bot API 10.4: (Optional) The TON wallet card of a private chat's user.
	TonWalletCard *TonWalletCard `json:"ton_wallet_card,omitempty"`

	// Bot API 10.4: (Optional) The theme of the screen for sharing a link to
	// the chat.
	ShareTheme *ChatShareTheme `json:"share_theme,omitempty"`

	// Bot API 10.1: (Optional) The bot that processes join request queries in
	// the chat; available to administrators only.
	GuardBot *User `json:"guard_bot,omitempty"`
}

// ChatBackground is a chat background.
//
// Bot API: ChatBackground.
type ChatBackground struct {
	// Type of the background
	Type BackgroundType `json:"type"`
}

// BackgroundType is the type of a chat background: a fill, a wallpaper,
// a pattern or a chat theme.
//
// Bot API: BackgroundType, BackgroundTypeFill, BackgroundTypeWallpaper, BackgroundTypePattern, BackgroundTypeChatTheme.
type BackgroundType struct {
	// Type of the background, always “fill”
	Type string `json:"type"`

	// The background fill
	Fill BackgroundFill `json:"fill,omitempty"`

	// Document with the wallpaper
	Document Document `json:"document,omitempty"`

	// Dimming of the background in dark themes, as a percentage; 0-100
	DarkThemeDimming int `json:"dark_theme_dimming,omitempty"`

	// Intensity of the pattern when it is shown above the filled background; 0-100
	Intensity int `json:"intensity,omitempty"`

	// (Optional) True, if the wallpaper is downscaled to fit in a 450x450
	// square and then box-blurred with radius 12
	Blurred bool `json:"is_blurred,omitempty"`

	// (Optional) True, if the background moves slightly when the device is tilted
	Moving bool `json:"is_moving,omitempty"`

	// (Optional) True, if the background fill must be applied only to the pattern itself.
	// All other pixels are black in this case. For dark themes only
	Inverted bool `json:"is_inverted,omitempty"`

	// Name of the chat theme, which is usually an emoji
	ThemeName string `json:"theme_name,omitempty"`
}

// BackgroundFill is the way a background is filled: a solid color, a
// gradient or a freeform gradient.
//
// Bot API: BackgroundFill, BackgroundFillSolid, BackgroundFillGradient, BackgroundFillFreeformGradient.
type BackgroundFill struct {
	// Type of the background fill.
	Type string `json:"type"`

	// The color of the background fill in the RGB24 format
	SolidColor int `json:"color,omitempty"`

	// Top color of the gradient in the RGB24 format
	GradientTopColor int `json:"top_color,omitempty"`

	// Bottom color of the gradient in the RGB24 format
	GradientBottomColor int `json:"bottom_color,omitempty"`

	// Clockwise rotation angle of the background fill in degrees; 0-359
	GradientRotationAngle int `json:"rotation_angle,omitempty"`

	// A list of the 3 or 4 base colors that are used to generate
	// the freeform gradient in the RGB24 format
	GradientColors []int `json:"colors,omitempty"`
}

// ExpireDate returns the moment of the link expiration in local time.
func (c *ChatInviteLink) ExpireDate() time.Time {
	return time.Unix(c.ExpireUnixtime, 0)
}

// Time returns the moment of chat join request sending in local time.
func (r ChatJoinRequest) Time() time.Time {
	return time.Unix(r.Unixtime, 0)
}

// Time returns the moment of the emoji status expiration.
func (c *Chat) Time() time.Time {
	return time.Unix(c.EmojiExpirationUnixtime, 0)
}

// InviteLink should be used to export chat's invite link.
//
// Bot API: exportChatInviteLink.
func (b *Bot) InviteLink(chat *Chat) (string, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("exportChatInviteLink", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Result string
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", wrapError(err)
	}
	return resp.Result, nil
}

// CreateInviteLink creates an additional invite link for a chat.
//
// Bot API: createChatInviteLink.
func (b *Bot) CreateInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}
	if link != nil {
		params["name"] = link.Name

		if link.ExpireUnixtime != 0 {
			params["expire_date"] = strconv.FormatInt(link.ExpireUnixtime, 10)
		}
		if link.MemberLimit > 0 {
			params["member_limit"] = strconv.Itoa(link.MemberLimit)
		} else if link.JoinRequest {
			params["creates_join_request"] = "true"
		}
	}

	data, err := b.Raw("createChatInviteLink", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result ChatInviteLink `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}

	return &resp.Result, nil
}

// EditInviteLink edits a non-primary invite link created by the bot.
//
// Bot API: editChatInviteLink.
func (b *Bot) EditInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}
	if link != nil {
		params["invite_link"] = link.InviteLink
		params["name"] = link.Name

		if link.ExpireUnixtime != 0 {
			params["expire_date"] = strconv.FormatInt(link.ExpireUnixtime, 10)
		}
		if link.MemberLimit > 0 {
			params["member_limit"] = strconv.Itoa(link.MemberLimit)
		} else if link.JoinRequest {
			params["creates_join_request"] = "true"
		}
	}

	data, err := b.Raw("editChatInviteLink", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result ChatInviteLink `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}

	return &resp.Result, nil
}

// RevokeInviteLink revokes an invite link created by the bot.
//
// Bot API: revokeChatInviteLink.
func (b *Bot) RevokeInviteLink(chat Recipient, link string) (*ChatInviteLink, error) {
	params := map[string]string{
		"chat_id":     chat.Recipient(),
		"invite_link": link,
	}

	data, err := b.Raw("revokeChatInviteLink", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result ChatInviteLink `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}

	return &resp.Result, nil
}

// Results of AnswerJoinRequest.
const (
	JoinRequestApprove = "approve"
	JoinRequestDecline = "decline"
	JoinRequestQueue   = "queue"
)

// AnswerJoinRequest resolves the join request query a bot assigned to process
// join requests received. result is JoinRequestApprove, JoinRequestDecline,
// or JoinRequestQueue to leave the decision to other administrators
// (Bot API 10.1).
//
// Bot API: answerChatJoinRequestQuery.
func (b *Bot) AnswerJoinRequest(req *ChatJoinRequest, result string) error {
	_, err := b.Raw("answerChatJoinRequestQuery", map[string]string{
		"chat_join_request_query_id": req.QueryID,
		"result":                     result,
	})
	return err
}

// SendJoinRequestWebApp shows the user who sent the join request query a
// Mini App at url before the bot decides; resolve the query afterwards with
// AnswerJoinRequest (Bot API 10.1).
//
// Bot API: sendChatJoinRequestWebApp.
func (b *Bot) SendJoinRequestWebApp(req *ChatJoinRequest, url string) error {
	_, err := b.Raw("sendChatJoinRequestWebApp", map[string]string{
		"chat_join_request_query_id": req.QueryID,
		"web_app_url":                url,
	})
	return err
}

// CreateSubscriptionInviteLink creates a subscription invite link for a
// channel, with link's Name, SubscriptionPeriod (currently always 2592000
// seconds) and SubscriptionPrice in Telegram Stars (Bot API 7.9).
//
// Bot API: createChatSubscriptionInviteLink.
func (b *Bot) CreateSubscriptionInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error) {
	params := map[string]string{
		"chat_id":             chat.Recipient(),
		"subscription_period": strconv.Itoa(link.SubscriptionPeriod),
		"subscription_price":  strconv.Itoa(link.SubscriptionPrice),
	}
	if link.Name != "" {
		params["name"] = link.Name
	}
	return b.inviteLink("createChatSubscriptionInviteLink", params)
}

// EditSubscriptionInviteLink renames a subscription invite link created by
// the bot; an empty Name removes its name (Bot API 7.9).
//
// Bot API: editChatSubscriptionInviteLink.
func (b *Bot) EditSubscriptionInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error) {
	params := map[string]string{
		"chat_id":     chat.Recipient(),
		"invite_link": link.InviteLink,
		"name":        link.Name,
	}
	return b.inviteLink("editChatSubscriptionInviteLink", params)
}

func (b *Bot) inviteLink(method string, params map[string]string) (*ChatInviteLink, error) {
	data, err := b.Raw(method, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *ChatInviteLink
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// ApproveJoinRequest approves a chat join request.
//
// Bot API: approveChatJoinRequest.
func (b *Bot) ApproveJoinRequest(chat Recipient, user *User) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"user_id": user.Recipient(),
	}

	data, err := b.Raw("approveChatJoinRequest", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// DeclineJoinRequest declines a chat join request.
//
// Bot API: declineChatJoinRequest.
func (b *Bot) DeclineJoinRequest(chat Recipient, user *User) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"user_id": user.Recipient(),
	}

	data, err := b.Raw("declineChatJoinRequest", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// SetGroupTitle should be used to update group title.
//
// Bot API: setChatTitle.
func (b *Bot) SetGroupTitle(chat *Chat, title string) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"title":   title,
	}

	_, err := b.Raw("setChatTitle", params)
	return err
}

// SetGroupDescription should be used to update group description.
//
// Bot API: setChatDescription.
func (b *Bot) SetGroupDescription(chat *Chat, description string) error {
	params := map[string]string{
		"chat_id":     chat.Recipient(),
		"description": description,
	}

	_, err := b.Raw("setChatDescription", params)
	return err
}

// SetGroupPhoto should be used to update group photo.
//
// Bot API: setChatPhoto.
func (b *Bot) SetGroupPhoto(chat *Chat, p *Photo) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.sendFiles("setChatPhoto", map[string]File{"photo": p.File}, params)
	return err
}

// SetGroupStickerSet should be used to update group's group sticker set.
//
// Bot API: setChatStickerSet.
func (b *Bot) SetGroupStickerSet(chat *Chat, setName string) error {
	params := map[string]string{
		"chat_id":          chat.Recipient(),
		"sticker_set_name": setName,
	}

	_, err := b.Raw("setChatStickerSet", params)
	return err
}

// SetGroupPermissions sets default chat permissions for all members.
//
// Bot API: setChatPermissions.
func (b *Bot) SetGroupPermissions(chat *Chat, perms Rights) error {
	params := map[string]interface{}{
		"chat_id":     chat.Recipient(),
		"permissions": perms,
	}
	if perms.Independent {
		params["use_independent_chat_permissions"] = true
	}

	_, err := b.Raw("setChatPermissions", params)
	return err
}

// DeleteGroupPhoto should be used to just remove group photo.
//
// Bot API: deleteChatPhoto.
func (b *Bot) DeleteGroupPhoto(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("deleteChatPhoto", params)
	return err
}

// DeleteGroupStickerSet should be used to just remove group sticker set.
//
// Bot API: deleteChatStickerSet.
func (b *Bot) DeleteGroupStickerSet(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("deleteChatStickerSet", params)
	return err
}

// UserRating is the rating of a user based on their Telegram Star spending
// (Bot API 9.3).
//
// Bot API: UserRating.
type UserRating struct {
	// Current level; higher is more trustworthy, negative is unreliable.
	Level int `json:"level"`

	// The rating, and the rating needed for the current and the next level;
	// NextLevelRating is zero at the maximum level.
	Rating             int `json:"rating"`
	CurrentLevelRating int `json:"current_level_rating"`
	NextLevelRating    int `json:"next_level_rating,omitempty"`
}
