package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Barehanded
//
//	Brobnar / Action / Rare / 1 Æmber
//	Play: Put each artifact on top of its owner's deck.
var Barehanded = card.New(
	"Barehanded",
	card.House.Brobnar,
	card.Type.Action,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 2),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ReturnToDeck{
			Target: card.Target.EachArtifact,
		}),
)
