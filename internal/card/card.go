// Package card is a small, ergonomic facade over the game engine for *authoring*
// cards. Card-set packages import this package (as `card`) so definitions read
// like card.New(...), card.House.Brobnar, card.DealDamage{...} without referring to the
// engine's package name directly.
//
// Almost everything here is a re-export — type aliases (which are fully
// interchangeable with the engine's types) and constant/function values — of
// identifiers in package engine, shaping the authoring API. The one bit of logic
// is the card registry (see register.go): New builds a card and enrolls it in
// the database, so a set package declares each card as `var X = card.New(...)`
// — one line, no init, no central list.
package card

import (
	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Types you name directly when authoring cards (structs and interfaces). The
// enum-like categories (house, card type, rarity, keyword, trigger, controller,
// and target) are exposed as grouped namespaces further down, not as loose constants.
type (
	// Definition is a card blueprint.
	Definition = engine.CardDefinition
	// Ability pairs a trigger with an effect.
	Ability = engine.Ability
	// StaticModifier is a continuous stat change applied by an upgrade.
	StaticModifier = engine.StaticModifier
	// Trait is a freeform label such as "Giant" or "Weapon".
	Trait = engine.Trait

	// Effect is a card effect (see DealDamage, GainAember, ...).
	Effect = engine.Effect
	// CreatureVerb is an action applied to a chosen creature (ReadyVerb, FightVerb).
	CreatureVerb = engine.CreatureVerb
	// Condition gates a Conditional effect (see OpponentAemberAtLeast).
	Condition = engine.Condition
	// Player is the relative player an effect targets: card.Controller or card.Opponent.
	Player = engine.Player

	DealDamage            = engine.DealDamage
	Exalt                 = engine.Exalt
	GainAember            = engine.GainAember
	LoseAember            = engine.LoseAember
	StealAember           = engine.StealAember
	CaptureAember         = engine.CaptureAember
	Draw                  = engine.Draw
	Heal                  = engine.Heal
	Destroy               = engine.Destroy
	Stun                  = engine.Stun
	Unstun                = engine.Unstun
	OnChosenCreature      = engine.OnChosenCreature
	Sequence              = engine.Sequence
	Conditional           = engine.Conditional
	OpponentAemberAtLeast = engine.OpponentAemberAtLeast
	OpponentAemberExactly = engine.OpponentAemberExactly
	ReadyVerb             = engine.ReadyVerb
	FightVerb             = engine.FightVerb
	ReturnToDeck          = engine.ReturnToDeck
	ReturnToHand          = engine.ReturnToHand
)

// Option configures a card at construction: either a gameplay option (the With*
// helpers below, which wrap the engine's options) or a Provenance tag. New splits
// them — gameplay options build the definition, provenance is recorded in the
// registry (see register.go).
type Option func(*builder)

// builder accumulates the pieces New assembles into a registered card.
type builder struct {
	opts []engine.CardOption
	prov []provenance.Ref
}

func gameplay(o engine.CardOption) Option { return func(b *builder) { b.opts = append(b.opts, o) } }

// Option helpers (thin wrappers over the engine's options). New lives in
// register.go because building a card there also enrolls it in the database.
var (
	WithPower       = func(p int) Option { return gameplay(engine.WithPower(p)) }
	WithArmor       = func(a int) Option { return gameplay(engine.WithArmor(a)) }
	WithTraits      = func(t ...Trait) Option { return gameplay(engine.WithTraits(t...)) }
	WithKeywords    = func(k ...engine.Keyword) Option { return gameplay(engine.WithKeywords(k...)) }
	WithAemberBonus = func(n int) Option { return gameplay(engine.WithAemberBonus(n)) }
	WithStatic      = func(m StaticModifier) Option { return gameplay(engine.WithStatic(m)) }
	WithAbility     = func(t engine.Trigger, e Effect) Option { return gameplay(engine.WithAbility(t, e)) }
)

// Provenance tags a card as derived from an original source card (set + collector
// number), for coverage tracking. Optional and repeatable — a card may draw from
// more than one original — e.g. card.Provenance(card.CotA, 1).
func Provenance(set provenance.SourceSet, number int) Option {
	return func(b *builder) { b.prov = append(b.prov, provenance.Ref{Set: set, Number: number}) }
}

// Source sets to tag a card's Provenance with, e.g. card.Provenance(card.CotA, 1).
var (
	CotA = provenance.CallOfTheArchons
	AoA  = provenance.AgeOfAscension
	WC   = provenance.WorldsCollide
	MM   = provenance.MassMutation
	DT   = provenance.DarkTidings
	WoE  = provenance.WindsOfExchange
	GR   = provenance.GrimReminders
	AS   = provenance.AemberSkies
	ToC  = provenance.TokensOfChange
	MoM  = provenance.MoreMutation
	Men  = provenance.Menagerie
	VM   = provenance.VaultMasters2025
	PV   = provenance.PropheticVisions
	CC   = provenance.CrucibleClash
	DM   = provenance.DraconianMeasures
)

// The enum-like categories are exposed as grouped namespaces so related values
// stay together and read unambiguously — card.House.Brobnar is plainly a house,
// card.Type.Action plainly a card type — all from this single `card` import.
// Treat these package-level vars as read-only.

// House groups the faction values, e.g. card.House.Brobnar.
var House = houses{
	Brobnar: engine.Brobnar,
	Dis:     engine.Dis,
	Logos:   engine.Logos,
	Mars:    engine.Mars,
	Sanctum: engine.Sanctum,
	Shadows: engine.Shadows,
	Untamed: engine.Untamed,
}

type houses struct {
	Brobnar, Dis, Logos, Mars, Sanctum, Shadows, Untamed engine.House
}

// Type groups the card-type values, e.g. card.Type.Creature.
var Type = cardTypes{
	Creature: engine.Creature,
	Action:   engine.Action,
	Artifact: engine.Artifact,
	Upgrade:  engine.Upgrade,
}

type cardTypes struct {
	Creature, Action, Artifact, Upgrade engine.CardType
}

// Rarity groups the rarity values, e.g. card.Rarity.Common.
var Rarity = rarities{
	Common:   engine.Common,
	Uncommon: engine.Uncommon,
	Rare:     engine.Rare,
	Special:  engine.Special,
}

type rarities struct {
	Common, Uncommon, Rare, Special engine.Rarity
}

// Keyword groups the keyword values, e.g. card.Keyword.Skirmish.
var Keyword = keywords{
	Skirmish: engine.Skirmish,
	Poison:   engine.Poison,
	Elusive:  engine.Elusive,
	Taunt:    engine.Taunt,
}

type keywords struct {
	Skirmish, Poison, Elusive, Taunt engine.Keyword
}

// Trigger groups the ability triggers, e.g. card.Trigger.Play or
// card.Trigger.AfterForgeKey.
var Trigger = triggers{
	Play:                engine.TriggerAfterPlay,
	Reap:                engine.TriggerAfterReap,
	Fight:               engine.TriggerAfterFight,
	BeforeFight:         engine.TriggerBeforeFight,
	Action:              engine.TriggerAction,
	AfterForgeKey:       engine.TriggerAfterForgeKey,
	AfterCreatureEnters: engine.TriggerAfterCreatureEnters,
	Destroyed:           engine.TriggerDestroyed,
}

type triggers struct {
	Play,
	Reap,
	Fight,
	BeforeFight,
	Action,
	AfterForgeKey,
	AfterCreatureEnters,
	Destroyed engine.Trigger
}

// Controller and Opponent are the two players an effect can target,
//
//	relative to the card's controller: card.Controller (the player who controls the card) or
//
// card.Opponent (their opponent).
var (
	Controller = engine.Controller
	Opponent   = engine.Opponent
)

// Target groups ready-made targets, e.g. card.Target.EachEnemyCreature.
var Target = targets{
	This:                      engine.Target{Kind: engine.TargetThisCreature},
	Triggering:                engine.Target{Kind: engine.TargetTriggeringCreature},
	Creature:                  engine.Target{Kind: engine.TargetChosenCreature},
	EachCreature:              engine.Target{Kind: engine.TargetEachCreature},
	EachFriendlyCreature:      engine.Target{Kind: engine.TargetEachFriendlyCreature},
	EachEnemyCreature:         engine.Target{Kind: engine.TargetEachEnemyCreature},
	EachArtifact:              engine.Target{Kind: engine.TargetEachArtifact},
	EachOtherFriendlyCreature: engine.Target{Kind: engine.TargetEachOtherFriendlyCreature},
}

type targets struct {
	This,
	Triggering,
	Creature,
	EachCreature,
	EachFriendlyCreature,
	EachEnemyCreature,
	EachArtifact,
	EachOtherFriendlyCreature engine.Target
}
