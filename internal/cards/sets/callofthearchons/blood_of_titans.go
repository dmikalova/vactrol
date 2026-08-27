package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Blood of Titans
//
//	Brobnar / Upgrade / Uncommon / 1 Æmber
//	This creature gains +5 power.
var BloodOfTitans = card.New(
	"Blood of Titans",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 50),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{PowerBonus: 5}),
)
