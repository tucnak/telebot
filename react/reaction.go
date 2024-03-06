package react

// EmojiType defines emoji types.
type EmojiType = string

// Currently available emojis.
var (
	ThumbUp                   = ReactionType{Emoji: "👍"}
	ThumbDown                 = ReactionType{Emoji: "👎"}
	Heart                     = ReactionType{Emoji: "❤"}
	Fire                      = ReactionType{Emoji: "🔥"}
	HeartEyes                 = ReactionType{Emoji: "😍"}
	ClappingHands             = ReactionType{Emoji: "👏"}
	GrinningFace              = ReactionType{Emoji: "😁"}
	ThinkingFace              = ReactionType{Emoji: "🤔"}
	ExplodingHead             = ReactionType{Emoji: "🤯"}
	ScreamingFace             = ReactionType{Emoji: "😱"}
	SwearingFace              = ReactionType{Emoji: "🤬"}
	CryingFace                = ReactionType{Emoji: "😢"}
	PartyPopper               = ReactionType{Emoji: "🎉"}
	StarStruck                = ReactionType{Emoji: "🤩"}
	VomitingFace              = ReactionType{Emoji: "🤮"}
	PileOfPoo                 = ReactionType{Emoji: "💩"}
	PrayingHands              = ReactionType{Emoji: "🙏"}
	OkHand                    = ReactionType{Emoji: "👌"}
	DoveOfPeace               = ReactionType{Emoji: "🕊"}
	ClownFace                 = ReactionType{Emoji: "🤡"}
	YawningFace               = ReactionType{Emoji: "🥱"}
	WoozyFace                 = ReactionType{Emoji: "🥴"}
	Whale                     = ReactionType{Emoji: "🐳"}
	HeartOnFire               = ReactionType{Emoji: "❤‍🔥"}
	MoonFace                  = ReactionType{Emoji: "🌚"}
	HotDog                    = ReactionType{Emoji: "🌭"}
	HundredPoints             = ReactionType{Emoji: "💯"}
	RollingOnTheFloorLaughing = ReactionType{Emoji: "🤣"}
	Lightning                 = ReactionType{Emoji: "⚡"}
	Banana                    = ReactionType{Emoji: "🍌"}
	Trophy                    = ReactionType{Emoji: "🏆"}
	BrokenHeart               = ReactionType{Emoji: "💔"}
	FaceWithRaisedEyebrow     = ReactionType{Emoji: "🤨"}
	NeutralFace               = ReactionType{Emoji: "😐"}
	Strawberry                = ReactionType{Emoji: "🍓"}
	Champagne                 = ReactionType{Emoji: "🍾"}
	KissMark                  = ReactionType{Emoji: "💋"}
	MiddleFinger              = ReactionType{Emoji: "🖕"}
	EvilFace                  = ReactionType{Emoji: "😈"}
	SleepingFace              = ReactionType{Emoji: "😴"}
	LoudlyCryingFace          = ReactionType{Emoji: "😭"}
	NerdFace                  = ReactionType{Emoji: "🤓"}
	Ghost                     = ReactionType{Emoji: "👻"}
	Engineer                  = ReactionType{Emoji: "👨‍💻"}
	Eyes                      = ReactionType{Emoji: "👀"}
	JackOLantern              = ReactionType{Emoji: "🎃"}
	NoMonkey                  = ReactionType{Emoji: "🙈"}
	SmilingFaceWithHalo       = ReactionType{Emoji: "😇"}
	FearfulFace               = ReactionType{Emoji: "😨"}
	Handshake                 = ReactionType{Emoji: "🤝"}
	WritingHand               = ReactionType{Emoji: "✍"}
	HuggingFace               = ReactionType{Emoji: "🤗"}
	Brain                     = ReactionType{Emoji: "🫡"}
	SantaClaus                = ReactionType{Emoji: "🎅"}
	ChristmasTree             = ReactionType{Emoji: "🎄"}
	Snowman                   = ReactionType{Emoji: "☃"}
	NailPolish                = ReactionType{Emoji: "💅"}
	ZanyFace                  = ReactionType{Emoji: "🤪"}
	Moai                      = ReactionType{Emoji: "🗿"}
	Cool                      = ReactionType{Emoji: "🆒"}
	HeartWithArrow            = ReactionType{Emoji: "💘"}
	HearMonkey                = ReactionType{Emoji: "🙉"}
	Unicorn                   = ReactionType{Emoji: "🦄"}
	FaceBlowingKiss           = ReactionType{Emoji: "😘"}
	Pill                      = ReactionType{Emoji: "💊"}
	SpeaklessMonkey           = ReactionType{Emoji: "🙊"}
	Sunglasses                = ReactionType{Emoji: "😎"}
	AlienMonster              = ReactionType{Emoji: "👾"}
	ManShrugging              = ReactionType{Emoji: "🤷‍♂️"}
	PersonShrugging           = ReactionType{Emoji: "🤷"}
	WomanShrugging            = ReactionType{Emoji: "🤷‍♀️"}
	PoutingFace               = ReactionType{Emoji: "😡"}
)

// ReactionType describes the type of reaction.
// Describes an instance of ReactionTypeCustomEmoji and ReactionTypeEmoji.
type ReactionType struct {
	// Type of the reaction, always “emoji”
	Type string `json:"type"`

	// Reaction emoji.
	Emoji EmojiType `json:"emoji,omitempty"`

	// Custom emoji identifier.
	CustomEmoji string `json:"custom_emoji_id,omitempty"`
}

// ReactionCount represents a reaction added to a message along
// with the number of times it was added.
type ReactionCount struct {
	// Type of the reaction.
	Type ReactionType `json:"type"`

	// Number of times the reaction was added.
	Count int `json:"total_count"`
}

// ReactionOptions represents an object of reaction options.
type ReactionOptions struct {
	// List of reaction types to set on the message.
	Reactions []ReactionType `json:"reaction"`

	// Pass True to set the reaction with a big animation.
	IsBig bool `json:"is_big"`
}
