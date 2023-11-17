package telebot

type KeyboardButtonRequestChat struct {
	ID int32 `json:"request_id"`

	IsChannel       bool    `json:"chat_is_channel,omitempty"`
	IsForum         bool    `json:"chat_is_forum,omitempty"`
	HasUsername     bool    `json:"chat_has_username,omitempty"`
	IsCreated       bool    `json:"chat_is_created,omitempty"`
	UserAdminRights *Rights `json:"user_administrator_rights,omitempty"`
	BotAdminRights  *Rights `json:"bot_administrator_rights,omitempty"`
	BotIsMember     bool    `json:"bot_is_member,omitempty"`
}

type ChatShared struct {
	ID     int32 `json:"request_id"`
	ChatID int64 `json:"chat_id"`
}
