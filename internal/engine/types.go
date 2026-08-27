// Package engine is the core rules engine for Vactrol, a KeyForge-style
// card game. This file defines the fundamental enumerations shared across the
// engine: houses, rarities, card types, traits, keywords, and ability triggers.
package engine

import "strings"

// House is the faction a card belongs to. It is a small integer enum so it can
// live inside the flat, value-copyable GameState without introducing pointers.
type House uint8

const (
	// HouseNone means no house is chosen/assigned.
	HouseNone House = iota
	Brobnar
	Dis
	Logos
	Mars
	Sanctum
	Shadows
	Untamed
)

// houseNames maps a House to its printed name, indexed by the enum value.
var houseNames = [...]string{
	"None", "Brobnar", "Dis", "Logos", "Mars", "Sanctum", "Shadows", "Untamed",
}

// String returns the printed house name.
func (h House) String() string {
	if int(h) < len(houseNames) {
		return houseNames[h]
	}
	return "Unknown"
}

// ParseHouse resolves a house name (case-insensitive) to its House value. The
// boolean result is false when the name matches no house.
func ParseHouse(name string) (House, bool) {
	for h := HouseNone + 1; int(h) < len(houseNames); h++ {
		if strings.EqualFold(houseNames[h], name) {
			return h, true
		}
	}
	return HouseNone, false
}

// Rarity describes how frequently a card appears when generating a deck. It only
// ever lives in the read-only card catalog, so a string is fine here.
type Rarity string

const (
	Common   Rarity = "Common"
	Uncommon Rarity = "Uncommon"
	Rare     Rarity = "Rare"
	Special  Rarity = "Special"
)

// CardType is one of the four Vactrol/KeyForge card types.
type CardType string

const (
	Creature CardType = "Creature"
	Action   CardType = "Action"
	Artifact CardType = "Artifact"
	Upgrade  CardType = "Upgrade"
)

// Trait is a flavor/type label printed on a card (e.g. "Giant", "Weapon").
// Traits carry no inherent rules meaning on their own; other cards reference them.
type Trait string

// Keyword is a rules shorthand a card can have (e.g. Skirmish, Poison).
type Keyword string

const (
	// Skirmish: this creature takes no damage when it is used to fight.
	Skirmish Keyword = "Skirmish"
	// Poison: any damage dealt to this creature destroys it.
	Poison Keyword = "Poison"
	// Elusive: defined for completeness; behavior not yet implemented.
	Elusive Keyword = "Elusive"
	// Taunt: defined for completeness; behavior not yet implemented.
	Taunt Keyword = "Taunt"
)

// Trigger identifies when an ability's effect resolves.
type Trigger int

const (
	// TriggerAfterPlay fires after a card is played from hand.
	TriggerAfterPlay Trigger = iota
	// TriggerAfterReap fires after a creature reaps.
	TriggerAfterReap
	// TriggerAfterFight fires after a creature fights (once damage is dealt).
	TriggerAfterFight
	// TriggerAction fires when a creature uses its "Action:" ability.
	TriggerAction
	// TriggerAfterForgeKey fires after this card's controller forges a key.
	TriggerAfterForgeKey
	// TriggerAfterCreatureEnters fires after any creature enters play.
	TriggerAfterCreatureEnters
	// TriggerDestroyed fires when the card is destroyed.
	TriggerDestroyed
	// TriggerBeforeFight fires when a creature is used to fight, before any
	// combat damage is dealt.
	TriggerBeforeFight
)

// prefix returns the printed text prefix for a trigger and whether the effect
// clause that follows should be capitalized (colon-style triggers start a new
// sentence; comma-style triggers continue one).
func (t Trigger) prefix() (text string, capitalizeEffect bool) {
	switch t {
	case TriggerAfterPlay:
		return "Play: ", true
	case TriggerAfterReap:
		return "Reap: ", true
	case TriggerAfterFight:
		return "Fight: ", true
	case TriggerBeforeFight:
		return "Before Fight: ", true
	case TriggerAction:
		return "Action: ", true
	case TriggerDestroyed:
		return "Destroyed: ", true
	case TriggerAfterForgeKey:
		return "After you forge a key, ", false
	case TriggerAfterCreatureEnters:
		return "After a creature enters play, ", false
	default:
		return "", true
	}
}
