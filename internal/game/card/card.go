// Package card is a small, ergonomic facade over the game engine for *authoring*
// cards. Card-set packages import this package (as `card`) so definitions read
// like card.New(...), card.House.Brobnar, card.DealDamage{...} without referring to the
// engine's package name directly.
//
// Almost everything here is a re-export — type aliases (which are fully
// interchangeable with the engine's types) and constant/function values — of
// identifiers in package game, shaping the authoring API. The one bit of logic
// is the card registry (see register.go): New builds a card and enrolls it in
// the database, so a set package declares each card as `var X = card.New(...)`
// — one line, no init, no central list.
package card

import "github.com/dmikalova/vactrol/internal/game"

// Types you name directly when authoring cards (structs and interfaces). The
// enum-like categories (house, card type, rarity, keyword, trigger, controller,
// and target) are exposed as grouped namespaces further down, not as loose constants.
type (
	// Definition is a card blueprint.
	Definition = game.CardDefinition
	// Option configures a Definition (see the With* helpers).
	Option = game.CardOption
	// Ability pairs a trigger with an effect.
	Ability = game.Ability
	// StaticModifier is a continuous stat change applied by an upgrade.
	StaticModifier = game.StaticModifier
	// Trait is a freeform label such as "Giant" or "Weapon".
	Trait = game.Trait

	// Effect is a card effect (see DealDamage, GainAember, ...).
	Effect = game.Effect
	// CreatureVerb is an action applied to a chosen creature (ReadyVerb, FightVerb).
	CreatureVerb = game.CreatureVerb

	DealDamage       = game.DealDamage
	Exalt            = game.Exalt
	GainAember       = game.GainAember
	OnChosenCreature = game.OnChosenCreature
	Sequence         = game.Sequence
	ReadyVerb        = game.ReadyVerb
	FightVerb        = game.FightVerb
	ReturnToDeck     = game.ReturnToDeck
)

// Option helpers (re-exported from the engine). The card constructor New lives
// in register.go because building a card there also enrolls it in the database.
var (
	WithPower       = game.WithPower
	WithArmor       = game.WithArmor
	WithTraits      = game.WithTraits
	WithKeywords    = game.WithKeywords
	WithAemberBonus = game.WithAemberBonus
	WithStatic      = game.WithStatic
	WithAbility     = game.WithAbility
)

// The enum-like categories are exposed as grouped namespaces so related values
// stay together and read unambiguously — card.House.Brobnar is plainly a house,
// card.Type.Action plainly a card type — all from this single `card` import.
// Treat these package-level vars as read-only.

// House groups the faction values, e.g. card.House.Brobnar.
var House = houses{
	Brobnar: game.Brobnar,
	Dis:     game.Dis,
	Logos:   game.Logos,
	Mars:    game.Mars,
	Sanctum: game.Sanctum,
	Shadows: game.Shadows,
	Untamed: game.Untamed,
}

type houses struct {
	Brobnar, Dis, Logos, Mars, Sanctum, Shadows, Untamed game.House
}

// Type groups the card-type values, e.g. card.Type.Creature.
var Type = cardTypes{
	Creature: game.Creature,
	Action:   game.Action,
	Artifact: game.Artifact,
	Upgrade:  game.Upgrade,
}

type cardTypes struct {
	Creature, Action, Artifact, Upgrade game.CardType
}

// Rarity groups the rarity values, e.g. card.Rarity.Common.
var Rarity = rarities{
	Common:   game.Common,
	Uncommon: game.Uncommon,
	Rare:     game.Rare,
	Special:  game.Special,
}

type rarities struct {
	Common, Uncommon, Rare, Special game.Rarity
}

// Keyword groups the keyword values, e.g. card.Keyword.Skirmish.
var Keyword = keywords{
	Skirmish: game.Skirmish,
	Poison:   game.Poison,
	Elusive:  game.Elusive,
	Taunt:    game.Taunt,
}

type keywords struct {
	Skirmish, Poison, Elusive, Taunt game.Keyword
}

// Trigger groups the ability triggers, e.g. card.Trigger.Play or
// card.Trigger.AfterForgeKey.
var Trigger = triggers{
	Play:                game.TriggerPlay,
	Reap:                game.TriggerReap,
	Fight:               game.TriggerFight,
	Action:              game.TriggerAction,
	AfterForgeKey:       game.TriggerAfterForgeKey,
	AfterCreatureEnters: game.TriggerAfterCreatureEnters,
	Destroyed:           game.TriggerDestroyed,
}

type triggers struct {
	Play, Reap, Fight, Action, AfterForgeKey, AfterCreatureEnters, Destroyed game.Trigger
}

// Controller groups the chosen-creature sides: card.Controller.Friendly / .Enemy.
var Controller = controllers{
	Friendly: game.Friendly,
	Enemy:    game.Enemy,
}

type controllers struct {
	Friendly, Enemy game.Controller
}

// Target groups ready-made targets, e.g. card.Target.EachEnemyCreature.
var Target = targets{
	This:              game.Target{Kind: game.TargetThisCreature},
	Triggering:        game.Target{Kind: game.TargetTriggeringCreature},
	EachCreature:      game.Target{Kind: game.TargetEachCreature},
	EachEnemyCreature: game.Target{Kind: game.TargetEachEnemyCreature},
	EachArtifact:      game.Target{Kind: game.TargetEachArtifact},
}

type targets struct {
	This, Triggering, EachCreature, EachEnemyCreature, EachArtifact game.Target
}
