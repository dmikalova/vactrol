package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dr. Escotera
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	Play: Gain 1 Æmber for each key your opponent has forged.
var DrEscotera = card.New(
	"Dr. Escotera",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 140),
	card.WithPower(4),
	card.WithTraits("Cyborg", "Scientist"),
	card.WithAbility(card.Trigger.Play, card.GainAember{Amount: 1, Per: card.OpponentForgedKeys{}}),
)
