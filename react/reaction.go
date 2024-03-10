package react

// EmojiType defines emoji types.
type EmojiType = string

// Currently available emojis.
var (
	ThumbUp                   = Reaction{Emoji: "👍"}
	ThumbDown                 = Reaction{Emoji: "👎"}
	Heart                     = Reaction{Emoji: "❤"}
	Fire                      = Reaction{Emoji: "🔥"}
	HeartEyes                 = Reaction{Emoji: "😍"}
	ClappingHands             = Reaction{Emoji: "👏"}
	GrinningFace              = Reaction{Emoji: "😁"}
	ThinkingFace              = Reaction{Emoji: "🤔"}
	ExplodingHead             = Reaction{Emoji: "🤯"}
	ScreamingFace             = Reaction{Emoji: "😱"}
	SwearingFace              = Reaction{Emoji: "🤬"}
	CryingFace                = Reaction{Emoji: "😢"}
	PartyPopper               = Reaction{Emoji: "🎉"}
	StarStruck                = Reaction{Emoji: "🤩"}
	VomitingFace              = Reaction{Emoji: "🤮"}
	PileOfPoo                 = Reaction{Emoji: "💩"}
	PrayingHands              = Reaction{Emoji: "🙏"}
	OkHand                    = Reaction{Emoji: "👌"}
	DoveOfPeace               = Reaction{Emoji: "🕊"}
	ClownFace                 = Reaction{Emoji: "🤡"}
	YawningFace               = Reaction{Emoji: "🥱"}
	WoozyFace                 = Reaction{Emoji: "🥴"}
	Whale                     = Reaction{Emoji: "🐳"}
	HeartOnFire               = Reaction{Emoji: "❤‍🔥"}
	MoonFace                  = Reaction{Emoji: "🌚"}
	HotDog                    = Reaction{Emoji: "🌭"}
	HundredPoints             = Reaction{Emoji: "💯"}
	RollingOnTheFloorLaughing = Reaction{Emoji: "🤣"}
	Lightning                 = Reaction{Emoji: "⚡"}
	Banana                    = Reaction{Emoji: "🍌"}
	Trophy                    = Reaction{Emoji: "🏆"}
	BrokenHeart               = Reaction{Emoji: "💔"}
	FaceWithRaisedEyebrow     = Reaction{Emoji: "🤨"}
	NeutralFace               = Reaction{Emoji: "😐"}
	Strawberry                = Reaction{Emoji: "🍓"}
	Champagne                 = Reaction{Emoji: "🍾"}
	KissMark                  = Reaction{Emoji: "💋"}
	MiddleFinger              = Reaction{Emoji: "🖕"}
	EvilFace                  = Reaction{Emoji: "😈"}
	SleepingFace              = Reaction{Emoji: "😴"}
	LoudlyCryingFace          = Reaction{Emoji: "😭"}
	NerdFace                  = Reaction{Emoji: "🤓"}
	Ghost                     = Reaction{Emoji: "👻"}
	Engineer                  = Reaction{Emoji: "👨‍💻"}
	Eyes                      = Reaction{Emoji: "👀"}
	JackOLantern              = Reaction{Emoji: "🎃"}
	NoMonkey                  = Reaction{Emoji: "🙈"}
	SmilingFaceWithHalo       = Reaction{Emoji: "😇"}
	FearfulFace               = Reaction{Emoji: "😨"}
	Handshake                 = Reaction{Emoji: "🤝"}
	WritingHand               = Reaction{Emoji: "✍"}
	HuggingFace               = Reaction{Emoji: "🤗"}
	Brain                     = Reaction{Emoji: "🫡"}
	SantaClaus                = Reaction{Emoji: "🎅"}
	ChristmasTree             = Reaction{Emoji: "🎄"}
	Snowman                   = Reaction{Emoji: "☃"}
	NailPolish                = Reaction{Emoji: "💅"}
	ZanyFace                  = Reaction{Emoji: "🤪"}
	Moai                      = Reaction{Emoji: "🗿"}
	Cool                      = Reaction{Emoji: "🆒"}
	HeartWithArrow            = Reaction{Emoji: "💘"}
	HearMonkey                = Reaction{Emoji: "🙉"}
	Unicorn                   = Reaction{Emoji: "🦄"}
	FaceBlowingKiss           = Reaction{Emoji: "😘"}
	Pill                      = Reaction{Emoji: "💊"}
	SpeaklessMonkey           = Reaction{Emoji: "🙊"}
	Sunglasses                = Reaction{Emoji: "😎"}
	AlienMonster              = Reaction{Emoji: "👾"}
	ManShrugging              = Reaction{Emoji: "🤷‍♂️"}
	PersonShrugging           = Reaction{Emoji: "🤷"}
	WomanShrugging            = Reaction{Emoji: "🤷‍♀️"}
	PoutingFace               = Reaction{Emoji: "😡"}
)

// Reaction describes the type of reaction.
// Describes an instance of ReactionTypeCustomEmoji and ReactionTypeEmoji.
type Reaction struct {
	// Type of the reaction, always “emoji”
	Type string `json:"type"`

	// Reaction emoji.
	Emoji EmojiType `json:"emoji,omitempty"`

	// Custom emoji identifier.
	CustomEmoji string `json:"custom_emoji_id,omitempty"`
}

// Count represents a reaction added to a message along
// with the number of times it was added.
type Count struct {
	// Type of the reaction.
	Type Reaction `json:"type"`

	// Number of times the reaction was added.
	Count int `json:"total_count"`
}

// Options represents an object of reaction options.
type Options struct {
	// List of reaction types to set on the message.
	Reactions []Reaction `json:"reaction"`

	// Pass True to set the reaction with a big animation.
	Big bool `json:"is_big"`
}
