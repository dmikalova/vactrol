package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Creeping Oblivion
//
//	House:  Dis
//	Type:   Action
//	Rarity: Rare
//
//	Play: Purge up to 2 cards from a discard pile.
var CreepingOblivion = card.New(
	"Creeping Oblivion",
	card.House.Dis,
	card.Type.Action,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 56),
	card.WithAbility(
		card.Trigger.Play, card.Purge{
			Zone:  card.Discard,
			Count: 2,
			UpTo:  true,
		}),
)
