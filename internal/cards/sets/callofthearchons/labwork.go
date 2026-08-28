package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Labwork
//
//	House:  Logos
//	Type:   Action
//	Rarity: Common
//	Æmber:  1
//
//	Play: Archive a card from your hand.
var Labwork = card.New(
	"Labwork",
	card.House.Logos,
	card.Type.Action,
	card.Rarity.Common,
	card.Provenance(card.CotA, 114),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.ArchiveFromHand{Count: 1}),
)
