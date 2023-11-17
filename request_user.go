package telebot

type KeyboardButtonRequestUser struct {
	ID int32 `json:"request_id"`

	IsBot     bool `json:"user_is_bot,omitempty"`
	IsPremium bool `json:"user_is_premium,omitempty"`
}

type UserShared struct {
	ID     int32 `json:"request_id"`
	UserID int64 `json:"user_id"`
}
