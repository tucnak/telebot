package telebot

// Query is an incoming inline query. When the user sends
// an empty query, your bot could return some default or
// trending results.
type Query struct {
	ID   string `json:"id"`
	From User   `json:"from"`
	Text string `json:"query"`
}

// Result ...
type Result interface {
	MarshalJSON() ([]byte, error)
}
