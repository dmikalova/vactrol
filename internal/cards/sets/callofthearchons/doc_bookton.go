package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Doc Bookton
//
//	Logos / Creature / Common / 5 Power / Human / Scientist
//	Reap: Draw a card.
var DocBookton = card.New(
	"Doc Bookton",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 139),
	card.WithPower(5),
	card.WithTraits("Human", "Scientist"),
	card.WithAbility(
		card.Trigger.Reap, card.Draw{Amount: 1}),
)
