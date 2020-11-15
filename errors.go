package telebot

import (
	"fmt"
	"strings"
)

type APIError struct {
	Code        int
	Description string
	Message     string
	Parameters  map[string]interface{}
}

type FloodError struct {
	*APIError
	RetryAfter int
}

// ʔ returns description of error.
// A tiny shortcut to make code clearer.
func (err *APIError) ʔ() string {
	return err.Description
}

// Error implements error interface.
func (err *APIError) Error() string {
	msg := err.Message
	if msg == "" {
		split := strings.Split(err.Description, ": ")
		if len(split) == 2 {
			msg = split[1]
		} else {
			msg = err.Description
		}
	}
	return fmt.Sprintf("telegram: %s (%d)", msg, err.Code)
}

// NewAPIError returns new APIError instance with given description.
// First element of msgs is Description. The second is optional Message.
func NewAPIError(code int, msgs ...string) *APIError {
	err := &APIError{Code: code}
	if len(msgs) >= 1 {
		err.Description = msgs[0]
	}
	if len(msgs) >= 2 {
		err.Message = msgs[1]
	}
	return err
}

var (
	// General errors
	ErrUnauthorized      = NewAPIError(401, "Unauthorized")
	ErrNotStartedByUser  = NewAPIError(403, "Forbidden: bot can't initiate conversation with a user")
	ErrBlockedByUser     = NewAPIError(401, "Forbidden: bot was blocked by the user")
	ErrUserIsDeactivated = NewAPIError(401, "Forbidden: user is deactivated")
	ErrNotFound          = NewAPIError(404, "Not Found")
	ErrInternal          = NewAPIError(500, "Internal Server Error")

	// Bad request errors
	ErrTooLarge             = NewAPIError(400, "Request Entity Too Large")
	ErrMessageTooLong       = NewAPIError(400, "Bad Request: message is too long")
	ErrToForwardNotFound    = NewAPIError(400, "Bad Request: message to forward not found")
	ErrToReplyNotFound      = NewAPIError(400, "Bad Request: reply message not found")
	ErrToDeleteNotFound     = NewAPIError(400, "Bad Request: message to delete not found")
	ErrEmptyMessage         = NewAPIError(400, "Bad Request: message must be non-empty")
	ErrEmptyText            = NewAPIError(400, "Bad Request: text is empty")
	ErrEmptyChatID          = NewAPIError(400, "Bad Request: chat_id is empty")
	ErrChatNotFound         = NewAPIError(400, "Bad Request: chat not found")
	ErrMessageNotModified   = NewAPIError(400, "Bad Request: message is not modified")
	ErrSameMessageContent   = NewAPIError(400, "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message")
	ErrCantEditMessage      = NewAPIError(400, "Bad Request: message can't be edited")
	ErrButtonDataInvalid    = NewAPIError(400, "Bad Request: BUTTON_DATA_INVALID")
	ErrWrongTypeOfContent   = NewAPIError(400, "Bad Request: wrong type of the web page content")
	ErrBadURLContent        = NewAPIError(400, "Bad Request: failed to get HTTP URL content")
	ErrWrongFileID          = NewAPIError(400, "Bad Request: wrong file identifier/HTTP URL specified")
	ErrWrongFileIDSymbol    = NewAPIError(400, "Bad Request: wrong remote file id specified: can't unserialize it. Wrong last symbol")
	ErrWrongFileIDLength    = NewAPIError(400, "Bad Request: wrong remote file id specified: Wrong string length")
	ErrWrongFileIDCharacter = NewAPIError(400, "Bad Request: wrong remote file id specified: Wrong character in the string")
	ErrWrongFileIDPadding   = NewAPIError(400, "Bad Request: wrong remote file id specified: Wrong padding in the string")
	ErrFailedImageProcess   = NewAPIError(400, "Bad Request: IMAGE_PROCESS_FAILED", "Image process failed")
	ErrInvalidStickerSet    = NewAPIError(400, "Bad Request: STICKERSET_INVALID", "Stickerset is invalid")
	ErrBadPollOptions       = NewAPIError(400, "Bad Request: expected an Array of String as options")
	ErrGroupMigrated        = NewAPIError(400, "Bad Request: group chat was upgraded to a supergroup chat")

	// No rights errors
	ErrNoRightsToRestrict     = NewAPIError(400, "Bad Request: not enough rights to restrict/unrestrict chat member")
	ErrNoRightsToSend         = NewAPIError(400, "Bad Request: have no rights to send a message")
	ErrNoRightsToSendPhoto    = NewAPIError(400, "Bad Request: not enough rights to send photos to the chat")
	ErrNoRightsToSendStickers = NewAPIError(400, "Bad Request: not enough rights to send stickers to the chat")
	ErrNoRightsToSendGifs     = NewAPIError(400, "Bad Request: CHAT_SEND_GIFS_FORBIDDEN", "sending GIFS is not allowed in this chat")
	ErrNoRightsToDelete       = NewAPIError(400, "Bad Request: message can't be deleted")
	ErrKickingChatOwner       = NewAPIError(400, "Bad Request: can't remove chat owner")

	// Super/groups errors
	ErrBotKickedFromGroup      = NewAPIError(403, "Forbidden: bot was kicked from the group chat")
	ErrBotKickedFromSuperGroup = NewAPIError(403, "Forbidden: bot was kicked from the supergroup chat")
)

// ErrByDescription returns APIError instance by given description.
func ErrByDescription(s string) error {
	switch s {
	case ErrUnauthorized.ʔ():
		return ErrUnauthorized
	case ErrNotStartedByUser.ʔ():
		return ErrNotStartedByUser
	case ErrNotFound.ʔ():
		return ErrNotFound
	case ErrUserIsDeactivated.ʔ():
		return ErrUserIsDeactivated
	case ErrToForwardNotFound.ʔ():
		return ErrToForwardNotFound
	case ErrToReplyNotFound.ʔ():
		return ErrToReplyNotFound
	case ErrMessageTooLong.ʔ():
		return ErrMessageTooLong
	case ErrBlockedByUser.ʔ():
		return ErrBlockedByUser
	case ErrToDeleteNotFound.ʔ():
		return ErrToDeleteNotFound
	case ErrEmptyMessage.ʔ():
		return ErrEmptyMessage
	case ErrEmptyText.ʔ():
		return ErrEmptyText
	case ErrEmptyChatID.ʔ():
		return ErrEmptyChatID
	case ErrChatNotFound.ʔ():
		return ErrChatNotFound
	case ErrMessageNotModified.ʔ():
		return ErrMessageNotModified
	case ErrSameMessageContent.ʔ():
		return ErrSameMessageContent
	case ErrCantEditMessage.ʔ():
		return ErrCantEditMessage
	case ErrButtonDataInvalid.ʔ():
		return ErrButtonDataInvalid
	case ErrBadPollOptions.ʔ():
		return ErrBadPollOptions
	case ErrNoRightsToRestrict.ʔ():
		return ErrNoRightsToRestrict
	case ErrNoRightsToSend.ʔ():
		return ErrNoRightsToSend
	case ErrNoRightsToSendPhoto.ʔ():
		return ErrNoRightsToSendPhoto
	case ErrNoRightsToSendStickers.ʔ():
		return ErrNoRightsToSendStickers
	case ErrNoRightsToSendGifs.ʔ():
		return ErrNoRightsToSendGifs
	case ErrNoRightsToDelete.ʔ():
		return ErrNoRightsToDelete
	case ErrKickingChatOwner.ʔ():
		return ErrKickingChatOwner
	case ErrBotKickedFromGroup.ʔ():
		return ErrKickingChatOwner
	case ErrBotKickedFromSuperGroup.ʔ():
		return ErrBotKickedFromSuperGroup
	case ErrWrongTypeOfContent.ʔ():
		return ErrWrongTypeOfContent
	case ErrBadURLContent.ʔ():
		return ErrBadURLContent
	case ErrWrongFileIDSymbol.ʔ():
		return ErrWrongFileIDSymbol
	case ErrWrongFileIDLength.ʔ():
		return ErrWrongFileIDLength
	case ErrWrongFileIDCharacter.ʔ():
		return ErrWrongFileIDCharacter
	case ErrWrongFileID.ʔ():
		return ErrWrongFileID
	case ErrTooLarge.ʔ():
		return ErrTooLarge
	case ErrWrongFileIDPadding.ʔ():
		return ErrWrongFileIDPadding
	case ErrFailedImageProcess.ʔ():
		return ErrFailedImageProcess
	case ErrInvalidStickerSet.ʔ():
		return ErrInvalidStickerSet
	case ErrGroupMigrated.ʔ():
		return ErrGroupMigrated
	default:
		return nil
	}
}
