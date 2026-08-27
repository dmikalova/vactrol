package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ghostly Hand
//
//	Shadows / Action / Common / 2 Æmber
//	Play: If your opponent has exactly 1 Æmber, steal 1 Æmber.
var GhostlyHand = card.New(
	"Ghostly Hand",
	card.House.Shadows,
	card.Type.Action,
	card.Rarity.Common,
	card.Provenance(card.CotA, 270),
	card.WithAemberBonus(2),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAemberExactly{Amount: 1},
			Then: card.StealAember{Amount: 1},
		}),
)
