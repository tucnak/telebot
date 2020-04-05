package telebot

import (
	"fmt"
	"regexp"
	"strings"
)

type APIError struct {
	Code        int
	Description string
	Message     string
}

// ʔ returns description of error.
// A tiny shortcut to make code clearier.
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
	return fmt.Sprintf("api error: %s (%d)", msg, err.Code)
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

var errorRx = regexp.MustCompile(`{.+"error_code":(\d+),"description":"(.+)"}`)

var (
	// Authorization errors
	ErrUnauthorized = NewAPIError(401, "Unauthorized")
	ErrBlockedByUsr = NewAPIError(401, "Forbidden: bot was blocked by the user")

	// Bad request errors
	ErrToForwardNotFound       = NewAPIError(400, "Bad Request: message to forward not found")
	ErrToReplyNotFound         = NewAPIError(400, "Bad Request: reply message not found")
	ErrMessageTooLong          = NewAPIError(400, "Bad Request: message is too long")
	ErrToDeleteNotFound        = NewAPIError(400, "Bad Request: message to delete not found")
	ErrEmptyMessage            = NewAPIError(400, "Bad Request: message must be non-empty")
	ErrNotFoundChat            = NewAPIError(400, "Bad Request: chat not found")
	ErrEmptyText               = NewAPIError(400, "Bad Request: text is empty")
	ErrEmptyChatID             = NewAPIError(400, "Bad Request: chat_id is empty")
	ErrMessageNotModified      = NewAPIError(400, "Bad Request: message is not modified")
	ErrWrongTypeOfContent      = NewAPIError(400, "Bad Request: wrong type of the web page content")
	ErrCantGetHTTPurlContent   = NewAPIError(400, "Bad Request: failed to get HTTP URL content")
	ErrWrongRemoteFileID       = NewAPIError(400, "Bad Request: wrong remote file id specified: can't unserialize it. Wrong last symbol")
	ErrFileIdTooShort          = NewAPIError(400, "Bad Request: wrong remote file id specified: Wrong string length")
	ErrWrongRemoteFileIDsymbol = NewAPIError(400, "Bad Request: wrong remote file id specified: Wrong character in the string")
	ErrWrongFileIdentifier     = NewAPIError(400, "Bad Request: wrong file identifier/HTTP URL specified")
	ErrWrongPadding            = NewAPIError(400, "Bad Request: wrong remote file id specified: Wrong padding in the string")
	ErrImageProcess            = NewAPIError(400, "Bad Request: IMAGE_PROCESS_FAILED", "Image process failed")
	ErrTooLarge                = NewAPIError(400, "Request Entity Too Large")
	ErrWrongStickerpack        = NewAPIError(400, "Bad Request: STICKERSET_INVALID", "Stickerset is invalid")

	// No rights errors
	ErrNoRightsToRestrict     = NewAPIError(400, "Bad Request: not enough rights to restrict/unrestrict chat member")
	ErrNoRightsToSendMsg      = NewAPIError(400, "Bad Request: have no rights to send a message")
	ErrNoRightsToSendPhoto    = NewAPIError(400, "Bad Request: not enough rights to send photos to the chat")
	ErrNoRightsToSendStickers = NewAPIError(400, "Bad Request: not enough rights to send stickers to the chat")
	ErrNoRightsToSendGifs     = NewAPIError(400, "Bad Request: CHAT_SEND_GIFS_FORBIDDEN", "sending GIFS is not allowed in this chat")
	ErrNoRightsToDelete       = NewAPIError(400, "Bad Request: message can't be deleted")
	ErrKickingChatOwner       = NewAPIError(400, "Bad Request: can't remove chat owner")

	// Super/groups errors
	ErrInteractKickedG    = NewAPIError(403, "Forbidden: bot was kicked from the group chat")
	ErrInteractKickedSprG = NewAPIError(403, "Forbidden: bot was kicked from the supergroup chat")
)

// errByDescription returns APIError instance by given description.
func errByDescription(s string) *APIError {
	switch s {
	case ErrUnauthorized.ʔ():
		return ErrUnauthorized
	case ErrToForwardNotFound.ʔ():
		return ErrToForwardNotFound
	case ErrToReplyNotFound.ʔ():
		return ErrToReplyNotFound
	case ErrMessageTooLong.ʔ():
		return ErrMessageTooLong
	case ErrBlockedByUsr.ʔ():
		return ErrBlockedByUsr
	case ErrToDeleteNotFound.ʔ():
		return ErrToDeleteNotFound
	case ErrEmptyMessage.ʔ():
		return ErrEmptyMessage
	case ErrEmptyText.ʔ():
		return ErrEmptyText
	case ErrEmptyChatID.ʔ():
		return ErrEmptyChatID
	case ErrNotFoundChat.ʔ():
		return ErrNotFoundChat
	case ErrMessageNotModified.ʔ():
		return ErrMessageNotModified
	case ErrNoRightsToRestrict.ʔ():
		return ErrNoRightsToRestrict
	case ErrNoRightsToSendMsg.ʔ():
		return ErrNoRightsToSendMsg
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
	case ErrInteractKickedG.ʔ():
		return ErrKickingChatOwner
	case ErrInteractKickedSprG.ʔ():
		return ErrInteractKickedSprG
	case ErrWrongTypeOfContent.ʔ():
		return ErrWrongTypeOfContent
	case ErrCantGetHTTPurlContent.ʔ():
		return ErrCantGetHTTPurlContent
	case ErrWrongRemoteFileID.ʔ():
		return ErrWrongRemoteFileID
	case ErrFileIdTooShort.ʔ():
		return ErrFileIdTooShort
	case ErrWrongRemoteFileIDsymbol.ʔ():
		return ErrWrongRemoteFileIDsymbol
	case ErrWrongFileIdentifier.ʔ():
		return ErrWrongFileIdentifier
	case ErrTooLarge.ʔ():
		return ErrTooLarge
	case ErrWrongPadding.ʔ():
		return ErrWrongPadding
	case ErrImageProcess.ʔ():
		return ErrImageProcess
	case ErrWrongStickerpack.ʔ():
		return ErrWrongStickerpack
	default:
		return nil
	}
}
