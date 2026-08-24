package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Ganymede Archivist
//
//	Logos / Creature / Common / 3 Power
//	Reap: Gain 1 Æmber.
var GanymedeArchivist = card.New(
	"Ganymede Archivist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.WithPower(3),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{Amount: 1}),
)
