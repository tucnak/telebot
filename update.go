package telebot

import "strings"

// Update object represents an incoming update.
//
// Bot API: Update.
type Update struct {
	ID int `json:"update_id"`

	Message                 *Message                  `json:"message,omitempty"`
	EditedMessage           *Message                  `json:"edited_message,omitempty"`
	ChannelPost             *Message                  `json:"channel_post,omitempty"`
	EditedChannelPost       *Message                  `json:"edited_channel_post,omitempty"`
	MessageReaction         *MessageReaction          `json:"message_reaction"`
	MessageReactionCount    *MessageReactionCount     `json:"message_reaction_count"`
	Callback                *Callback                 `json:"callback_query,omitempty"`
	Query                   *Query                    `json:"inline_query,omitempty"`
	InlineResult            *InlineResult             `json:"chosen_inline_result,omitempty"`
	ShippingQuery           *ShippingQuery            `json:"shipping_query,omitempty"`
	PreCheckoutQuery        *PreCheckoutQuery         `json:"pre_checkout_query,omitempty"`
	Poll                    *Poll                     `json:"poll,omitempty"`
	PollAnswer              *PollAnswer               `json:"poll_answer,omitempty"`
	MyChatMember            *ChatMemberUpdate         `json:"my_chat_member,omitempty"`
	ChatMember              *ChatMemberUpdate         `json:"chat_member,omitempty"`
	ChatJoinRequest         *ChatJoinRequest          `json:"chat_join_request,omitempty"`
	Boost                   *BoostUpdated             `json:"chat_boost"`
	BoostRemoved            *BoostRemoved             `json:"removed_chat_boost"`
	BusinessConnection      *BusinessConnection       `json:"business_connection"`
	BusinessMessage         *Message                  `json:"business_message"`
	EditedBusinessMessage   *Message                  `json:"edited_business_message"`
	DeletedBusinessMessages *BusinessMessagesDeleted  `json:"deleted_business_messages"`
	GuestMessage            *Message                  `json:"guest_message"`
	PurchasedPaidMedia      *PaidMediaPurchased       `json:"purchased_paid_media"`
	Subscription            *BotSubscriptionUpdated   `json:"subscription,omitempty"`
	GenerationStopped       *MessageGenerationStopped `json:"stopped_message_generation,omitempty"`
	WidgetQuery             *WidgetQuery              `json:"widget_query,omitempty"`
	SharedWidgetQuery       *SharedWidgetQuery        `json:"shared_widget_query,omitempty"`
	ManagedBot              *ManagedBotUpdated        `json:"managed_bot,omitempty"`
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

	if u.Message != nil {
		m := u.Message

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

		if m.Checklist != nil {
			b.handle(OnChecklist, c)
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
		if m.Gift != nil {
			b.handle(OnGift, c)
			return
		}
		if m.UniqueGift != nil {
			b.handle(OnUniqueGift, c)
			return
		}
		if m.GiftUpgradeSent != nil {
			b.handle(OnGiftUpgradeSent, c)
			return
		}
		if m.TonWalletTransfer != nil {
			b.handle(OnTonWalletTransfer, c)
			return
		}
		if m.PaidMessagesWaived != nil {
			b.handle(OnPaidMessagesWaived, c)
			return
		}
		if m.ManagedBotCreated != nil {
			b.handle(OnManagedBotCreated, c)
			return
		}
		if m.SuggestedPostApproved != nil {
			b.handle(OnSuggestedPostApproved, c)
			return
		}
		if m.SuggestedPostApprovalFailed != nil {
			b.handle(OnSuggestedPostApprovalFailed, c)
			return
		}
		if m.SuggestedPostDeclined != nil {
			b.handle(OnSuggestedPostDeclined, c)
			return
		}
		if m.SuggestedPostPaid != nil {
			b.handle(OnSuggestedPostPaid, c)
			return
		}
		if m.SuggestedPostRefunded != nil {
			b.handle(OnSuggestedPostRefunded, c)
			return
		}
		if m.PassportData != nil {
			b.handle(OnPassportData, c)
			return
		}
		if m.PaidMessagePriceChanged != nil {
			b.handle(OnPaidMessagePriceChanged, c)
			return
		}
		if m.ChecklistTasksDone != nil {
			b.handle(OnChecklistTasksDone, c)
			return
		}
		if m.ChecklistTasksAdded != nil {
			b.handle(OnChecklistTasksAdded, c)
			return
		}
		if m.DirectMessagePriceChanged != nil {
			b.handle(OnDirectMessagePriceChanged, c)
			return
		}
		if m.PollOptionAdded != nil {
			b.handle(OnPollOptionAdded, c)
			return
		}
		if m.PollOptionDeleted != nil {
			b.handle(OnPollOptionDeleted, c)
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
		if m.ChatOwnerLeft != nil {
			b.handle(OnChatOwnerLeft, c)
			return
		}
		if m.ChatOwnerChanged != nil {
			b.handle(OnChatOwnerChanged, c)
			return
		}
		if m.CommunityChatAdded != nil {
			b.handle(OnCommunityChatAdded, c)
			return
		}
		if m.CommunityChatRemoved != nil {
			b.handle(OnCommunityChatRemoved, c)
			return
		}
		if m.CommunityChatJoined != nil {
			b.handle(OnCommunityChatJoined, c)
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

	if u.EditedMessage != nil {
		b.handle(OnEdited, c)
		return
	}

	if u.ChannelPost != nil {
		m := u.ChannelPost

		if m.PinnedMessage != nil {
			b.handle(OnPinned, c)
			return
		}

		b.handle(OnChannelPost, c)
		return
	}

	if u.EditedChannelPost != nil {
		b.handle(OnEditedChannelPost, c)
		return
	}

	if u.Callback != nil {
		if data := u.Callback.Data; data != "" && data[0] == '\f' {
			match := cbackRx.FindAllStringSubmatch(data, -1)
			if match != nil {
				unique, payload := match[0][1], match[0][3]
				if handler, ok := b.handlers["\f"+unique]; ok {
					u.Callback.Unique = unique
					u.Callback.Data = payload
					b.runHandler(handler, c)
					return
				}
			}
		}

		b.handle(OnCallback, c)
		return
	}

	if u.Query != nil {
		b.handle(OnQuery, c)
		return
	}

	if u.InlineResult != nil {
		b.handle(OnInlineResult, c)
		return
	}

	if u.ShippingQuery != nil {
		b.handle(OnShipping, c)
		return
	}

	if u.PreCheckoutQuery != nil {
		b.handle(OnCheckout, c)
		return
	}

	if u.Poll != nil {
		b.handle(OnPoll, c)
		return
	}
	if u.PollAnswer != nil {
		b.handle(OnPollAnswer, c)
		return
	}

	if u.MyChatMember != nil {
		b.handle(OnMyChatMember, c)
		return
	}
	if u.ChatMember != nil {
		b.handle(OnChatMember, c)
		return
	}
	if u.ChatJoinRequest != nil {
		b.handle(OnChatJoinRequest, c)
		return
	}

	if u.Boost != nil {
		b.handle(OnBoost, c)
		return
	}
	if u.BoostRemoved != nil {
		b.handle(OnBoostRemoved, c)
		return
	}

	if u.BusinessConnection != nil {
		b.handle(OnBusinessConnection, c)
		return
	}
	if u.BusinessMessage != nil {
		b.handle(OnBusinessMessage, c)
		return
	}
	if u.EditedBusinessMessage != nil {
		b.handle(OnEditedBusinessMessage, c)
		return
	}
	if u.DeletedBusinessMessages != nil {
		b.handle(OnDeletedBusinessMessages, c)
		return
	}
	if u.GuestMessage != nil {
		b.handle(OnGuestMessage, c)
		return
	}
	if u.PurchasedPaidMedia != nil {
		b.handle(OnPurchasedPaidMedia, c)
		return
	}
	if u.Subscription != nil {
		b.handle(OnSubscription, c)
		return
	}
	if u.MessageReaction != nil {
		b.handle(OnReaction, c)
		return
	}
	if u.MessageReactionCount != nil {
		b.handle(OnReactionCount, c)
		return
	}
	if u.GenerationStopped != nil {
		b.handle(OnGenerationStopped, c)
		return
	}
	if u.WidgetQuery != nil {
		b.handle(OnWidgetQuery, c)
		return
	}
	if u.SharedWidgetQuery != nil {
		b.handle(OnSharedWidgetQuery, c)
		return
	}
	if u.ManagedBot != nil {
		b.handle(OnManagedBot, c)
		return
	}
}

func (b *Bot) handle(end string, c Context) bool {
	if handler, ok := b.handlers[end]; ok {
		b.runHandler(handler, c)
		return true
	}
	return false
}

func (b *Bot) handleMedia(c Context) bool {
	var (
		m     = c.Message()
		fired = true
	)

	// A live photo also carries its still in Photo, so without an
	// OnLivePhoto handler it is handled as a photo.
	if m.LivePhoto != nil && b.handle(OnLivePhoto, c) {
		return true
	}

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
