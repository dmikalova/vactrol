package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Radiant Truth
//
//	House:  Sanctum
//	Type:   Action
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Stun each enemy creature that is not on a flank.
var RadiantTruth = card.New(
	"Radiant Truth",
	card.House.Sanctum,
	card.Type.Action,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 224),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Stun{Target: card.Target.EachEnemyCreature.NotOnFlank()}),
)
