package telebot

import (
	"fmt"
	"regexp"
	"strings"
)

type ApiError struct {
	Code 		int
	Description string

	message 	string

}

func (err *ApiError) ʔ() string {
	return err.Description
}

func badRequest(description ...string) *ApiError {
	if len(description) == 1 {
		return &ApiError{400, description[0], ""}
	}

	return &ApiError{400, description[0], description[1]}
}

func (err *ApiError) Error() string {
	canonical := err.message
	if canonical == ""{
		butchered := strings.Split(err.Description, ": ")
		if len(butchered) == 2{
			canonical = butchered[1]
		}else{
			canonical = err.Description
		}
	}
	return fmt.Sprintf("api error: %s (%d)", canonical, err.Code)
}

var (
	//RegExp
    apiErrorRx = regexp.MustCompile(`{.+"error_code":(\d+),"description":"(.+)"}`)

	//bot creation errors
	ErrUnauthorized = &ApiError{401,"Unauthorized",""}

	// not found etc
	ErrToForwardNotFound = badRequest("Bad Request: message to forward not found")
	ErrToReplyNotFound   = badRequest("Bad Request: reply message not found")
	ErrMessageTooLong    = badRequest("Bad Request: message is too long")
	ErrBlockedByUsr      = &ApiError{401,"Forbidden: bot was blocked by the user",""}
	ErrToDeleteNotFound  = badRequest("Bad Request: message to delete not found")
	ErrEmptyMessage      = badRequest("Bad Request: message must be non-empty")
	//checking
	ErrEmptyText          = badRequest("Bad Request: text is empty")
	ErrEmptyChatID        = badRequest("Bad Request: chat_id is empty")
	ErrNotFoundChat       = badRequest("Bad Request: chat not found")
	ErrMessageNotModified = badRequest("Bad Request: message is not modified")


	// Rigts Errors
	ErrNoRightsToRestrict     = badRequest("Bad Request: not enough rights to restrict/unrestrict chat member")
	ErrNoRightsToSendMsg      = badRequest("Bad Request: have no rights to send a message")
	ErrNoRightsToSendPhoto    = badRequest("Bad Request: not enough rights to send photos to the chat")
	ErrNoRightsToSendStickers = badRequest("Bad Request: not enough rights to send stickers to the chat")
	ErrNoRightsToSendGifs     = badRequest("Bad Request: CHAT_SEND_GIFS_FORBIDDEN","sending GIFS is not allowed in this chat")
	ErrNoRightsToDelete       = badRequest("Bad Request: message can't be deleted")
	ErrKickingChatOwner       = badRequest("Bad Request: can't remove chat owner")

	// Interacting with group/supergroup after being kicked
	ErrInteractKickedG    = &ApiError{403,"Forbidden: bot was kicked from the group chat",""}
	ErrInteractKickedSprG = &ApiError{403,"Forbidden: bot was kicked from the supergroup chat",""}

	// file errors etc
	ErrWrongTypeOfContent      = badRequest("Bad Request: wrong type of the web page content")
	ErrCantGetHTTPurlContent   = badRequest("Bad Request: failed to get HTTP URL content")
	ErrWrongRemoteFileID       = badRequest("Bad Request: wrong remote file id specified: can't unserialize it. Wrong last symbol")
	ErrFileIdTooShort    	   = badRequest("Bad Request: wrong remote file id specified: Wrong string length")
	ErrWrongRemoteFileIDsymbol = badRequest("Bad Request: wrong remote file id specified: Wrong character in the string")
	ErrWrongFileIdentifier     = badRequest("Bad Request: wrong file identifier/HTTP URL specified")
	ErrTooLarge                = badRequest("Request Entity Too Large")
	ErrWrongPadding            = badRequest("Bad Request: wrong remote file id specified: Wrong padding in the string") // not my
	ErrImageProcess            = badRequest("Bad Request: IMAGE_PROCESS_FAILED", "Image process failed")

	// sticker errors
	ErrWrongStickerpack = badRequest("Bad Request: STICKERSET_INVALID","Stickerset is invalid")
)