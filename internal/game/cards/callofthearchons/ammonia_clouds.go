package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Ammonia Clouds
//
//	Mars / Action / Common
//	Play: Deal 3 damage to each creature.
var AmmoniaClouds = card.New(
	"Ammonia Clouds",
	card.House.Mars,
	card.Type.Action,
	card.Rarity.Common,
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Target: card.Target.EachCreature,
		}),
)
