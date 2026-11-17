package telebot

import (
	"strconv"
	"time"
	"unicode/utf16"
)

// Message object represents a message.
//
// Bot API: Message, MaybeInaccessibleMessage, InaccessibleMessage, MessageId.
type Message struct {
	ID int `json:"message_id"`

	// (Optional) Unique identifier of a message thread to which the message belongs; for supergroups only
	ThreadID int `json:"message_thread_id"`

	// For message sent to channels, Sender will be nil
	Sender *User `json:"from"`

	// Unixtime, use Message.Time() to get time.Time
	Unixtime int64 `json:"date"`

	// Conversation the message belongs to.
	Chat *Chat `json:"chat"`

	// Sender of the message, sent on behalf of a chat.
	SenderChat *Chat `json:"sender_chat"`

	// (Optional) For ephemeral messages, the only user who can see the
	// message besides the bot (Bot API 10.2).
	ReceiverUser *User `json:"receiver_user,omitempty"`

	// (Optional) For ephemeral messages, the identifier of the message, used
	// instead of ID by the EditEphemeral* methods and DeleteEphemeral.
	EphemeralID int `json:"ephemeral_message_id,omitempty"`

	// The user that called the bot in guest mode.
	GuestUser *User `json:"guest_bot_caller_user"`

	// The chat that called the bot in guest mode.
	GuestChat *Chat `json:"guest_bot_caller_chat"`

	// For forwarded messages, sender of the original message.
	OriginalSender *User `json:"forward_from"`

	// For forwarded messages, chat of the original message when
	// forwarded from a channel.
	OriginalChat *Chat `json:"forward_from_chat"`

	// For forwarded messages, identifier of the original message
	// when forwarded from a channel.
	OriginalMessageID int `json:"forward_from_message_id"`

	// For forwarded messages, signature of the post author.
	OriginalSignature string `json:"forward_signature"`

	// For forwarded messages, sender's name from users who
	// disallow adding a link to their account.
	OriginalSenderName string `json:"forward_sender_name"`

	// For forwarded messages, unixtime of the original message.
	OriginalUnixtime int `json:"forward_date"`

	// For information about the original message for forwarded messages.
	Origin *MessageOrigin `json:"forward_origin"`

	// Message is a channel post that was automatically forwarded to the connected discussion group.
	AutomaticForward bool `json:"is_automatic_forward"`

	// For replies, ReplyTo represents the original message.
	//
	// Note that the Message object in this field will not
	// contain further ReplyTo fields even if it
	// itself is a reply.
	ReplyTo *Message `json:"reply_to_message"`

	// (Optional) For replies to a story, the original story
	Story *Story `json:"story"`

	// (Optional) Information about the message that is being replied to,
	// which may come from another chat or forum topic.
	ExternalReply *ExternalReply `json:"external_reply"`

	// (Optional) For replies that quote part of the original message,
	// the quoted part of the message.
	Quote *TextQuote `json:"quote"`

	// Shows through which bot the message was sent.
	Via *User `json:"via_bot"`

	// For replies to a story, the original story.
	ReplyToStory *Story `json:"reply_to_story"`

	// (Optional) Persistent identifier of the poll option being replied to
	// (Bot API 9.6).
	ReplyToPollOptionID string `json:"reply_to_poll_option_id,omitempty"`

	// (Optional) Identifier of the checklist task being replied to
	// (Bot API 9.2).
	ReplyToChecklistTaskID int `json:"reply_to_checklist_task_id,omitempty"`

	// (Optional) Time of last edit in Unix.
	LastEdit int64 `json:"edit_date"`

	// (Optional) True, if the message is sent to a forum topic.
	TopicMessage bool `json:"is_topic_message"`

	// (Optional) Message can't be forwarded.
	Protected bool `json:"has_protected_content,omitempty"`

	// (Optional) True, if the message was sent by an implicit action,
	// for example, as an away or a greeting business message, or as a scheduled message
	FromOffline bool `json:"is_from_offline,omitempty"`

	// AlbumID is the unique identifier of a media message group
	// this message belongs to.
	AlbumID string `json:"media_group_id"`

	// Author signature (in channels).
	Signature string `json:"author_signature"`

	// For a text message, the actual UTF-8 text of the message.
	Text string `json:"text"`

	// For registered commands, will contain the string payload:
	//
	// Ex: `/command <payload>` or `/command@botname <payload>`
	Payload string `json:"-"`

	// For text messages, special entities like usernames, URLs, bot commands,
	// etc. that appear in the text.
	Entities Entities `json:"entities,omitempty"`

	// (Optional) PreviewOptions used for link preview generation for the message,
	// if it is a text message and link preview options were changed.
	PreviewOptions *PreviewOptions `json:"link_preview_options,omitempty"`

	// (Optional) Unique identifier of the message effect added to the message
	EffectID string `json:"effect_id"`

	// Some messages containing media, may as well have a caption.
	Caption string `json:"caption,omitempty"`

	// For messages with a caption, special entities like usernames, URLs,
	// bot commands, etc. that appear in the caption.
	CaptionEntities Entities `json:"caption_entities,omitempty"`

	// For an audio recording, information about it.
	Audio *Audio `json:"audio"`

	// For a general file, information about it.
	Document *Document `json:"document"`

	// Message contains paid media; information about the paid media
	PaidMedia PaidMedias `json:"paid_media"`

	// For a photo, all available sizes (thumbnails).
	Photo *Photo `json:"photo"`

	// For a sticker, information about it.
	Sticker *Sticker `json:"sticker"`

	// For a voice message, information about it.
	Voice *Voice `json:"voice"`

	// For a video note, information about it.
	VideoNote *VideoNote `json:"video_note"`

	// For a video, information about it.
	Video *Video `json:"video"`

	// (Optional) Message is a live photo (Bot API 10.0).
	LivePhoto *LivePhoto `json:"live_photo,omitempty"`

	// (Optional) Message is a checklist (Bot API 9.1).
	Checklist *Checklist `json:"checklist,omitempty"`

	// For an animation, information about it.
	Animation *Animation `json:"animation"`

	// For a contact, contact information itself.
	Contact *Contact `json:"contact"`

	// For a location, its longitude and latitude.
	Location *Location `json:"location"`

	// For a venue, information about it.
	Venue *Venue `json:"venue"`

	// For a poll, information the native poll.
	Poll *Poll `json:"poll"`

	// For a game, information about it.
	Game *Game `json:"game"`

	// For a dice, information about it.
	Dice *Dice `json:"dice"`

	// (Optional) The message is a scheduled giveaway message.
	Giveaway *Giveaway `json:"giveaway"`

	// (Optional) A giveaway with public winners was completed.
	GiveawayWinners *GiveawayWinners `json:"giveaway_winners"`

	// (Optional) Service message: a scheduled giveaway was created.
	GiveawayCreated *GiveawayCreated `json:"giveaway_created"`

	// (Optional) Service message: a giveaway without public winners was completed.
	GiveawayCompleted *GiveawayCompleted `json:"giveaway_completed"`

	// (Optional) Unique identifier of the business connection from which the message
	// was received. If non-empty, the message belongs to a chat of the corresponding
	// business account that is independent from any potential bot chat which might
	// share the same identifier.
	BusinessConnectionID string `json:"business_connection_id"`

	// Unique identifier for answering the guest query.
	GuestQueryID string `json:"guest_query_id"`

	// (Optional) For a rich message, the rich formatted content (Bot API 10.1).
	RichMessage *RichMessage `json:"rich_message,omitempty"`

	// (Optional) Message is an HTML widget (Bot API 10.4).
	Widget *Widget `json:"widget,omitempty"`

	// (Optional) The bot that actually sent the message on behalf of the business account.
	// Available only for outgoing messages sent on behalf of the connected business account.
	BusinessBot *User `json:"sender_business_bot"`

	// For a service message, represents a user,
	// that just got added to chat, this message came from.
	//
	// Sender leads to User, capable of invite.
	//
	// UserJoined might be the Bot itself.
	UserJoined *User `json:"new_chat_member"`

	// For a service message, represents a user,
	// that just left chat, this message came from.
	//
	// If user was kicked, Sender leads to a User,
	// capable of this kick.
	//
	// UserLeft might be the Bot itself.
	UserLeft *User `json:"left_chat_member"`

	// For a service message, represents a new title
	// for chat this message came from.
	//
	// Sender would lead to a User, capable of change.
	NewGroupTitle string `json:"new_chat_title"`

	// For a service message, represents all available
	// thumbnails of the new chat photo.
	//
	// Sender would lead to a User, capable of change.
	NewGroupPhoto *Photo `json:"new_chat_photo"`

	// For a service message, new members that were added to
	// the group or supergroup and information about them
	// (the bot itself may be one of these members).
	UsersJoined []User `json:"new_chat_members"`

	// For a service message, true if chat photo just
	// got removed.
	//
	// Sender would lead to a User, capable of change.
	GroupPhotoDeleted bool `json:"delete_chat_photo"`

	// For a service message, true if group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	GroupCreated bool `json:"group_chat_created"`

	// For a service message, true if supergroup has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	SuperGroupCreated bool `json:"supergroup_chat_created"`

	// For a service message, true if channel has been created.
	//
	// You would receive such a message if you are one of
	// initial channel administrators.
	//
	// Sender would lead to creator of the chat.
	ChannelCreated bool `json:"channel_chat_created"`

	// For a service message, the destination (supergroup) you
	// migrated to.
	//
	// You would receive such a message when your chat has migrated
	// to a supergroup.
	//
	// Sender would lead to creator of the migration.
	MigrateTo int64 `json:"migrate_to_chat_id"`

	// For a service message, the Origin (normal group) you migrated
	// from.
	//
	// You would receive such a message when your chat has migrated
	// to a supergroup.
	//
	// Sender would lead to creator of the migration.
	MigrateFrom int64 `json:"migrate_from_chat_id"`

	// Specified message was pinned. Note that the Message object
	// in this field will not contain further ReplyTo fields even
	// if it is itself a reply.
	PinnedMessage *Message `json:"pinned_message"`

	// Message is an invoice for a payment.
	Invoice *Invoice `json:"invoice"`

	// Message is a service message about a successful payment.
	Payment *Payment `json:"successful_payment"`

	// Message is a service message about a refunded payment, information about the payment.
	RefundedPayment *RefundedPayment `json:"refunded_payment"`

	// Service message: a regular gift was sent or received.
	Gift *GiftInfo `json:"gift,omitempty"`

	// Service message: a unique gift was sent or received.
	UniqueGift *UniqueGiftInfo `json:"unique_gift,omitempty"`

	// Service message: an upgrade was bought for a gift after it was sent
	// (Bot API 9.3).
	GiftUpgradeSent *GiftInfo `json:"gift_upgrade_sent,omitempty"`

	// Service message: a transfer with a TON wallet (Bot API 10.4).
	TonWalletTransfer *TonWalletTransfer `json:"ton_wallet_transfer,omitempty"`

	// (Optional) Topic of a channel's direct messages chat the message
	// belongs to (Bot API 9.2).
	DirectMessagesTopic *DirectMessagesTopic `json:"direct_messages_topic,omitempty"`

	// (Optional) The message is a paid post, which must not be deleted for
	// 24 hours for the payment to arrive (Bot API 9.2).
	IsPaidPost bool `json:"is_paid_post,omitempty"`

	// (Optional) The message is a suggested post (Bot API 9.2).
	SuggestedPostInfo *SuggestedPostInfo `json:"suggested_post_info,omitempty"`

	// Service messages about suggested posts (Bot API 9.2).
	SuggestedPostApproved       *SuggestedPostApprovalInfo   `json:"suggested_post_approved,omitempty"`
	SuggestedPostApprovalFailed *SuggestedPostApprovalFailed `json:"suggested_post_approval_failed,omitempty"`
	SuggestedPostDeclined       *SuggestedPostDeclineInfo    `json:"suggested_post_declined,omitempty"`
	SuggestedPostPaid           *SuggestedPostPaid           `json:"suggested_post_paid,omitempty"`
	SuggestedPostRefunded       *SuggestedPostRefunded       `json:"suggested_post_refunded,omitempty"`

	// (Optional) Telegram Passport data shared with the bot.
	PassportData *PassportData `json:"passport_data,omitempty"`

	// (Optional) Stars the sender paid to send the message (Bot API 9.0).
	PaidStarCount int `json:"paid_star_count,omitempty"`

	// Service message: the price of messages in the chat changed
	// (Bot API 9.0).
	PaidMessagePriceChanged *PaidMessagePriceChanged `json:"paid_message_price_changed,omitempty"`

	// Service message: checklist tasks were marked as done or not done
	// (Bot API 9.1).
	ChecklistTasksDone *ChecklistTasksDone `json:"checklist_tasks_done,omitempty"`

	// Service message: tasks were added to a checklist (Bot API 9.1).
	ChecklistTasksAdded *ChecklistTasksAdded `json:"checklist_tasks_added,omitempty"`

	// Service message: the price of direct messages to the channel changed
	// (Bot API 9.1).
	DirectMessagePriceChanged *DirectMessagePriceChanged `json:"direct_message_price_changed,omitempty"`

	// Service message: an option was added to a poll (Bot API 9.6).
	PollOptionAdded *PollOptionAdded `json:"poll_option_added,omitempty"`

	// Service message: an option was deleted from a poll (Bot API 9.6).
	PollOptionDeleted *PollOptionDeleted `json:"poll_option_deleted,omitempty"`

	// Service message: the user created a bot managed by this one
	// (Bot API 9.6).
	ManagedBotCreated *ManagedBotCreated `json:"managed_bot_created,omitempty"`

	// Service message: the sender used a link that lets them message the
	// user without paying Telegram Stars (Bot API 10.4).
	PaidMessagesWaived *PaidMessagesWaived `json:"paid_messages_waived,omitempty"`

	// For a service message, a user was shared with the bot.
	UserShared *RecipientShared `json:"users_shared,omitempty"`

	// For a service message, a chat was shared with the bot.
	ChatShared *RecipientShared `json:"chat_shared,omitempty"`

	// The domain name of the website on which the user has logged in.
	ConnectedWebsite string `json:"connected_website,omitempty"`

	// For a service message, a video chat started in the chat.
	VideoChatStarted *VideoChatStarted `json:"video_chat_started,omitempty"`

	// For a service message, a video chat ended in the chat.
	VideoChatEnded *VideoChatEnded `json:"video_chat_ended,omitempty"`

	// For a service message, some users were invited in the video chat.
	VideoChatParticipants *VideoChatParticipants `json:"video_chat_participants_invited,omitempty"`

	// For a service message, a video chat schedule in the chat.
	VideoChatScheduled *VideoChatScheduled `json:"video_chat_scheduled,omitempty"`

	// For a data sent by a Web App.
	WebAppData *WebAppData `json:"web_app_data,omitempty"`

	// For a service message, represents the content of a service message,
	// sent whenever a user in the chat triggers a proximity alert set by another user.
	ProximityAlert *ProximityAlert `json:"proximity_alert_triggered,omitempty"`

	// For a service message, represents about a change in auto-delete timer settings.
	AutoDeleteTimer *AutoDeleteTimer `json:"message_auto_delete_timer_changed,omitempty"`

	// Inline keyboard attached to the message.
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`

	// Service message: user boosted the chat.
	BoostAdded *BoostAdded `json:"boost_added"`

	// Service message: chat background set
	ChatBackground ChatBackground `json:"chat_background_set"`

	// If the sender of the message boosted the chat, the number of boosts
	// added by the user.
	SenderBoosts int `json:"sender_boost_count"`

	// Bot API 9.5: (Optional) Tag or custom title of the sender of the message;
	// for supergroups only.
	SenderTag string `json:"sender_tag,omitempty"`

	// Service message: forum topic created
	TopicCreated *Topic `json:"forum_topic_created,omitempty"`

	// Service message: forum topic closed
	TopicClosed *ForumTopicClosed `json:"forum_topic_closed,omitempty"`

	// Service message: forum topic reopened
	TopicReopened *Topic `json:"forum_topic_reopened,omitempty"`

	// Service message: forum topic deleted
	TopicEdited *Topic `json:"forum_topic_edited,omitempty"`

	// Service message: general forum topic hidden
	GeneralTopicHidden *GeneralForumTopicHidden `json:"general_forum_topic_hidden,omitempty"`

	// Service message: general forum topic unhidden
	GeneralTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden,omitempty"`

	// Service message: represents spoiler information about the message.
	HasMediaSpoiler bool `json:"has_media_spoiler,omitempty"`

	// (Optional) Pass True, if the caption must be shown above the message media
	CaptionAbove bool `json:"show_caption_above_media"`

	// Service message: the user allowed the bot added to the attachment menu to write messages
	WriteAccessAllowed *WriteAccessAllowed `json:"write_access_allowed,omitempty"`

	// Bot API 9.4: Service message: chat owner has left.
	ChatOwnerLeft *ChatOwnerLeft `json:"chat_owner_left,omitempty"`

	// Bot API 9.4: Service message: chat owner has changed.
	ChatOwnerChanged *ChatOwnerChanged `json:"chat_owner_changed,omitempty"`

	// Service message: the chat or the bot was added to a community.
	CommunityChatAdded *CommunityChatAdded `json:"community_chat_added,omitempty"`

	// Service message: the chat or the bot was removed from a community.
	CommunityChatRemoved *CommunityChatRemoved `json:"community_chat_removed,omitempty"`

	// Service message: a user joined the chat from a community.
	CommunityChatJoined *CommunityChatJoined `json:"community_chat_joined,omitempty"`
}

// ChatOwnerLeft describes a service message about the chat owner leaving the chat.
//
// Bot API: ChatOwnerLeft.
type ChatOwnerLeft struct {
	// (Optional) The user who will become the new owner of the chat
	// if the previous owner does not return to the chat.
	NewOwner *User `json:"new_owner,omitempty"`
}

// ChatOwnerChanged describes a service message about an ownership change in the chat.
//
// Bot API: ChatOwnerChanged.
type ChatOwnerChanged struct {
	// The new owner of the chat.
	NewOwner *User `json:"new_owner"`
}

// MessageEntity object represents "special" parts of text messages,
// including hashtags, usernames, URLs, etc.
//
// Bot API: MessageEntity.
type MessageEntity struct {
	// Specifies entity type.
	Type EntityType `json:"type"`

	// Offset in UTF-16 code units to the start of the entity.
	Offset int `json:"offset"`

	// Length of the entity in UTF-16 code units.
	Length int `json:"length"`

	// (Optional) For EntityTextLink entity type only.
	//
	// URL will be opened after user taps on the text.
	URL string `json:"url,omitempty"`

	// (Optional) For EntityTMention entity type only.
	User *User `json:"user,omitempty"`

	// (Optional) For EntityCodeBlock entity type only.
	Language string `json:"language,omitempty"`

	// (Optional) For EntityCustomEmoji entity type only.
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`

	// (Optional) For EntityHashtag and EntityCashtag entity types only.
	// Username of the chat where the hashtag or cashtag search should be performed.
	ChatUsername string `json:"chat_username,omitempty"`

	// Bot API 9.5: (Optional) For EntityDateTime only,
	// the Unix time associated with the entity.
	UnixTime int64 `json:"unix_time,omitempty"`

	// Bot API 9.5: (Optional) For EntityDateTime only,
	// the string that defines the formatting of the date and time.
	DateTimeFormat string `json:"date_time_format,omitempty"`
}

// EntityType is a MessageEntity type.
type EntityType string

const (
	EntityMention       EntityType = "mention"
	EntityTMention      EntityType = "text_mention"
	EntityHashtag       EntityType = "hashtag"
	EntityCashtag       EntityType = "cashtag"
	EntityCommand       EntityType = "bot_command"
	EntityURL           EntityType = "url"
	EntityEmail         EntityType = "email"
	EntityPhone         EntityType = "phone_number"
	EntityBold          EntityType = "bold"
	EntityItalic        EntityType = "italic"
	EntityUnderline     EntityType = "underline"
	EntityStrikethrough EntityType = "strikethrough"
	EntityCode          EntityType = "code"
	EntityCodeBlock     EntityType = "pre"
	EntityTextLink      EntityType = "text_link"
	EntitySpoiler       EntityType = "spoiler"
	EntityCustomEmoji   EntityType = "custom_emoji"
	EntityBlockquote    EntityType = "blockquote"
	EntityEBlockquote   EntityType = "expandable_blockquote"
	EntityDateTime      EntityType = "date_time"
	EntityBankCard      EntityType = "bank_card_number" // Bot API 10.4
	EntityTonAddress    EntityType = "ton_address"      // Bot API 10.4
)

// Entities are used to set message's text entities as a send option.
type Entities []MessageEntity

// ProximityAlert sent whenever a user in the chat triggers
// a proximity alert set by another user.
//
// Bot API: ProximityAlertTriggered.
type ProximityAlert struct {
	Traveler *User `json:"traveler,omitempty"`
	Watcher  *User `json:"watcher,omitempty"`
	Distance int   `json:"distance"`
}

// AutoDeleteTimer represents a service message about a change in auto-delete timer settings.
//
// Bot API: MessageAutoDeleteTimerChanged.
type AutoDeleteTimer struct {
	Unixtime int `json:"message_auto_delete_time"`
}

// Inaccessible shows whether the message is InaccessibleMessage object,
// which Telegram marks with a zero date.
func (m *Message) Inaccessible() bool {
	return m.Unixtime == 0
}

// MessageSig satisfies Editable interface (see Editable.)
func (m *Message) MessageSig() (string, int64) {
	return strconv.Itoa(m.ID), m.Chat.ID
}

// Time returns the moment of message creation in local time.
func (m *Message) Time() time.Time {
	return time.Unix(m.Unixtime, 0)
}

// LastEdited returns time.Time of last edit.
func (m *Message) LastEdited() time.Time {
	return time.Unix(m.LastEdit, 0)
}

// IsForwarded says whether message is forwarded copy of another
// message or not.
func (m *Message) IsForwarded() bool {
	return m.Origin != nil || m.OriginalSender != nil || m.OriginalChat != nil
}

// IsReply says whether message is a reply to another message.
func (m *Message) IsReply() bool {
	return m.ReplyTo != nil
}

// Private returns true, if it's a personal message.
func (m *Message) Private() bool {
	return m.Chat.Type == ChatPrivate
}

// FromGroup returns true, if message came from a group OR a supergroup.
func (m *Message) FromGroup() bool {
	return m.Chat.Type == ChatGroup || m.Chat.Type == ChatSuperGroup
}

// FromChannel returns true, if message came from a channel.
func (m *Message) FromChannel() bool {
	return m.Chat.Type == ChatChannel
}

// IsService returns true, if message is a service message,
// returns false otherwise.
//
// Service messages are automatically sent messages, which
// typically occur on some global action. For instance, when
// anyone leaves the chat or chat title changes.
func (m *Message) IsService() bool {
	fact := false

	fact = fact || m.UserJoined != nil
	fact = fact || len(m.UsersJoined) > 0
	fact = fact || m.UserLeft != nil
	fact = fact || m.NewGroupTitle != ""
	fact = fact || m.NewGroupPhoto != nil
	fact = fact || m.GroupPhotoDeleted
	fact = fact || m.GroupCreated || m.SuperGroupCreated
	fact = fact || (m.MigrateTo != m.MigrateFrom)

	return fact
}

// EntityText returns the substring of the message identified by the
// given MessageEntity.
//
// It's safer than manually slicing Text because Telegram uses
// UTF-16 indices whereas Go string are []byte.
func (m *Message) EntityText(e MessageEntity) string {
	text := m.Text
	if text == "" {
		text = m.Caption
	}

	a := utf16.Encode([]rune(text))
	off, end := e.Offset, e.Offset+e.Length

	if off < 0 || end > len(a) {
		return ""
	}

	return string(utf16.Decode(a[off:end]))
}

// Media returns the message's media if it contains either photo,
// voice, audio, animation, sticker, document, video or video note.
func (m *Message) Media() Media {
	switch {
	case m.Photo != nil:
		return m.Photo
	case m.Voice != nil:
		return m.Voice
	case m.Audio != nil:
		return m.Audio
	case m.Animation != nil:
		return m.Animation
	case m.Sticker != nil:
		return m.Sticker
	case m.Document != nil:
		return m.Document
	case m.Video != nil:
		return m.Video
	case m.VideoNote != nil:
		return m.VideoNote
	default:
		return nil
	}
}

// MessageReaction object represents a change of a reaction on a message performed by a user.
//
// Bot API: MessageReactionUpdated.
type MessageReaction struct {
	// The chat containing the message the user reacted to.
	Chat *Chat `json:"chat"`

	// Unique identifier of the message inside the chat.
	MessageID int `json:"message_id"`

	// (Optional) The user that changed the reaction,
	// if the user isn't anonymous
	User *User `json:"user"`

	// (Optional) The chat on behalf of which the reaction was changed,
	// if the user is anonymous.
	ActorChat *Chat `json:"actor_chat"`

	// Date of the change in Unix time.
	DateUnixtime int64 `json:"date"`

	// Previous list of reaction types that were set by the user.
	OldReaction []Reaction `json:"old_reaction"`

	// New list of reaction types that have been set by the user.
	NewReaction []Reaction `json:"new_reaction"`
}

func (mu *MessageReaction) Time() time.Time {
	return time.Unix(mu.DateUnixtime, 0)
}

// MessageReactionCount represents reaction changes on a message with
// anonymous reactions.
//
// Bot API: MessageReactionCountUpdated.
type MessageReactionCount struct {
	// The chat containing the message.
	Chat *Chat `json:"chat"`

	// Unique message identifier inside the chat.
	MessageID int `json:"message_id"`

	// Date of the change in Unix time.
	DateUnixtime int64 `json:"date"`

	// List of reactions that are present on the message.
	Reactions []*ReactionCount `json:"reactions"`
}

// Time returns the moment of change in local time.
func (mc *MessageReactionCount) Time() time.Time {
	return time.Unix(mc.DateUnixtime, 0)
}

// TextQuote contains information about the quoted part of a message that is
// replied to by the given message.
//
// Bot API: TextQuote.
type TextQuote struct {
	// Text of the quoted part of a message that is replied to by the given message.
	Text string `json:"text"`

	// (Optional) Special entities that appear in the quote.
	// Currently, only bold, italic, underline, strikethrough, spoiler,
	// and custom_emoji entities are kept in quotes.
	Entities []MessageEntity `json:"entities"`

	// Approximate quote position in the original message in UTF-16 code units
	// as specified by the sender.
	Position int `json:"position"`

	// (Optional) True, if the quote was chosen manually by the message sender.
	// Otherwise, the quote was added automatically by the server.
	Manual bool `json:"is_manual"`

	// (Optional) The quote as rich text, if it was chosen from a rich
	// message; Text may then be empty (Bot API 10.4).
	RichTexts []RichText `json:"rich_texts,omitempty"`
}

// MessageOrigin a message reference that has been sent originally by a known user.
//
// Bot API: MessageOrigin, MessageOriginUser, MessageOriginHiddenUser, MessageOriginChat, MessageOriginChannel.
type MessageOrigin struct {
	// Type of the message origin, always “channel”.
	Type string `json:"type"`

	// Date the message was sent originally in Unix time.
	DateUnixtime int64 `json:"date"`

	// User that sent the message originally.
	Sender *User `json:"sender_user,omitempty"`

	// Name of the user that sent the message originally.
	SenderUsername string `json:"sender_user_name,omitempty"`

	// Chat that sent the message originally.
	SenderChat *Chat `json:"sender_chat,omitempty"`

	// Channel chat to which the message was originally sent.
	Chat *Chat `json:"chat,omitempty"`

	// Unique message identifier inside the chat.
	MessageID int `json:"message_id,omitempty"`

	// (Optional) For messages originally sent by an anonymous chat administrator,
	// original message author signature.
	Signature string `json:"author_signature,omitempty"`
}

// Time returns the moment of message that was sent originally in local time.
func (mo *MessageOrigin) Time() time.Time {
	return time.Unix(mo.DateUnixtime, 0)
}

// ExternalReply contains information about a message that is being replied to,
// which may come from another chat or forum topic.
//
// Bot API: ExternalReplyInfo.
type ExternalReply struct {
	// Origin of the message replied to by the given message.
	Origin *MessageOrigin `json:"origin"`

	// (Optional) Chat the original message belongs to.
	// Available only if the chat is a supergroup or a channel.
	Chat *Chat `json:"chat"`

	// (Optional) Unique message identifier inside the original chat.
	// Available only if the original chat is a supergroup or a channel.
	MessageID int `json:"message_id"`

	// (Optional) PreviewOptions used for link preview generation for the original message,
	// if it is a text message.
	PreviewOptions *PreviewOptions `json:"link_preview_options"`

	// (Optional) Message is an animation, information about the animation.
	Animation *Animation `json:"animation"`

	// (Optional) Message is an audio file, information about the file.
	Audio *Audio `json:"audio"`

	// (Optional) Message is a general file, information about the file.
	Document *Document `json:"document"`

	// (Optional) Message is a photo, available sizes of the photo.
	Photo []Photo `json:"photo"`

	// (Optional) Message is a sticker, information about the sticker.
	Sticker *Sticker `json:"sticker"`

	// (Optional) Message is a forwarded story.
	Story *Story `json:"story"`

	// (Optional) Message is a video, information about the video.
	Video *Video `json:"video"`

	// (Optional) Message is a live photo (Bot API 10.0).
	LivePhoto *LivePhoto `json:"live_photo,omitempty"`

	// (Optional) Message is a checklist (Bot API 9.1).
	Checklist *Checklist `json:"checklist,omitempty"`

	// (Optional) Message is a video note, information about the video message.
	Note *VideoNote `json:"video_note"`

	// (Optional) Message is a voice message, information about the file.
	Voice *Voice `json:"voice"`

	// (Optional) Message is a shared contact, information about the contact.
	Contact *Contact `json:"contact"`

	// (Optional) Message is a dice with random value.
	Dice *Dice `json:"dice"`

	//( Optional) Message is a game, information about the game.
	Game *Game `json:"game"`

	// (Optional) Message is a venue, information about the venue.
	Venue *Venue `json:"venue"`

	// (Optional) Message is a native poll, information about the poll.
	Poll *Poll `json:"poll"`

	// (Optional) Message is a shared location, information about the location.
	Location *Location `json:"location"`

	// (Optional) Message is an invoice for a payment, information about the invoice.
	Invoice *Invoice `json:"invoice"`

	// (Optional) Message is a scheduled giveaway, information about the giveaway.
	Giveaway *Giveaway `json:"giveaway"`

	// (Optional) A giveaway with public winners was completed.
	GiveawayWinners *GiveawayWinners `json:"giveaway_winners"`

	// Message contains paid media; information about the paid media
	PaidMedia PaidMedias `json:"paid_media"`

	// (Optional) True, if the message media is covered by a spoiler animation.
	HasMediaSpoiler bool `json:"has_media_spoiler"`

	// (Optional) Message is an HTML widget (Bot API 10.4).
	Widget *Widget `json:"widget,omitempty"`
}

// ReplyParams describes reply parameters for the message that is being sent.
//
// Bot API: ReplyParameters.
type ReplyParams struct {
	// Identifier of the message that will be replied to in the current chat,
	// or in the chat chat_id if it is specified. It is omitted when replying
	// to an ephemeral message.
	MessageID int `json:"message_id,omitempty"`

	// (Optional) Identifier of an incoming ephemeral message, such as an
	// ephemeral command, to reply to within 15 seconds; the reply must be
	// ephemeral too (Bot API 10.2). The bot's own ephemeral messages can't
	// be replied to.
	EphemeralMessageID int `json:"ephemeral_message_id,omitempty"`

	// (Optional) If the message to be replied to is from a different chat,
	// unique identifier for the chat or username of the channel.
	ChatID int64 `json:"chat_id,omitempty"`

	// Optional. Pass True if the message should be sent even if the specified message
	// to be replied to is not found; can be used only for replies in the
	// same chat and forum topic.
	AllowWithoutReply bool `json:"allow_sending_without_reply,omitempty"`

	// (Optional) Quoted part of the message to be replied to; 0-1024 characters after
	// entities parsing. The quote must be an exact substring of the message to be replied to,
	// including bold, italic, underline, strikethrough, spoiler, and custom_emoji entities.
	// The message will fail to send if the quote isn't found in the original message.
	Quote string `json:"quote,omitempty"`

	// (Optional) Mode for parsing entities in the quote.
	QuoteParseMode ParseMode `json:"quote_parse_mode,omitempty"`

	// (Optional) A JSON-serialized list of special entities that appear in the quote.
	// It can be specified instead of quote_parse_mode.
	QuoteEntities []MessageEntity `json:"quote_entities,omitempty"`

	// (Optional) Position of the quote in the original message in UTF-16 code units.
	QuotePosition int `json:"quote_position,omitempty"`

	// (Optional) Quoted parts of a rich message, one per quoted text in the
	// order they appear; if set, Quote, QuoteParseMode and QuoteEntities are
	// ignored (Bot API 10.4).
	QuoteRichTexts []RichText `json:"quote_rich_texts,omitempty"`

	// (Optional) Persistent identifier of the poll option to reply to
	// (Bot API 9.6).
	PollOptionID string `json:"poll_option_id,omitempty"`

	// (Optional) Identifier of the checklist task to reply to (Bot API 9.2).
	ChecklistTaskID int `json:"checklist_task_id,omitempty"`
}

// PaidMessagesWaived is a service message about the sender using a link that
// lets them message the user without paying Telegram Stars (Bot API 10.4).
// It holds no information.
//
// Bot API: PaidMessagesWaived.
type PaidMessagesWaived struct{}

// DirectMessagePriceChanged is a service message about a change in the price
// of direct messages to a channel (Bot API 9.1).
//
// Bot API: DirectMessagePriceChanged.
type DirectMessagePriceChanged struct {
	// Direct messages are enabled for the channel.
	Enabled bool `json:"are_direct_messages_enabled"`

	// (Optional) Stars users must pay for each direct message.
	StarCount int `json:"direct_message_star_count,omitempty"`
}

// PaidMessagePriceChanged is a service message about a change in the Stars
// users pay per message in a supergroup (Bot API 9.0).
//
// Bot API: PaidMessagePriceChanged.
type PaidMessagePriceChanged struct {
	StarCount int `json:"paid_message_star_count"`
}
