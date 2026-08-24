package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Blood Money
//
//	Brobnar / Action / Uncommon
//	Play: Exalt an enemy creature 2 times.
var BloodMoney = card.New(
	"Blood Money",
	card.House.Brobnar,
	card.Type.Action,
	card.Rarity.Uncommon,
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Controller: card.Controller.Enemy,
			Times:      2,
		}),
)
