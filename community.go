package telebot

// Community is a group of supergroups, channels and bots linked together
// around a shared topic or audience (Bot API 10.2).
//
// Bot API: Community.
type Community struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// CommunityChatAdded is a service message about a chat or a bot being added
// to a community.
//
// Bot API: CommunityChatAdded.
type CommunityChatAdded struct {
	// The community the chat or the bot now belongs to.
	Community Community `json:"community"`
}

// CommunityChatRemoved is a service message about a chat or a bot being
// removed from a community. It holds no information yet.
//
// Bot API: CommunityChatRemoved.
type CommunityChatRemoved struct{}

// CommunityChatJoined is a service message about a user joining the chat
// from a community (Bot API 10.3).
//
// Bot API: CommunityChatJoined.
type CommunityChatJoined struct {
	// The community from which the chat was joined.
	Community Community `json:"community"`
}
