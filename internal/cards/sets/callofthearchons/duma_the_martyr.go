package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Duma the Martyr
//
//	Sanctum / Creature / Rare / 3 Power / Human
//	Destroyed: Fully heal each other friendly creature and draw 2 cards.
var DumaTheMartyr = card.New(
	"Duma the Martyr",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 242),
	card.WithPower(3),
	card.WithTraits("Human"),
	card.WithAbility(
		card.Trigger.Destroyed, card.Sequence{
			Effects: []card.Effect{
				card.Heal{Fully: true, Target: card.Target.EachOtherFriendlyCreature},
				card.Draw{Amount: 2},
			},
		}),
)
