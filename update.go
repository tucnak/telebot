package telebot

import "strings"

// Update object represents an incoming update.
type Update struct {
	ID int `json:"update_id"`

	Message                 *Message                 `json:"message,omitempty"`
	EditedMessage           *Message                 `json:"edited_message,omitempty"`
	ChannelPost             *Message                 `json:"channel_post,omitempty"`
	EditedChannelPost       *Message                 `json:"edited_channel_post,omitempty"`
	MessageReaction         *MessageReaction         `json:"message_reaction"`
	MessageReactionCount    *MessageReactionCount    `json:"message_reaction_count"`
	Callback                *Callback                `json:"callback_query,omitempty"`
	Query                   *Query                   `json:"inline_query,omitempty"`
	InlineResult            *InlineResult            `json:"chosen_inline_result,omitempty"`
	ShippingQuery           *ShippingQuery           `json:"shipping_query,omitempty"`
	PreCheckoutQuery        *PreCheckoutQuery        `json:"pre_checkout_query,omitempty"`
	Poll                    *Poll                    `json:"poll,omitempty"`
	PollAnswer              *PollAnswer              `json:"poll_answer,omitempty"`
	MyChatMember            *ChatMemberUpdate        `json:"my_chat_member,omitempty"`
	ChatMember              *ChatMemberUpdate        `json:"chat_member,omitempty"`
	ChatJoinRequest         *ChatJoinRequest         `json:"chat_join_request,omitempty"`
	Boost                   *BoostUpdated            `json:"chat_boost"`
	BoostRemoved            *BoostRemoved            `json:"removed_chat_boost"`
	BusinessConnection      *BusinessConnection      `json:"business_connection"`
	BusinessMessage         *Message                 `json:"business_message"`
	EditedBusinessMessage   *Message                 `json:"edited_business_message"`
	DeletedBusinessMessages *BusinessMessagesDeleted `json:"deleted_business_messages"`
}

func (u Update) String() string {
	switch {
	case u.Message != nil:
		return OnMessage
	case u.EditedMessage != nil:
		return OnEdited
	case u.ChannelPost != nil:
		if u.ChannelPost.PinnedMessage != nil {
			return OnPinned
		}
		return OnChannelPost
	case u.EditedChannelPost != nil:
		return OnEditedChannelPost
	case u.Callback != nil:
		return OnCallback
	case u.Query != nil:
		return OnQuery
	case u.InlineResult != nil:
		return OnInlineResult
	case u.ShippingQuery != nil:
		return OnShipping
	case u.PreCheckoutQuery != nil:
		return OnCheckout
	case u.Poll != nil:
		return OnPoll
	case u.PollAnswer != nil:
		return OnPollAnswer
	case u.MyChatMember != nil:
		return OnMyChatMember
	case u.ChatMember != nil:
		return OnChatMember
	case u.ChatJoinRequest != nil:
		return OnChatJoinRequest
	case u.Boost != nil:
		return OnBoost
	case u.BoostRemoved != nil:
		return OnBoostRemoved
	default:
		return ""
	}
}

// ProcessUpdate processes a single incoming update.
// A started bot calls this function automatically.
func (b *Bot) ProcessUpdate(u Update) {
	b.ProcessContext(b.NewContext(u))
}

// ProcessContext processes the given context.
// A started bot calls this function automatically.
func (b *Bot) ProcessContext(c Context) {
	u := c.Update()

	switch true {
	case b.flowManager.IsFollowed(c.Recipient()): // handle flow
		if u.Callback != nil {
			b.handleCallback(u)
		}

		if h := b.flowManager.MakeProcessing(c, u); h != nil {
			if err := applyMiddleware(h, b.group.middleware...)(c); err != nil {
				b.OnError(err, c)
				return
			}
		}

		if b.flowManager.Close(c.Recipient()) {
			return
		}

		f := b.flowManager.MakeTransition(c, u)
		if f != nil {
			b.runHandler(f.Enter(b.group.middleware...), c)

			return
		}
	case u.ChannelPost != nil:
		m := u.ChannelPost

		if m.PinnedMessage != nil {
			b.handle(OnPinned, c)
			return
		}

		b.handle(OnChannelPost, c)
		return
	case u.Callback != nil: // processing callbacks
		unique, handled := b.handleCallback(u)
		if !handled {
			b.handle(OnCallback, c)
			return
		}

		if handler, ok := b.handlers["\f"+unique]; ok {
			b.runHandler(handler, c)
			return
		}

		f := b.flowManager.Get(unique)
		if f != nil {
			b.flowManager.Start(c, f)

			b.runHandler(f.Enter(b.group.middleware...), c)
			return
		}
	case u.String() != "": // processing all other handlers
		b.handle(OnQuery, c)
		return
	}
}

func (b *Bot) handle(end string, c Context) bool {
	handler, ok := b.handlers[end]
	if !ok {
		return false
	}

	b.runHandler(handler, c)
	return true
}

func (b *Bot) handleCallback(u Update) (string, bool) {
	data := u.Callback.Data
	if data == "" || data[0] != '\f' {
		return "", false
	}

	match := cbackRx.FindAllStringSubmatch(data, -1)
	if match == nil {
		return "", false
	}

	unique, payload := match[0][1], match[0][3]

	u.Callback.Unique = unique
	u.Callback.Data = payload

	return unique, true
}

func (b *Bot) handleMedia(c Context) bool {
	var (
		m     = c.Message()
		fired = true
	)

	switch {
	case m.Photo != nil:
		fired = b.handle(OnPhoto, c)
	case m.Voice != nil:
		fired = b.handle(OnVoice, c)
	case m.Audio != nil:
		fired = b.handle(OnAudio, c)
	case m.Animation != nil:
		fired = b.handle(OnAnimation, c)
	case m.Document != nil:
		fired = b.handle(OnDocument, c)
	case m.Sticker != nil:
		fired = b.handle(OnSticker, c)
	case m.Video != nil:
		fired = b.handle(OnVideo, c)
	case m.VideoNote != nil:
		fired = b.handle(OnVideoNote, c)
	default:
		return false
	}

	if !fired {
		return b.handle(OnMedia, c)
	}

	return true
}

func (b *Bot) handleMessage(c Context, m *Message) {
	if m.PinnedMessage != nil {
		b.handle(OnPinned, c)
		return
	}

	if m.Origin != nil {
		b.handle(OnForward, c)
	}

	// Commands
	if m.Text != "" {
		// Filtering malicious messages
		if m.Text[0] == '\a' {
			return
		}

		match := cmdRx.FindAllStringSubmatch(m.Text, -1)
		if match != nil {
			// Syntax: "</command>@<bot> <payload>"
			command, botName := match[0][1], match[0][3]

			if botName != "" && !strings.EqualFold(b.Me.Username, botName) {
				return
			}

			m.Payload = match[0][5]

			if b.flowManager.IsRegistred(command) {
				f := b.flowManager.Get(command)

				b.flowManager.Start(c, f)
				b.runHandler(f.Enter(b.group.middleware...), c)
				defer b.flowManager.Close(c.Recipient())

				return
			}

			if b.handle(command, c) {
				return
			}
		}

		// 1:1 satisfaction
		if b.handle(m.Text, c) {
			return
		}

		if m.ReplyTo != nil {
			b.handle(OnReply, c)
		}

		b.handle(OnText, c)
		return
	}

	if b.handleMedia(c) {
		return
	}

	if m.Contact != nil {
		b.handle(OnContact, c)
		return
	}
	if m.Location != nil {
		b.handle(OnLocation, c)
		return
	}
	if m.Venue != nil {
		b.handle(OnVenue, c)
		return
	}
	if m.Game != nil {
		b.handle(OnGame, c)
		return
	}
	if m.Dice != nil {
		b.handle(OnDice, c)
		return
	}
	if m.Invoice != nil {
		b.handle(OnInvoice, c)
		return
	}
	if m.Payment != nil {
		b.handle(OnPayment, c)
		return
	}
	if m.RefundedPayment != nil {
		b.handle(OnRefund, c)
		return
	}
	if m.TopicCreated != nil {
		b.handle(OnTopicCreated, c)
		return
	}
	if m.TopicReopened != nil {
		b.handle(OnTopicReopened, c)
		return
	}
	if m.TopicClosed != nil {
		b.handle(OnTopicClosed, c)
		return
	}
	if m.TopicEdited != nil {
		b.handle(OnTopicEdited, c)
		return
	}
	if m.GeneralTopicHidden != nil {
		b.handle(OnGeneralTopicHidden, c)
		return
	}
	if m.GeneralTopicUnhidden != nil {
		b.handle(OnGeneralTopicUnhidden, c)
		return
	}
	if m.WriteAccessAllowed != nil {
		b.handle(OnWriteAccessAllowed, c)
		return
	}

	wasAdded := (m.UserJoined != nil && m.UserJoined.ID == b.Me.ID) ||
		(m.UsersJoined != nil && isUserInList(b.Me, m.UsersJoined))
	if m.GroupCreated || m.SuperGroupCreated || wasAdded {
		b.handle(OnAddedToGroup, c)
		return
	}

	if m.UserJoined != nil {
		b.handle(OnUserJoined, c)
		return
	}
	if m.UsersJoined != nil {
		for _, user := range m.UsersJoined {
			m.UserJoined = &user
			b.handle(OnUserJoined, c)
		}
		return
	}
	if m.UserLeft != nil {
		b.handle(OnUserLeft, c)
		return
	}

	if m.UserShared != nil {
		b.handle(OnUserShared, c)
		return
	}
	if m.ChatShared != nil {
		b.handle(OnChatShared, c)
		return
	}

	if m.NewGroupTitle != "" {
		b.handle(OnNewGroupTitle, c)
		return
	}
	if m.NewGroupPhoto != nil {
		b.handle(OnNewGroupPhoto, c)
		return
	}
	if m.GroupPhotoDeleted {
		b.handle(OnGroupPhotoDeleted, c)
		return
	}

	if m.GroupCreated {
		b.handle(OnGroupCreated, c)
		return
	}
	if m.SuperGroupCreated {
		b.handle(OnSuperGroupCreated, c)
		return
	}
	if m.ChannelCreated {
		b.handle(OnChannelCreated, c)
		return
	}

	if m.MigrateTo != 0 {
		m.MigrateFrom = m.Chat.ID
		b.handle(OnMigration, c)
		return
	}

	if m.VideoChatStarted != nil {
		b.handle(OnVideoChatStarted, c)
		return
	}
	if m.VideoChatEnded != nil {
		b.handle(OnVideoChatEnded, c)
		return
	}
	if m.VideoChatParticipants != nil {
		b.handle(OnVideoChatParticipants, c)
		return
	}
	if m.VideoChatScheduled != nil {
		b.handle(OnVideoChatScheduled, c)
		return
	}

	if m.WebAppData != nil {
		b.handle(OnWebApp, c)
		return
	}

	if m.ProximityAlert != nil {
		b.handle(OnProximityAlert, c)
		return
	}
	if m.AutoDeleteTimer != nil {
		b.handle(OnAutoDeleteTimer, c)
		return
	}
}

func (b *Bot) runHandler(h HandlerFunc, c Context) {
	f := func() {
		if err := h(c); err != nil {
			b.OnError(err, c)
		}
	}
	if b.synchronous {
		f()
	} else {
		go f()
	}
}

func isUserInList(user *User, list []User) bool {
	for _, user2 := range list {
		if user.ID == user2.ID {
			return true
		}
	}
	return false
}
