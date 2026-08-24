package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Dust Imp
//
//	Untamed / Creature / Common / 1 Power
//	Reap: Gain 1 Æmber.
var DustImp = card.New(
	"Dust Imp",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.WithPower(1),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{Amount: 1}),
)
