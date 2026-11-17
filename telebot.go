// Package telebot is a framework for Telegram bots.
//
// Example:
//
//	package main
//
//	import (
//		"time"
//		tele "gopkg.in/telebot.v4"
//	)
//
//	func main() {
//		b, err := tele.NewBot(tele.Settings{
//			Token:  "...",
//			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
//		})
//		if err != nil {
//			return
//		}
//
//		b.Handle("/start", func(c tele.Context) error {
//			return c.Send("Hello world!")
//		})
//
//		b.Start()
//	}
package telebot

import "errors"

var (
	ErrBadRecipient    = errors.New("telebot: recipient is nil")
	ErrUnsupportedWhat = errors.New("telebot: unsupported what argument")
	ErrCouldNotUpdate  = errors.New("telebot: could not fetch new updates")
	ErrTrueResult      = errors.New("telebot: result is True")
	ErrBadContext      = errors.New("telebot: context does not contain message")
)

const DefaultApiURL = "https://api.telegram.org"

// These are one of the possible events Handle() can deal with.
//
// For convenience, all Telebot-provided endpoints start with
// an "alert" character \a.
const (
	// Basic message handlers.
	OnText                        = "\atext"
	OnForward                     = "\aforward"
	OnReply                       = "\areply"
	OnEdited                      = "\aedited"
	OnPhoto                       = "\aphoto"
	OnAudio                       = "\aaudio"
	OnAnimation                   = "\aanimation"
	OnDocument                    = "\adocument"
	OnSticker                     = "\asticker"
	OnVideo                       = "\avideo"
	OnVoice                       = "\avoice"
	OnLivePhoto                   = "\alive_photo"
	OnChecklist                   = "\achecklist"
	OnVideoNote                   = "\avideo_note"
	OnContact                     = "\acontact"
	OnLocation                    = "\alocation"
	OnVenue                       = "\avenue"
	OnDice                        = "\adice"
	OnInvoice                     = "\ainvoice"
	OnPayment                     = "\apayment"
	OnRefund                      = "\arefund"
	OnGift                        = "\agift"
	OnUniqueGift                  = "\aunique_gift"
	OnGiftUpgradeSent             = "\agift_upgrade_sent"
	OnTonWalletTransfer           = "\aton_wallet_transfer"
	OnPaidMessagesWaived          = "\apaid_messages_waived"
	OnManagedBotCreated           = "\amanaged_bot_created"
	OnPollOptionAdded             = "\apoll_option_added"
	OnSuggestedPostApproved       = "\asuggested_post_approved"
	OnSuggestedPostApprovalFailed = "\asuggested_post_approval_failed"
	OnSuggestedPostDeclined       = "\asuggested_post_declined"
	OnSuggestedPostPaid           = "\asuggested_post_paid"
	OnSuggestedPostRefunded       = "\asuggested_post_refunded"
	OnPassportData                = "\apassport_data"
	OnPaidMessagePriceChanged     = "\apaid_message_price_changed"
	OnChecklistTasksDone          = "\achecklist_tasks_done"
	OnChecklistTasksAdded         = "\achecklist_tasks_added"
	OnDirectMessagePriceChanged   = "\adirect_message_price_changed"
	OnPollOptionDeleted           = "\apoll_option_deleted"
	OnGame                        = "\agame"
	OnPoll                        = "\apoll"
	OnPollAnswer                  = "\apoll_answer"
	OnPinned                      = "\apinned"
	OnChannelPost                 = "\achannel_post"
	OnEditedChannelPost           = "\aedited_channel_post"
	OnTopicCreated                = "\atopic_created"
	OnTopicReopened               = "\atopic_reopened"
	OnTopicClosed                 = "\atopic_closed"
	OnTopicEdited                 = "\atopic_edited"
	OnGeneralTopicHidden          = "\ageneral_topic_hidden"
	OnGeneralTopicUnhidden        = "\ageneral_topic_unhidden"
	OnWriteAccessAllowed          = "\awrite_access_allowed"
	OnChatOwnerLeft               = "\achat_owner_left"
	OnChatOwnerChanged            = "\achat_owner_changed"
	OnCommunityChatAdded          = "\acommunity_chat_added"
	OnCommunityChatRemoved        = "\acommunity_chat_removed"
	OnCommunityChatJoined         = "\acommunity_chat_joined"

	OnAddedToGroup      = "\aadded_to_group"
	OnUserJoined        = "\auser_joined"
	OnUserLeft          = "\auser_left"
	OnUserShared        = "\auser_shared"
	OnChatShared        = "\achat_shared"
	OnNewGroupTitle     = "\anew_chat_title"
	OnNewGroupPhoto     = "\anew_chat_photo"
	OnGroupPhotoDeleted = "\achat_photo_deleted"
	OnGroupCreated      = "\agroup_created"
	OnSuperGroupCreated = "\asupergroup_created"
	OnChannelCreated    = "\achannel_created"

	// OnMigration happens when group switches to
	// a supergroup. You might want to update
	// your internal references to this chat
	// upon switching as its ID will change.
	OnMigration = "\amigration"

	OnMedia           = "\amedia"
	OnCallback        = "\acallback"
	OnQuery           = "\aquery"
	OnInlineResult    = "\ainline_result"
	OnShipping        = "\ashipping_query"
	OnCheckout        = "\apre_checkout_query"
	OnMyChatMember    = "\amy_chat_member"
	OnChatMember      = "\achat_member"
	OnChatJoinRequest = "\achat_join_request"
	OnProximityAlert  = "\aproximity_alert_triggered"
	OnAutoDeleteTimer = "\amessage_auto_delete_timer_changed"
	OnWebApp          = "\aweb_app"

	OnVideoChatStarted      = "\avideo_chat_started"
	OnVideoChatEnded        = "\avideo_chat_ended"
	OnVideoChatParticipants = "\avideo_chat_participants_invited"
	OnVideoChatScheduled    = "\avideo_chat_scheduled"

	OnBoost        = "\aboost_updated"
	OnBoostRemoved = "\aboost_removed"

	OnBusinessConnection      = "\abusiness_connection"
	OnBusinessMessage         = "\abusiness_message"
	OnEditedBusinessMessage   = "\aedited_business_message"
	OnDeletedBusinessMessages = "\adeleted_business_messages"
	OnGuestMessage            = "\aguest_message"
	OnPurchasedPaidMedia      = "\apurchased_paid_media"
	OnSubscription            = "\asubscription"
	OnReaction                = "\amessage_reaction"
	OnReactionCount           = "\amessage_reaction_count"
	OnGenerationStopped       = "\astopped_message_generation"
	OnWidgetQuery             = "\awidget_query"
	OnSharedWidgetQuery       = "\ashared_widget_query"
	OnManagedBot              = "\amanaged_bot"
)

// ChatAction is a client-side status indicating bot activity.
type ChatAction string

const (
	Typing            ChatAction = "typing"
	UploadingPhoto    ChatAction = "upload_photo"
	UploadingVideo    ChatAction = "upload_video"
	UploadingAudio    ChatAction = "upload_audio"
	UploadingDocument ChatAction = "upload_document"
	UploadingVNote    ChatAction = "upload_video_note"
	RecordingVideo    ChatAction = "record_video"
	RecordingAudio    ChatAction = "record_audio"
	RecordingVNote    ChatAction = "record_video_note"
	FindingLocation   ChatAction = "find_location"
	ChoosingSticker   ChatAction = "choose_sticker"
)

// ParseMode determines the way client applications treat the text of the message
type ParseMode = string

const (
	ModeDefault    ParseMode = ""
	ModeMarkdown   ParseMode = "Markdown"
	ModeMarkdownV2 ParseMode = "MarkdownV2"
	ModeHTML       ParseMode = "HTML"
)

// M is a shortcut for map[string]interface{}.
// Useful for passing arguments to the layout functions.
type M = map[string]interface{}

// Flag returns a pointer to the given bool.
// Useful for passing the three-state flags to a Bot API.
// For example, see ReplyRecipient type.
func Flag(b bool) *bool {
	return &b
}
