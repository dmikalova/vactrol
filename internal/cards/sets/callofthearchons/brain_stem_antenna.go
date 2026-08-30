package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Brain Stem Antenna
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains, "After you play a Mars creature, ready this creature and for the remainder of the turn it belongs to house Mars."
var BrainStemAntenna = card.New(
	"Brain Stem Antenna",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 209),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.AfterCreatureEnters, Effect: card.ReadyAndBelongToHouseAfterYouPlayCreature{House: card.House.Mars}},
		},
	}),
)
