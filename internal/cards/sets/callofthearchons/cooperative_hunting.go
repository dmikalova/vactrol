package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Cooperative Hunting
//
//	House:  Untamed
//	Type:   Action
//	Rarity: Common
//
//	Play: For each friendly creature in play, deal 1 damage to a creature.
var CooperativeHunting = card.New(
	"Cooperative Hunting",
	card.House.Untamed,
	card.Type.Action,
	card.Rarity.Common,
	card.Provenance(card.CotA, 319),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			Per:    card.FriendlyCreaturesInPlay{},
			Target: card.Target.Creature,
		}),
)
