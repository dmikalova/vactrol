package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Begone!
//
//	House:  Sanctum
//	Type:   Action
//	Rarity: Rare
//
//	Play: Choose one:
//	- Destroy each Dis creature
//	- Gain 1 Æmber.
var Begone = card.New(
	"Begone!",
	card.House.Sanctum,
	card.Type.Action,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 212),
	card.WithAbility(card.Trigger.Play, card.ChooseOne{Options: []card.Effect{
		card.Destroy{Target: card.Target.EachCreature.OfHouse(card.House.Dis)},
		card.GainAember{Amount: 1},
	}}),
)
