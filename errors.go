package telebot

type ApiError struct {
	ErrorCode string
}

func (err *ApiError) Error() string {
	return err.ErrorCode
}

var (
	//bot creation errors
	ErrUnauthorized = &ApiError{"Unauthorized"}

	// not found etc
	ErrToForwardNotFound = &ApiError{"Bad Request: message to forward not found"}
	ErrToReplyNotFound   = &ApiError{"Bad Request: reply message not found"}
	ErrMessageTooLong    = &ApiError{"Bad Request: message is too long"}
	ErrBlockedByUsr      = &ApiError{"Forbidden: bot was blocked by the user"}
	ErrToDeleteNotFound  = &ApiError{"Bad Request: message to delete not found"}
	ErrEmptyMessage      = &ApiError{"Bad Request: message must be non-empty"}
	//checking
	ErrEmptyText          = &ApiError{"Bad Request: text is empty"}
	ErrEmptyChatID        = &ApiError{"Bad Request: chat_id is empty"}
	ErrNotFoundChat       = &ApiError{"Bad Request: chat not found"}
	ErrMessageNotModified = &ApiError{"Bad Request: message is not modified"}
	//

	// Rigts Errors
	ErrNoRightsToRestrict     = &ApiError{"Bad Request: not enough rights to restrict/unrestrict chat member"}
	ErrNoRightsToSendMsg      = &ApiError{"Bad Request: have no rights to send a message"}
	ErrNoRightsToSendPhoto    = &ApiError{"Bad Request: not enough rights to send photos to the chat"}
	ErrNoRightsToSendStickers = &ApiError{"Bad Request: not enough rights to send stickers to the chat"}
	ErrNoRightsToSendGifs     = &ApiError{"Bad Request: CHAT_SEND_GIFS_FORBIDDEN"}
	ErrNoRightsToDelete       = &ApiError{"Bad Request: message can't be deleted"}
	ErrKickingChatOwner       = &ApiError{"Bad Request: can't remove chat owner"}

	// Interacting with group/supergroup after being kicked
	ErrInteractKickedG    = &ApiError{"api error: Forbidden: bot was kicked from the group chat"}
	ErrInteractKickedSprG = &ApiError{"Forbidden: bot was kicked from the supergroup chat"}

	// file errors etc
	ErrWrongTypeOfContent      = &ApiError{"Bad Request: wrong type of the web page content"}
	ErrCantGetHTTPurlContent   = &ApiError{"Bad Request: failed to get HTTP URL content"}
	ErrWrongRemoteFileID       = &ApiError{"Bad Request: wrong remote file id specified: can't unserialize it. Wrong last symbol"}
	ErrWrongRemoteFileIDlen    = &ApiError{"Bad Request: wrong remote file id specified: Wrong string length"}
	ErrWrongRemoteFileIDsymbol = &ApiError{"Bad Request: wrong remote file id specified: Wrong character in the string"}
	ErrWrongFileIdentifier     = &ApiError{"Bad Request: wrong file identifier/HTTP URL specified"}
	ErrTooLarge                = &ApiError{"Request Entity Too Large"}
	ErrWrongPadding            = &ApiError{"Bad Request: wrong remote file id specified: Wrong padding in the string"} // not my
	ErrImageProcess            = &ApiError{"Bad Request: IMAGE_PROCESS_FAILED"}

	// sticker errors
	ErrWrongStickerpack = &ApiError{"Bad Request: STICKERSET_INVALID"}
)
