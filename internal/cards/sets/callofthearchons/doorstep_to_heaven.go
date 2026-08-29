package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Doorstep to Heaven
//
//	House:  Sanctum
//	Type:   Action
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Each player with 6 Æmber or more loses all but 5 Æmber.
var DoorstepToHeaven = card.New(
	"Doorstep to Heaven",
	card.House.Sanctum,
	card.Type.Action,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 217),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.EachPlayerLosesAllBut{Keep: 5}),
)
