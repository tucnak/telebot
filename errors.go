package telebot

type Error struct {
	ErrorCode string
}

var (
	//bot creation errors
	ErrUnauthorized = &Error{"Unauthorized"}

	// not found etc
	ErrToForwardNotFound = &Error{"Bad Request: message to forward not found"}
	ErrToReplyNotFound   = &Error{"Bad Request: reply message not found"}
	ErrMessageTooLong    = &Error{"Bad Request: message is too long"}
	ErrBlockedByUsr      = &Error{"Forbidden: bot was blocked by the user"}
	ErrToDeleteNotFound  = &Error{"Bad Request: message to delete not found"}
	ErrEmptyMessage      = &Error{"Bad Request: message must be non-empty"}
	//checking
	ErrEmptyText          = &Error{"Bad Request: text is empty"}
	ErrEmptyChatID        = &Error{"Bad Request: chat_id is empty"}
	ErrNotFoundChat       = &Error{"Bad Request: chat not found"}
	ErrMessageNotModified = &Error{"Bad Request: message is not modified"}
	//

	// Rigts Errors
	ErrNoRightsToRestrict     = &Error{"Bad Request: not enough rights to restrict/unrestrict chat member"}
	ErrNoRightsToSendMsg      = &Error{"Bad Request: have no rights to send a message"}
	ErrNoRightsToSendPhoto    = &Error{"Bad Request: not enough rights to send photos to the chat"}
	ErrNoRightsToSendStickers = &Error{"Bad Request: not enough rights to send stickers to the chat"}
	ErrNoRightsToSendGifs     = &Error{"Bad Request: CHAT_SEND_GIFS_FORBIDDEN"}
	ErrNoRightsToDelete       = &Error{"Bad Request: message can't be deleted"}
	ErrKickingChatOwner       = &Error{"Bad Request: can't remove chat owner"}

	// Interacting with group/supergroup after being kicked
	ErrInteractKickedG    = &Error{"api error: Forbidden: bot was kicked from the group chat"}
	ErrInteractKickedSprG = &Error{"Forbidden: bot was kicked from the supergroup chat"}

	// file errors etc
	ErrWrongTypeOfContent      = &Error{"Bad Request: wrong type of the web page content"}
	ErrCantGetHTTPurlContent   = &Error{"Bad Request: failed to get HTTP URL content"}
	ErrWrongRemoteFileID       = &Error{"Bad Request: wrong remote file id specified: can't unserialize it. Wrong last symbol"}
	ErrWrongRemoteFileIDlen    = &Error{"Bad Request: wrong remote file id specified: Wrong string length"}
	ErrWrongRemoteFileIDsymbol = &Error{"Bad Request: wrong remote file id specified: Wrong character in the string"}
	ErrWrongFileIdentifier     = &Error{"Bad Request: wrong file identifier/HTTP URL specified"}
	ErrTooLarge                = &Error{"Request Entity Too Large"}
	ErrWrongPadding            = &Error{"Bad Request: wrong remote file id specified: Wrong padding in the string"} // not my
	ErrImageProcess            = &Error{"Bad Request: IMAGE_PROCESS_FAILED"}

	// sticker errors
	ErrWrongStickerpack = &Error{"Bad Request: STICKERSET_INVALID"}
)
