package telebot

import "io"

// API is the interface that wraps all basic methods for interacting
// with Telegram Bot API.
type API interface {
	Raw(method string, payload interface{}) ([]byte, error)

	Accept(query *PreCheckoutQuery, errorMessage ...string) error
	AddStickerToSet(of Recipient, name string, sticker InputSticker) error
	AdminsOf(chat *Chat) ([]ChatMember, error)
	Answer(query *Query, resp *QueryResponse) error
	AnswerWebApp(query *Query, r Result) (*WebAppMessage, error)
	ApproveJoinRequest(chat Recipient, user *User) error
	Ban(chat *Chat, member *ChatMember, revokeMessages ...bool) error
	BanSenderChat(chat *Chat, sender Recipient) error
	BusinessConnection(id string) (*BusinessConnection, error)
	ChatByID(id int64) (*Chat, error)
	ChatByUsername(name string) (*Chat, error)
	ChatMemberOf(chat, user Recipient) (*ChatMember, error)
	Close() (bool, error)
	CloseGeneralTopic(chat *Chat) error
	CloseTopic(chat *Chat, topic *Topic) error
	Commands(opts ...interface{}) ([]Command, error)
	Copy(to Recipient, msg Editable, opts ...interface{}) (*Message, error)
	CopyMany(to Recipient, msgs []Editable, opts ...*SendOptions) ([]Message, error)
	CreateInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error)
	CreateInvoiceLink(i Invoice) (string, error)
	CreateStickerSet(of Recipient, set *StickerSet) error
	CreateTopic(chat *Chat, topic *Topic) (*Topic, error)
	CustomEmojiStickers(ids []string) ([]Sticker, error)
	DeclineJoinRequest(chat Recipient, user *User) error
	DefaultRights(forChannels bool) (*Rights, error)
	Delete(msg Editable) error
	DeleteCommands(opts ...interface{}) error
	DeleteGroupPhoto(chat *Chat) error
	DeleteGroupStickerSet(chat *Chat) error
	DeleteMany(msgs []Editable) error
	DeleteSticker(sticker string) error
	DeleteStickerSet(name string) error
	DeleteTopic(chat *Chat, topic *Topic) error
	Download(file *File, localFilename string) error
	Edit(msg Editable, what interface{}, opts ...interface{}) (*Message, error)
	EditCaption(msg Editable, caption string, opts ...interface{}) (*Message, error)
	EditGeneralTopic(chat *Chat, topic *Topic) error
	EditInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error)
	EditMedia(msg Editable, media Inputtable, opts ...interface{}) (*Message, error)
	EditReplyMarkup(msg Editable, markup *ReplyMarkup) (*Message, error)
	EditTopic(chat *Chat, topic *Topic) error
	File(file *File) (io.ReadCloser, error)
	FileByID(fileID string) (File, error)
	Forward(to Recipient, msg Editable, opts ...interface{}) (*Message, error)
	ForwardMany(to Recipient, msgs []Editable, opts ...*SendOptions) ([]Message, error)
	GameScores(user Recipient, msg Editable) ([]GameHighScore, error)
	HideGeneralTopic(chat *Chat) error
	InviteLink(chat *Chat) (string, error)
	Leave(chat Recipient) error
	Len(chat *Chat) (int, error)
	Logout() (bool, error)
	MenuButton(chat *User) (*MenuButton, error)
	MyDescription(language string) (*BotInfo, error)
	MyName(language string) (*BotInfo, error)
	MyShortDescription(language string) (*BotInfo, error)
	Notify(to Recipient, action ChatAction, threadID ...int) error
	Pin(msg Editable, opts ...interface{}) error
	ProfilePhotosOf(user *User) ([]Photo, error)
	Promote(chat *Chat, member *ChatMember) error
	React(to Recipient, msg Editable, r Reactions) error
	RefundStars(to Recipient, chargeID string) error
	RemoveWebhook(dropPending ...bool) error
	ReopenGeneralTopic(chat *Chat) error
	ReopenTopic(chat *Chat, topic *Topic) error
	ReplaceStickerInSet(of Recipient, stickerSet, oldSticker string, sticker InputSticker) (bool, error)
	Reply(to *Message, what interface{}, opts ...interface{}) (*Message, error)
	Respond(c *Callback, resp ...*CallbackResponse) error
	Restrict(chat *Chat, member *ChatMember) error
	RevokeInviteLink(chat Recipient, link string) (*ChatInviteLink, error)
	Send(to Recipient, what interface{}, opts ...interface{}) (*Message, error)
	SendAlbum(to Recipient, a Album, opts ...interface{}) ([]Message, error)
	SendPaid(to Recipient, stars int, a PaidAlbum, opts ...interface{}) (*Message, error)
	SetAdminTitle(chat *Chat, user *User, title string) error
	SetCommands(opts ...interface{}) error
	SetCustomEmojiStickerSetThumb(name, id string) error
	SetDefaultRights(rights Rights, forChannels bool) error
	SetGameScore(user Recipient, msg Editable, score GameHighScore) (*Message, error)
	SetGroupDescription(chat *Chat, description string) error
	SetGroupPermissions(chat *Chat, perms Rights) error
	SetGroupStickerSet(chat *Chat, setName string) error
	SetGroupTitle(chat *Chat, title string) error
	SetMenuButton(chat *User, mb interface{}) error
	SetMyDescription(desc, language string) error
	SetMyName(name, language string) error
	SetMyShortDescription(desc, language string) error
	SetStickerEmojis(sticker string, emojis []string) error
	SetStickerKeywords(sticker string, keywords []string) error
	SetStickerMaskPosition(sticker string, mask MaskPosition) error
	SetStickerPosition(sticker string, position int) error
	SetStickerSetThumb(of Recipient, set *StickerSet) error
	SetStickerSetTitle(s StickerSet) error
	SetWebhook(w *Webhook) error
	Ship(query *ShippingQuery, what ...interface{}) error
	StarTransactions(offset, limit int) ([]StarTransaction, error)
	StickerSet(name string) (*StickerSet, error)
	StopLiveLocation(msg Editable, opts ...interface{}) (*Message, error)
	StopPoll(msg Editable, opts ...interface{}) (*Poll, error)
	TopicIconStickers() ([]Sticker, error)
	Unban(chat *Chat, user *User, forBanned ...bool) error
	UnbanSenderChat(chat *Chat, sender Recipient) error
	UnhideGeneralTopic(chat *Chat) error
	Unpin(chat Recipient, messageID ...int) error
	UnpinAll(chat Recipient) error
	UnpinAllGeneralTopicMessages(chat *Chat) error
	UnpinAllTopicMessages(chat *Chat, topic *Topic) error
	UploadSticker(to Recipient, format StickerSetFormat, f File) (*File, error)
	UserBoosts(chat, user Recipient) ([]Boost, error)
	Webhook() (*Webhook, error)
}
