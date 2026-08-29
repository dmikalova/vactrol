package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dance of Doom
//
//	House:  Dis
//	Type:   Action
//	Rarity: Rare
//
//	Play: Choose a creature - destroy each creature with the same power as the chosen creature.
var DanceOfDoom = card.New(
	"Dance of Doom",
	card.House.Dis,
	card.Type.Action,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 57),
	card.WithAbility(
		card.Trigger.Play, card.DestroySamePower{}),
)
