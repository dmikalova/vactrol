package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dust Imp
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Imp
//
//	Reap: Gain 1 Æmber.
var DustImp = card.New(
	"Dust Imp",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 83),
	card.WithPower(1),
	card.WithTraits("Imp"),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{Amount: 1}),
)
