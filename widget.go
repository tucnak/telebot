package telebot

import (
	"strconv"
)

// Widget is a message with an HTML widget (Bot API 10.4). As received, it
// describes the widget of Message.Widget. As sent, it is a Sendable:
//
//	b.Send(user, &tele.Widget{HTML: "<h1>Hi</h1>", Height: 200, AllowQueries: true})
//
// A Widget with an ID and no HTML re-sends that shared widget.
//
// In production, widgets can be sent only to private chats, and are always
// protected from forwarding.
//
// Bot API: Widget.
type Widget struct {
	// Identifier of a shared widget; empty if the widget isn't shared.
	ID string `json:"id,omitempty"`

	// (Optional) Size in logical pixels; zero width fills the message bubble,
	// zero height lets the widget grow.
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	// HTML code of the widget, 1-32768 bytes and at most 4096 bytes when
	// compressed with zlib.
	HTML string `json:"html,omitempty"`

	// (Optional) Data available wherever the widget is shown; all messages
	// with a shared widget have the same public data.
	PublicData string `json:"public_data,omitempty"`

	// (Optional) Data specific to the message.
	ChatData string `json:"chat_data,omitempty"`

	// The fields below apply when sending.

	// (Optional) Caption of the widget message; parse mode and entities come
	// from the send options.
	Caption string `json:"-"`

	// (Optional) Show the caption above the widget.
	CaptionAbove bool `json:"-"`

	// (Optional) Allow the widget to send queries to the bot.
	AllowQueries bool `json:"-"`

	// (Optional) Create a shared widget, which can be sent again by ID.
	IsShared bool `json:"-"`
}

// Send delivers the widget through bot b to recipient.
//
// Bot API: sendWidget, sendSharedWidget.
func (w *Widget) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}
	if w.Caption != "" {
		params["caption"] = w.Caption
	}
	if w.ChatData != "" {
		params["chat_data"] = w.ChatData
	}
	if w.CaptionAbove {
		params["show_caption_above_media"] = "true"
	}

	method := "sendWidget"
	if w.HTML == "" && w.ID != "" {
		method = "sendSharedWidget"
		params["widget_id"] = w.ID
	} else {
		params["html"] = w.HTML
		if w.Width != 0 {
			params["width"] = strconv.Itoa(w.Width)
		}
		if w.Height != 0 {
			params["height"] = strconv.Itoa(w.Height)
		}
		if w.PublicData != "" {
			params["public_data"] = w.PublicData
		}
		if w.AllowQueries {
			params["allow_queries"] = "true"
		}
		if w.IsShared {
			params["is_shared"] = "true"
		}
	}
	b.embedSendOptions(params, opt)

	data, err := b.Raw(method, params)
	if err != nil {
		return nil, err
	}
	return extractMessage(data)
}

// Widget query types.
const (
	WidgetQueryData            = "data"
	WidgetQuerySignData        = "sign_data"
	WidgetQuerySentTransaction = "sent_transaction"
)

// WidgetQuery is a query from a widget in a message sent by or via the bot
// (Bot API 10.4), delivered to OnWidgetQuery. Answer it with
// AnswerWidget within 15 seconds.
//
// Bot API: WidgetQuery.
type WidgetQuery struct {
	ID     string `json:"id"`
	Sender *User  `json:"from"`

	// (Optional) Chat and message with the widget; omitted for inline
	// messages, which have InlineMessageID instead.
	Chat            *Chat  `json:"chat,omitempty"`
	MessageID       int    `json:"message_id,omitempty"`
	InlineMessageID string `json:"inline_message_id,omitempty"`

	// Global identifier of the chat the widget's message was sent to.
	ChatInstance string `json:"chat_instance"`

	// (Optional) For inline messages, the type of the chat: "sender",
	// "private", "group", "supergroup" or "channel".
	ChatType string `json:"chat_type,omitempty"`

	// What the query carries.
	Type WidgetQueryType `json:"type"`

	// Data of the query, 1-64 bytes.
	Data string `json:"data"`

	// (Optional) The shared widget that sent the query.
	WidgetID string `json:"widget_id,omitempty"`
}

// MessageSig lets a query from a widget be passed to EditWidgetData.
func (q *WidgetQuery) MessageSig() (string, int64) {
	if q.InlineMessageID != "" || q.Chat == nil {
		return q.InlineMessageID, 0
	}
	return strconv.Itoa(q.MessageID), q.Chat.ID
}

// SharedWidgetQuery is a query from a forwarded message with a shared widget
// created by the bot (Bot API 10.4), delivered to OnSharedWidgetQuery. Answer
// it with AnswerWidget.
//
// Bot API: SharedWidgetQuery.
type SharedWidgetQuery struct {
	ID       string          `json:"id"`
	Sender   *User           `json:"from"`
	Type     WidgetQueryType `json:"type"`
	Data     string          `json:"data"`
	WidgetID string          `json:"widget_id"`
}

// WidgetQueryType tells what a widget query carries: data for the bot,
// possibly with the user's TON wallet card ("data"), data signed by the
// user's TON wallet ("sign_data"), or a transaction sent from it
// ("sent_transaction").
//
// Bot API: WidgetQueryType, WidgetQueryTypeData, WidgetQueryTypeSignData, WidgetQueryTypeSentTransaction.
type WidgetQueryType struct {
	// WidgetQueryData, WidgetQuerySignData or WidgetQuerySentTransaction.
	Type string `json:"type"`

	// (Optional) For "data", the user's TON wallet card, if it was sent.
	WalletCard *TonWalletCard `json:"wallet_card,omitempty"`

	// For "sign_data" and "sent_transaction", the user's wallet.
	Wallet *WidgetQueryWallet `json:"wallet,omitempty"`

	// For "sign_data", the signed data, JSON-serialized. The data is signed
	// for the domain bot<bot_id>.
	SignData string `json:"sign_data,omitempty"`

	// For "sent_transaction", the msg_hash of the transaction.
	MsgHash string `json:"msg_hash,omitempty"`
}

// WidgetQueryWallet is a user's TON wallet in a widget query.
//
// Bot API: WidgetQueryWallet.
type WidgetQueryWallet struct {
	Address string `json:"address"`

	// Balance in the smallest units of the cryptocurrency, as a decimal
	// number.
	Amount string `json:"amount"`

	// (Optional) Hex-encoded public key and base64-encoded StateInit of the
	// wallet contract; both omitted if unknown.
	PublicKey string `json:"public_key,omitempty"`
	StateInit string `json:"state_init,omitempty"`
}

// AnswerWidget answers a WidgetQuery or SharedWidgetQuery by its ID.
// Only the widget that sent the query receives data, 0-1024 bytes; the
// widget's own data is not changed.
//
// Bot API: answerWidgetQuery.
func (b *Bot) AnswerWidget(queryID, data string) error {
	params := map[string]string{
		"widget_query_id": queryID,
	}
	if data != "" {
		params["data"] = data
	}
	_, err := b.Raw("answerWidgetQuery", params)
	return err
}

// EditWidgetData changes the data of the widget in msg, a message sent by or
// via the bot. A nil publicData or chatData is left unchanged; an empty one
// is removed.
//
// Bot API: editMessageWidgetData.
func (b *Bot) EditWidgetData(msg Editable, publicData, chatData *string) error {
	msgID, chatID := msg.MessageSig()

	params := map[string]string{}
	if chatID == 0 {
		params["inline_message_id"] = msgID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = msgID
	}
	if publicData != nil {
		params["public_data"] = *publicData
	}
	if chatData != nil {
		params["chat_data"] = *chatData
	}

	_, err := b.Raw("editMessageWidgetData", params)
	return err
}

// EditSharedWidgetData changes the public data of a shared widget created by
// the bot; an empty publicData removes it. cacheTime, in seconds (5-3600), is
// how long clients may cache the data; zero keeps the current value.
//
// Bot API: editSharedWidgetData.
func (b *Bot) EditSharedWidgetData(widgetID, publicData string, cacheTime int) error {
	params := map[string]string{
		"widget_id":   widgetID,
		"public_data": publicData,
	}
	if cacheTime != 0 {
		params["cache_time"] = strconv.Itoa(cacheTime)
	}

	_, err := b.Raw("editSharedWidgetData", params)
	return err
}

// InputWidgetMessageContent is a widget to send as the result of an inline
// query (Bot API 10.4). Telegram does not accept it in production yet.
//
// Bot API: InputWidgetMessageContent.
type InputWidgetMessageContent struct {
	HTML         string          `json:"html"`
	Width        int             `json:"width,omitempty"`
	Height       int             `json:"height,omitempty"`
	PublicData   string          `json:"public_data,omitempty"`
	ChatData     string          `json:"chat_data,omitempty"`
	Caption      string          `json:"caption,omitempty"`
	ParseMode    ParseMode       `json:"parse_mode,omitempty"`
	Entities     []MessageEntity `json:"caption_entities,omitempty"`
	CaptionAbove bool            `json:"show_caption_above_media,omitempty"`
	AllowQueries bool            `json:"allow_queries,omitempty"`
	IsShared     bool            `json:"is_shared,omitempty"`
}

// IsInputMessageContent marks InputWidgetMessageContent as InputMessageContent.
func (input *InputWidgetMessageContent) IsInputMessageContent() bool {
	return true
}

// InputSharedWidgetMessageContent is a shared widget created by the bot to
// send as the result of an inline query (Bot API 10.4). Telegram does not
// accept it in production yet.
//
// Bot API: InputSharedWidgetMessageContent.
type InputSharedWidgetMessageContent struct {
	WidgetID     string          `json:"widget_id"`
	ChatData     string          `json:"chat_data,omitempty"`
	Caption      string          `json:"caption,omitempty"`
	ParseMode    ParseMode       `json:"parse_mode,omitempty"`
	Entities     []MessageEntity `json:"caption_entities,omitempty"`
	CaptionAbove bool            `json:"show_caption_above_media,omitempty"`
}

// IsInputMessageContent marks InputSharedWidgetMessageContent as
// InputMessageContent.
func (input *InputSharedWidgetMessageContent) IsInputMessageContent() bool {
	return true
}
