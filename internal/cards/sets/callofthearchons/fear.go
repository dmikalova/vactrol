package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Fear
//
//	House:  Dis
//	Type:   Action
//	Rarity: Common
//
//	Play: Put an enemy creature into its owner's hand.
var Fear = card.New(
	"Fear",
	card.House.Dis,
	card.Type.Action,
	card.Rarity.Common,
	card.Provenance(card.CotA, 58),
	card.WithAbility(
		card.Trigger.Play, card.MoveFromPlay{
			Target:      card.Target.EnemyCreature,
			Destination: card.To.Hand,
		}),
)
