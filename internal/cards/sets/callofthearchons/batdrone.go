package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Batdrone
//
//	Logos / Creature / Common / 2 Power / Robot / Skirmish
//	Fight: Steal 1 Æmber.
var Batdrone = card.New(
	"Batdrone",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 136),
	card.WithPower(2),
	card.WithTraits("Robot"),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.StealAember{Amount: 1}),
)
