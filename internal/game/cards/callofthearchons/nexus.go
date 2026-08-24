package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Nexus
//
//	Logos / Creature / Common / 4 Power
var Nexus = card.New(
	"Nexus",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.WithPower(4),
)
