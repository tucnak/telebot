package telebot

import "io"

// API is the interface that wraps all basic methods for interacting
// with Telegram Bot API.
type API interface {
	Raw(method string, payload interface{}) ([]byte, error)
	RawFiles(method string, params map[string]string, files map[string]File) ([]byte, error)

	Accept(query *PreCheckoutQuery, errorMessage ...string) error
	AddStickerToSet(of Recipient, name string, sticker InputSticker) error
	AdminsOf(chat *Chat, returnBots ...bool) ([]ChatMember, error)
	Answer(query *Query, resp *QueryResponse) error
	AnswerGuest(msg *Message, result Result) (*SentGuestMessage, error)
	AnswerJoinRequest(req *ChatJoinRequest, result string) error
	AnswerWebApp(query *Query, r Result) (*WebAppMessage, error)
	AnswerWidget(queryID, data string) error
	ApproveSuggestedPost(msg Editable, sendDate int64) error
	ApproveJoinRequest(chat Recipient, user *User) error
	AvailableGifts() ([]Gift, error)
	Ban(chat *Chat, member *ChatMember, revokeMessages ...bool) error
	BanSenderChat(chat *Chat, sender Recipient) error
	BusinessConnection(id string) (*BusinessConnection, error)
	BusinessGifts(connectionID string, q GiftQuery) (*OwnedGifts, error)
	BusinessStarBalance(connectionID string) (*StarAmount, error)
	ChatByID(id int64) (*Chat, error)
	ChatByUsername(name string) (*Chat, error)
	ChatFullInfo(chat Recipient) (*ChatFullInfo, error)
	ChatGifts(chat Recipient, q GiftQuery) (*OwnedGifts, error)
	ChatMemberOf(chat, user Recipient) (*ChatMember, error)
	Close() (bool, error)
	CloseGeneralTopic(chat *Chat) error
	CloseTopic(chat *Chat, topic *Topic) error
	Commands(opts ...interface{}) ([]Command, error)
	ConvertGift(connectionID, ownedGiftID string) error
	Copy(to Recipient, msg Editable, opts ...interface{}) (*Message, error)
	CopyMany(to Recipient, msgs []Editable, opts ...*SendOptions) ([]Message, error)
	CreateInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error)
	CreateInvoiceLink(i Invoice) (string, error)
	CreateStickerSet(of Recipient, set *StickerSet) error
	CreateSubscriptionInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error)
	CreateTopic(chat *Chat, topic *Topic) (*Topic, error)
	CustomEmojiStickers(ids []string) ([]Sticker, error)
	DeclineJoinRequest(chat Recipient, user *User) error
	DeclineSuggestedPost(msg Editable, comment string) error
	DefaultRights(forChannels bool) (*Rights, error)
	Delete(msg Editable) error
	DeleteAllReactions(chat, actor Recipient) error
	DeleteBusinessGift(connectionID, ownedGiftID string) error
	DeleteBusinessMessages(connectionID string, msgs []Editable) error
	DeleteCommands(opts ...interface{}) error
	DeleteEphemeral(msg *Message) error
	DeleteGroupPhoto(chat *Chat) error
	DeleteGroupStickerSet(chat *Chat) error
	DeleteMany(msgs []Editable) error
	DeleteReaction(msg Editable, actor Recipient) error
	DeleteSticker(sticker string) error
	DeleteStickerSet(name string) error
	DeleteStory(connectionID string, storyID int) error
	DeleteTopic(chat *Chat, topic *Topic) error
	Download(file *File, localFilename string) error
	Edit(msg Editable, what interface{}, opts ...interface{}) (*Message, error)
	EditCaption(msg Editable, caption string, opts ...interface{}) (*Message, error)
	EditEphemeral(msg *Message, what interface{}, opts ...interface{}) error
	EditEphemeralCaption(msg *Message, caption string, opts ...interface{}) error
	EditEphemeralMedia(msg *Message, media Inputtable, opts ...interface{}) error
	EditEphemeralReplyMarkup(msg *Message, markup *ReplyMarkup) error
	EditGeneralTopic(chat *Chat, topic *Topic) error
	EditInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error)
	EditMedia(msg Editable, media Inputtable, opts ...interface{}) (*Message, error)
	EditReplyMarkup(msg Editable, markup *ReplyMarkup, opts ...interface{}) (*Message, error)
	EditSharedWidgetData(widgetID, publicData string, cacheTime int) error
	EditStory(connectionID string, storyID int, story StoryParams) (*Story, error)
	EditSubscriptionInviteLink(chat Recipient, link *ChatInviteLink) (*ChatInviteLink, error)
	EditTopic(chat *Chat, topic *Topic) error
	File(file *File) (io.ReadCloser, error)
	FileByID(fileID string) (File, error)
	Forward(to Recipient, msg Editable, opts ...interface{}) (*Message, error)
	ForwardMany(to Recipient, msgs []Editable, opts ...*SendOptions) ([]Message, error)
	GameScores(user Recipient, msg Editable) ([]GameHighScore, error)
	GiftPremiumSubscription(user Recipient, months, stars int, text string, opts ...interface{}) error
	HideGeneralTopic(chat *Chat) error
	InviteLink(chat *Chat) (string, error)
	Leave(chat Recipient) error
	Len(chat *Chat) (int, error)
	Logout() (bool, error)
	ManagedBotAccessSettings(bot *User) (*BotAccessSettings, error)
	ManagedBotToken(bot *User) (string, error)
	MenuButton(chat *User) (*MenuButton, error)
	MyDescription(language string) (*BotInfo, error)
	MyName(language string) (*BotInfo, error)
	MyShortDescription(language string) (*BotInfo, error)
	Notify(to Recipient, action ChatAction, opts ...interface{}) error
	PersonalChatMessages(user Recipient, limit int) ([]Message, error)
	Pin(msg Editable, opts ...interface{}) error
	PostStory(connectionID string, story StoryParams) (*Story, error)
	PrepareButton(user Recipient, button ReplyButton) (*PreparedKeyboardButton, error)
	PrepareMessage(user Recipient, result Result, opts ...interface{}) (*PreparedInlineMessage, error)
	ProfileAudiosOf(user *User, offset, limit int) (*UserProfileAudios, error)
	ProfilePhotosOf(user *User, page ...int) ([]Photo, error)
	Promote(chat *Chat, member *ChatMember) error
	React(to Recipient, msg Editable, r Reactions) error
	ReadBusinessMessage(connectionID string, msg Editable) error
	RefundStars(to Recipient, chargeID string) error
	EditUserStarSubscription(user Recipient, chargeID string, isCanceled bool) error
	EditWidgetData(msg Editable, publicData, chatData *string) error
	RemoveBusinessProfilePhoto(connectionID string, public bool) error
	RemoveChatVerification(chat *Chat) error
	RemoveMyProfilePhoto() error
	RemoveUserVerification(user Recipient) error
	RemoveWebhook(dropPending ...bool) error
	ReopenGeneralTopic(chat *Chat) error
	ReopenTopic(chat *Chat, topic *Topic) error
	ReplaceManagedBotToken(bot *User) (string, error)
	ReplaceStickerInSet(of Recipient, stickerSet, oldSticker string, sticker InputSticker) (bool, error)
	Reply(to *Message, what interface{}, opts ...interface{}) (*Message, error)
	RepostStory(connectionID string, from *Story, story StoryParams) (*Story, error)
	Respond(c *Callback, resp ...*CallbackResponse) error
	Restrict(chat *Chat, member *ChatMember) error
	RevokeInviteLink(chat Recipient, link string) (*ChatInviteLink, error)
	Send(to Recipient, what interface{}, opts ...interface{}) (*Message, error)
	SendAlbum(to Recipient, a Album, opts ...interface{}) ([]Message, error)
	SendDraft(to Recipient, draftID int, text string, opts ...interface{}) error
	SendJoinRequestWebApp(req *ChatJoinRequest, url string) error
	SendPaidMedia(to Recipient, stars int, a PaidAlbum, opts ...interface{}) (*Message, error)
	SendRichDraft(to Recipient, draftID int, rich *InputRichMessage, opts ...interface{}) error
	SendGift(to Recipient, giftID string, opts ...interface{}) error
	SetAdminTitle(chat *Chat, user *User, title string) error
	SetManagedBotAccessSettings(bot *User, settings BotAccessSettings) error
	SetMemberTag(chat *Chat, user *User, tag string) error
	SetBusinessBio(connectionID, bio string) error
	SetBusinessGiftSettings(connectionID string, showGiftButton bool, types AcceptedGiftTypes) error
	SetBusinessName(connectionID, firstName, lastName string) error
	SetBusinessProfilePhoto(connectionID string, photo InputProfilePhoto, public bool) error
	SetBusinessUsername(connectionID, username string) error
	SetCommands(opts ...interface{}) error
	SetCustomEmojiStickerSetThumb(name, id string) error
	SetDefaultRights(rights Rights, forChannels bool) error
	SetGameScore(user Recipient, msg Editable, score GameHighScore) (*Message, error)
	SetGroupDescription(chat *Chat, description string) error
	SetGroupPermissions(chat *Chat, perms Rights) error
	SetGroupPhoto(chat *Chat, p *Photo) error
	SetGroupStickerSet(chat *Chat, setName string) error
	SetGroupTitle(chat *Chat, title string) error
	SetMenuButton(chat *User, mb interface{}) error
	SetMyDescription(desc, language string) error
	SetMyName(name, language string) error
	SetMyProfilePhoto(photo InputProfilePhoto) error
	SetMyShortDescription(desc, language string) error
	SetPassportDataErrors(user Recipient, errs []PassportElementError) error
	SetStickerEmojis(sticker string, emojis []string) error
	SetStickerKeywords(sticker string, keywords []string) error
	SetStickerMaskPosition(sticker string, mask MaskPosition) error
	SetStickerPosition(sticker string, position int) error
	SetStickerSetThumb(of Recipient, set *StickerSet) error
	SetStickerSetTitle(s StickerSet) error
	SetUserEmojiStatus(user Recipient, emojiStatusCustomEmojiID string, expirationDate ...int64) error
	SetWebhook(w *Webhook) error
	ShareThemeBackdrops() ([]ChatShareThemeBackdrop, error)
	Ship(query *ShippingQuery, what ...interface{}) error
	StarBalance() (*StarAmount, error)
	StarTransactions(offset, limit int) ([]StarTransaction, error)
	StickerSet(name string) (*StickerSet, error)
	StopLiveLocation(msg Editable, opts ...interface{}) (*Message, error)
	StopPoll(msg Editable, opts ...interface{}) (*Poll, error)
	TopicIconStickers() ([]Sticker, error)
	TransferBusinessStars(connectionID string, stars int) error
	TransferGift(connectionID, ownedGiftID string, newOwner Recipient, stars int) error
	Unban(chat *Chat, user *User, forBanned ...bool) error
	UnbanSenderChat(chat *Chat, sender Recipient) error
	UniqueGift(name string) (*UniqueGift, error)
	UnhideGeneralTopic(chat *Chat) error
	Unpin(chat Recipient, opts ...interface{}) error
	UnpinAll(chat Recipient) error
	UnpinAllGeneralTopicMessages(chat *Chat) error
	UnpinAllTopicMessages(chat *Chat, topic *Topic) error
	UpgradeGift(connectionID, ownedGiftID string, keepDetails bool, stars int) error
	UploadSticker(to Recipient, format StickerSetFormat, f File) (*File, error)
	UserBoosts(chat, user Recipient) ([]Boost, error)
	UserGifts(user Recipient, q GiftQuery) (*OwnedGifts, error)
	Verify(to Recipient, description string, opts ...interface{}) error
	Webhook() (*Webhook, error)
}
