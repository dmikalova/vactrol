package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Burn the Stockpile
//
//	Brobnar / Action / Uncommon
//	Play: If your opponent has 7 Æmber or more, they lose 4 Æmber.
var BurnTheStockpile = card.New(
	"Burn the Stockpile",
	card.House.Brobnar,
	card.Type.Action,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 5),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAemberAtLeast{Amount: 7},
			Then: card.LoseAember{Player: card.Opponent, Amount: 4},
		}),
)
