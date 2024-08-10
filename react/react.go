package react

import (
	tele "gopkg.in/telebot.v3"
)

type Reaction = tele.Reaction

func React(r ...Reaction) tele.Reactions {
	return tele.Reactions{Reactions: r}
}

// Currently available emojis.
var (
	ThumbUp                   = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👍"}
	ThumbDown                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👎"}
	Heart                     = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "❤"}
	Fire                      = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🔥"}
	HeartEyes                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😍"}
	ClappingHands             = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👏"}
	GrinningFace              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😁"}
	ThinkingFace              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤔"}
	ExplodingHead             = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤯"}
	ScreamingFace             = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😱"}
	SwearingFace              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤬"}
	CryingFace                = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😢"}
	PartyPopper               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🎉"}
	StarStruck                = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤩"}
	VomitingFace              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤮"}
	PileOfPoo                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "💩"}
	PrayingHands              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🙏"}
	OkHand                    = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👌"}
	DoveOfPeace               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🕊"}
	ClownFace                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤡"}
	YawningFace               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🥱"}
	WoozyFace                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🥴"}
	Whale                     = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🐳"}
	HeartOnFire               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "❤‍🔥"}
	MoonFace                  = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🌚"}
	HotDog                    = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🌭"}
	HundredPoints             = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "💯"}
	RollingOnTheFloorLaughing = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤣"}
	Lightning                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "⚡"}
	Banana                    = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🍌"}
	Trophy                    = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🏆"}
	BrokenHeart               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "💔"}
	FaceWithRaisedEyebrow     = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤨"}
	NeutralFace               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😐"}
	Strawberry                = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🍓"}
	Champagne                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🍾"}
	KissMark                  = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "💋"}
	MiddleFinger              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🖕"}
	EvilFace                  = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😈"}
	SleepingFace              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😴"}
	LoudlyCryingFace          = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😭"}
	NerdFace                  = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤓"}
	Ghost                     = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👻"}
	Engineer                  = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👨‍💻"}
	Eyes                      = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👀"}
	JackOLantern              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🎃"}
	NoMonkey                  = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🙈"}
	SmilingFaceWithHalo       = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😇"}
	FearfulFace               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😨"}
	Handshake                 = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤝"}
	WritingHand               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "✍"}
	HuggingFace               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤗"}
	Brain                     = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🫡"}
	SantaClaus                = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🎅"}
	ChristmasTree             = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🎄"}
	Snowman                   = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "☃"}
	NailPolish                = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "💅"}
	ZanyFace                  = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤪"}
	Moai                      = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🗿"}
	Cool                      = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🆒"}
	HeartWithArrow            = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "💘"}
	HearMonkey                = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🙉"}
	Unicorn                   = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🦄"}
	FaceBlowingKiss           = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😘"}
	Pill                      = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "💊"}
	SpeaklessMonkey           = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🙊"}
	Sunglasses                = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😎"}
	AlienMonster              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "👾"}
	ManShrugging              = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤷‍♂️"}
	PersonShrugging           = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤷"}
	WomanShrugging            = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "🤷‍♀️"}
	PoutingFace               = Reaction{Type: tele.ReactionTypeEmoji, Emoji: "😡"}
)
