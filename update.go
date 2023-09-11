package telebot

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Update object represents an incoming update.
type Update struct {
	ID int `json:"update_id"`

	Message           *Message          `json:"message,omitempty"`
	EditedMessage     *Message          `json:"edited_message,omitempty"`
	ChannelPost       *Message          `json:"channel_post,omitempty"`
	EditedChannelPost *Message          `json:"edited_channel_post,omitempty"`
	Callback          *Callback         `json:"callback_query,omitempty"`
	Query             *Query            `json:"inline_query,omitempty"`
	InlineResult      *InlineResult     `json:"chosen_inline_result,omitempty"`
	ShippingQuery     *ShippingQuery    `json:"shipping_query,omitempty"`
	PreCheckoutQuery  *PreCheckoutQuery `json:"pre_checkout_query,omitempty"`
	Poll              *Poll             `json:"poll,omitempty"`
	PollAnswer        *PollAnswer       `json:"poll_answer,omitempty"`
	MyChatMember      *ChatMemberUpdate `json:"my_chat_member,omitempty"`
	ChatMember        *ChatMemberUpdate `json:"chat_member,omitempty"`
	ChatJoinRequest   *ChatJoinRequest  `json:"chat_join_request,omitempty"`
}

var errHandlerNotFound = errors.New("handler not found")
var errPartiallyHandled = errors.New("only part of message was handled")
var errMaliciousInput = errors.New("malicious message from user")
var errCommandForAnotherBot = errors.New("command is sent to other bot, ignore")

// ProcessUpdateCtx processes a single incoming update and allows to pass standard context.
// A started bot calls this function automatically.
// returns errHandlerNotFound, errPartiallyHandled, errMaliciousInput, errMaliciousInput
func (b *Bot) ProcessUpdateCtx(ctx context.Context, u Update) error {
	c := b.NewContext(u)

	if u.Message != nil {
		m := u.Message

		if m.PinnedMessage != nil {
			if b.handleCtx(ctx, OnPinned, c) {
				return nil
			}
			return fmt.Errorf("pinned: %w", errHandlerNotFound)
		}

		// Commands
		if m.Text != "" {
			// Filtering malicious messages
			if m.Text[0] == '\a' {
				return errMaliciousInput
			}

			match := cmdRx.FindAllStringSubmatch(m.Text, -1)
			if match != nil {
				// Syntax: "</command>@<bot> <payload>"
				command, botName := match[0][1], match[0][3]

				if botName != "" && !strings.EqualFold(b.Me.Username, botName) {
					return fmt.Errorf("'%s' != '%s': %w", b.Me.Username, botName, errCommandForAnotherBot)
				}

				m.Payload = match[0][5]
				if b.handleCtx(ctx, command, c) {
					return nil
				}
			}

			// 1:1 satisfaction
			if b.handleCtx(ctx, m.Text, c) {
				return nil
			}

			if b.handleCtx(ctx, OnText, c) {
				return nil
			}
			return fmt.Errorf("text: %w", errHandlerNotFound)
		}

		if b.handleMediaCtx(ctx, c) {
			return nil
		}

		if m.Contact != nil {
			if b.handleCtx(ctx, OnContact, c) {
				return nil
			}
			return fmt.Errorf("contact: %w", errHandlerNotFound)
		}
		if m.Location != nil {
			if b.handleCtx(ctx, OnLocation, c) {
				return nil
			}
			return fmt.Errorf("location: %w", errHandlerNotFound)
		}
		if m.Venue != nil {
			if b.handleCtx(ctx, OnVenue, c) {
				return nil
			}
			return fmt.Errorf("venue: %w", errHandlerNotFound)
		}
		if m.Game != nil {
			if b.handleCtx(ctx, OnGame, c) {
				return nil
			}
			return fmt.Errorf("game: %w", errHandlerNotFound)
		}
		if m.Dice != nil {
			if b.handleCtx(ctx, OnDice, c) {
				return nil
			}
			return fmt.Errorf("dice: %w", errHandlerNotFound)
		}
		if m.Invoice != nil {
			if b.handleCtx(ctx, OnInvoice, c) {
				return nil
			}
			return fmt.Errorf("invoice: %w", errHandlerNotFound)
		}
		if m.Payment != nil {
			if b.handleCtx(ctx, OnPayment, c) {
				return nil
			}
			return fmt.Errorf("payment: %w", errHandlerNotFound)
		}

		wasAdded := (m.UserJoined != nil && m.UserJoined.ID == b.Me.ID) ||
			(m.UsersJoined != nil && isUserInList(b.Me, m.UsersJoined))
		if m.GroupCreated || m.SuperGroupCreated || wasAdded {
			if b.handleCtx(ctx, OnAddedToGroup, c) {
				return nil
			}
			return fmt.Errorf("added to group: %w", errHandlerNotFound)
		}

		if m.UserJoined != nil {
			if b.handleCtx(ctx, OnUserJoined, c) {
				return nil
			}
			return fmt.Errorf("user joined: %w", errHandlerNotFound)
		}

		if m.UsersJoined != nil {
			var allHandled, anyHandled = true, false
			for _, user := range m.UsersJoined {
				m.UserJoined = &user
				if handled := b.handleCtx(ctx, OnUserJoined, c); handled {
					anyHandled = true
				} else {
					allHandled = false
				}
			}
			if allHandled {
				return nil
			}
			if anyHandled {
				return fmt.Errorf("users joined: %w", errPartiallyHandled)
			}
			return fmt.Errorf("users joined: %w", errHandlerNotFound)
		}

		if m.UserLeft != nil {
			if b.handleCtx(ctx, OnUserLeft, c) {
				return nil
			}
			return fmt.Errorf("user left: %w", errHandlerNotFound)
		}

		if m.NewGroupTitle != "" {
			if b.handleCtx(ctx, OnNewGroupTitle, c) {
				return nil
			}
			return fmt.Errorf("new group title: %w", errHandlerNotFound)
		}

		if m.NewGroupPhoto != nil {
			if b.handleCtx(ctx, OnNewGroupPhoto, c) {
				return nil
			}
			return fmt.Errorf("new group photo: %w", errHandlerNotFound)
		}

		if m.GroupPhotoDeleted {
			if b.handleCtx(ctx, OnGroupPhotoDeleted, c) {
				return nil
			}
			return fmt.Errorf("group photo deleted: %w", errHandlerNotFound)
		}

		if m.GroupCreated {
			if b.handleCtx(ctx, OnGroupCreated, c) {
				return nil
			}
			return fmt.Errorf("group created: %w", errHandlerNotFound)
		}

		if m.SuperGroupCreated {
			if b.handleCtx(ctx, OnSuperGroupCreated, c) {
				return nil
			}
			return fmt.Errorf("super group created: %w", errHandlerNotFound)
		}

		if m.ChannelCreated {
			if b.handleCtx(ctx, OnChannelCreated, c) {
				return nil
			}
			return fmt.Errorf("channel created: %w", errHandlerNotFound)
		}

		if m.MigrateTo != 0 {
			m.MigrateFrom = m.Chat.ID
			if b.handleCtx(ctx, OnMigration, c) {
				return nil
			}
			return fmt.Errorf("migration: %w", errHandlerNotFound)
		}

		if m.VideoChatStarted != nil {
			if b.handleCtx(ctx, OnVideoChatStarted, c) {
				return nil
			}
			return fmt.Errorf("video chat started: %w", errHandlerNotFound)
		}

		if m.VideoChatEnded != nil {
			if b.handleCtx(ctx, OnVideoChatEnded, c) {
				return nil
			}
			return fmt.Errorf("video chat ended: %w", errHandlerNotFound)
		}

		if m.VideoChatParticipants != nil {
			if b.handleCtx(ctx, OnVideoChatParticipants, c) {
				return nil
			}
			return fmt.Errorf("video chat participants: %w", errHandlerNotFound)
		}

		if m.VideoChatScheduled != nil {
			if b.handleCtx(ctx, OnVideoChatScheduled, c) {
				return nil
			}
			return fmt.Errorf("video chat scheduled: %w", errHandlerNotFound)
		}

		if m.WebAppData != nil {
			if b.handleCtx(ctx, OnWebApp, c) {
				return nil
			}
			return fmt.Errorf("web app: %w", errHandlerNotFound)
		}

		if m.ProximityAlert != nil {
			if b.handleCtx(ctx, OnProximityAlert, c) {
				return nil
			}
			return fmt.Errorf("proximity alert: %w", errHandlerNotFound)
		}

		if m.AutoDeleteTimer != nil {
			if b.handleCtx(ctx, OnAutoDeleteTimer, c) {
				return nil
			}
			return fmt.Errorf("autodelete timer: %w", errHandlerNotFound)
		}
	}

	if u.EditedMessage != nil {
		if b.handleCtx(ctx, OnEdited, c) {
			return nil
		}
		return fmt.Errorf("edit: %w", errHandlerNotFound)
	}

	if u.ChannelPost != nil {
		m := u.ChannelPost

		if m.PinnedMessage != nil {
			if b.handleCtx(ctx, OnPinned, c) {
				return nil
			}
			return fmt.Errorf("channel post: pinned: %w", errHandlerNotFound)
		}

		if b.handleCtx(ctx, OnChannelPost, c) {
			return nil
		}
		return fmt.Errorf("channel post: %w", errHandlerNotFound)
	}

	if u.EditedChannelPost != nil {
		if b.handleCtx(ctx, OnEditedChannelPost, c) {
			return nil
		}
		return fmt.Errorf("channel post: edit: %w", errHandlerNotFound)
	}

	if u.Callback != nil {
		if data := u.Callback.Data; data != "" && data[0] == '\f' {
			match := cbackRx.FindAllStringSubmatch(data, -1)
			if match != nil {
				unique, payload := match[0][1], match[0][3]
				if handler, ok := b.handlers["\f"+unique]; ok {
					u.Callback.Unique = unique
					u.Callback.Data = payload
					b.runHandlerCtx(ctx, handler, c)
					return nil
				}
			}
		}

		if b.handleCtx(ctx, OnCallback, c) {
			return nil
		}
		return fmt.Errorf("callback: %w", errHandlerNotFound)
	}

	if u.Query != nil {
		if b.handleCtx(ctx, OnQuery, c) {
			return nil
		}
		return fmt.Errorf("query: %w", errHandlerNotFound)
	}

	if u.InlineResult != nil {
		if b.handleCtx(ctx, OnInlineResult, c) {
			return nil
		}
		return fmt.Errorf("inline result: %w", errHandlerNotFound)
	}

	if u.ShippingQuery != nil {
		if b.handleCtx(ctx, OnShipping, c) {
			return nil
		}
		return fmt.Errorf("shipping: %w", errHandlerNotFound)
	}

	if u.PreCheckoutQuery != nil {
		if b.handleCtx(ctx, OnCheckout, c) {
			return nil
		}
		return fmt.Errorf("checkout: %w", errHandlerNotFound)
	}

	if u.Poll != nil {
		if b.handleCtx(ctx, OnPoll, c) {
			return nil
		}
		return fmt.Errorf("poll: %w", errHandlerNotFound)
	}

	if u.PollAnswer != nil {
		if b.handleCtx(ctx, OnPollAnswer, c) {
			return nil
		}
		return fmt.Errorf("poll answer: %w", errHandlerNotFound)
	}

	if u.MyChatMember != nil {
		if b.handleCtx(ctx, OnMyChatMember, c) {
			return nil
		}
		return fmt.Errorf("my chat member: %w", errHandlerNotFound)
	}

	if u.ChatMember != nil {
		if b.handleCtx(ctx, OnChatMember, c) {
			return nil
		}
		return fmt.Errorf("chat member: %w", errHandlerNotFound)
	}

	if u.ChatJoinRequest != nil {
		if b.handleCtx(ctx, OnChatJoinRequest, c) {
			return nil
		}
		return fmt.Errorf("chat join request: %w", errHandlerNotFound)
	}
	return fmt.Errorf("some new unknown message type: %w", errHandlerNotFound)
}

func (b *Bot) handleCtx(ctx context.Context, end string, c Context) bool {
	if handler, ok := b.handlers[end]; ok {
		b.runHandlerCtx(ctx, handler, c)
		return true
	}
	return false
}

func (b *Bot) handleMediaCtx(ctx context.Context, c Context) bool {
	var (
		m     = c.Message()
		fired = true
	)

	switch {
	case m.Photo != nil:
		fired = b.handleCtx(ctx, OnPhoto, c)
	case m.Voice != nil:
		fired = b.handleCtx(ctx, OnVoice, c)
	case m.Audio != nil:
		fired = b.handleCtx(ctx, OnAudio, c)
	case m.Animation != nil:
		fired = b.handleCtx(ctx, OnAnimation, c)
	case m.Document != nil:
		fired = b.handleCtx(ctx, OnDocument, c)
	case m.Sticker != nil:
		fired = b.handleCtx(ctx, OnSticker, c)
	case m.Video != nil:
		fired = b.handleCtx(ctx, OnVideo, c)
	case m.VideoNote != nil:
		fired = b.handleCtx(ctx, OnVideoNote, c)
	default:
		return false
	}

	if !fired {
		return b.handleCtx(ctx, OnMedia, c)
	}

	return true
}

func (b *Bot) runHandlerCtx(ctx context.Context, h HandlerFunc, c Context) {
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
