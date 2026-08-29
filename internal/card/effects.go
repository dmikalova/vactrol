package card

import "github.com/dmikalova/vactrol/internal/engine"

// The card effect "AST" re-exported for authoring: the effect nodes you nest
// inside an ability, plus the CreatureVerb, Condition, and Count helpers some of
// them take. Each mirrors a type in the engine's effect_*.go — see there for what
// it does. They are named directly, e.g. card.DealDamage{Amount: 2, ...}.
type (
	// Effect is a card effect (see DealDamage, GainAember, ...).
	Effect = engine.Effect
	// CreatureVerb is an action applied to a chosen creature (ReadyVerb, FightVerb).
	CreatureVerb = engine.CreatureVerb
	// Condition gates a Conditional effect (see OpponentAemberAtLeast).
	Condition = engine.Condition

	DealDamage                = engine.DealDamage
	SplashDamage              = engine.SplashDamage
	Exalt                     = engine.Exalt
	GainAember                = engine.GainAember
	LoseAember                = engine.LoseAember
	StealAember               = engine.StealAember
	CaptureAember             = engine.CaptureAember
	Draw                      = engine.Draw
	Heal                      = engine.Heal
	CreaturesHealed           = engine.CreaturesHealed
	Destroy                   = engine.Destroy
	DestroySamePower          = engine.DestroySamePower
	Then                      = engine.Then
	May                       = engine.May
	Purge                     = engine.Purge
	PurgeFromHand             = engine.PurgeFromHand
	PurgeCreature             = engine.PurgeCreature
	AddPowerCounter           = engine.AddPowerCounter
	EachPlayerLosesAllBut     = engine.EachPlayerLosesAllBut
	EachPlayerLosesHalfAember = engine.EachPlayerLosesHalfAember
	GainChains                = engine.GainChains
	Stun                      = engine.Stun
	Unstun                    = engine.Unstun
	Exhaust                   = engine.Exhaust
	OnChooseCreature          = engine.OnChooseCreature
	Sequence                  = engine.Sequence
	ChooseOne                 = engine.ChooseOne
	ChooseHouseThen           = engine.ChooseHouseThen
	Conditional               = engine.Conditional
	RepeatWhile               = engine.RepeatWhile
	OpponentAemberAtLeast     = engine.OpponentAemberAtLeast
	OpponentAemberExactly     = engine.OpponentAemberExactly
	OpponentAemberMoreThanYou = engine.OpponentAemberMoreThanYou
	ReadyVerb                 = engine.ReadyVerb
	FightVerb                 = engine.FightVerb
	UseVerb                   = engine.UseVerb
	StunVerb                  = engine.StunVerb
	ExhaustVerb               = engine.ExhaustVerb
	ReadyCreatures            = engine.ReadyCreatures
	MoveFromPlay              = engine.MoveFromPlay
	MoveArtifactsToHand       = engine.MoveArtifactsToHand
	ArchiveFromHand           = engine.ArchiveFromHand
	ArchiveTopOfDeck          = engine.ArchiveTopOfDeck
	DiscardArchives           = engine.DiscardArchives
	MoveFromDiscard           = engine.MoveFromDiscard
	DiscardHand               = engine.DiscardHand
	Reveal                    = engine.Reveal
	CannotFight               = engine.CannotFight
	GrantFightForChosenHouse  = engine.GrantFightForChosenHouse
	ForceOpponentActiveHouse  = engine.ForceOpponentActiveHouse
	ForRemainderOfTurn        = engine.ForRemainderOfTurn
	Instead                   = engine.Instead
	ForgeKey                  = engine.ForgeKey
	OpponentForgedKeys        = engine.OpponentForgedKeys
	CardsInArchives           = engine.CardsInArchives
	FriendlyCreaturesInPlay   = engine.FriendlyCreaturesInPlay
	CardsRevealed             = engine.CardsRevealed
)

// Event groups the game events a lasting "for the remainder of the turn" effect
// attaches to (see ForRemainderOfTurn and Instead), e.g.
// card.ForRemainderOfTurn{On: card.Event.CreaturePlayed, Do: card.GainAember{...}}.
var Event = events{
	CreaturePlayed: engine.EventCreaturePlayed,
	Reap:           engine.EventReap,
	ReapAember:     engine.EventReapAember,
}

type events struct {
	CreaturePlayed, Reap, ReapAember engine.Event
}

// Steal is the replacement that makes gaining Æmber steal it from the opponent
// instead, for card.Instead (Dimension Door).
var Steal = engine.Steal
