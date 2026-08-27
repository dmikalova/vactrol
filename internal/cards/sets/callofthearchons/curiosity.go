package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Curiosity
//
//	Untamed / Action / Rare / 1 Æmber
//	Play: Destroy each Scientist trait creature.
var Curiosity = card.New(
	"Curiosity",
	card.House.Untamed,
	card.Type.Action,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 320),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{Target: card.Target.EachCreature.WithTrait("Scientist")}),
)
