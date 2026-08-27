package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ganger Chieftain
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
var GangerChieftain = card.New(
	"Ganger Chieftain",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 33),
	card.WithTraits("Giant"),
	card.WithPower(5),
)
