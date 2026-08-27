package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Charette
//
//	Dis / Creature / Common / 4 Power / Demon
//	Play: Charette captures 3 Æmber.
var Charette = card.New(
	"Charette",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 81),
	card.WithPower(4),
	card.WithTraits("Demon"),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{Amount: 3}),
)
