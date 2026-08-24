package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Brute Strength
//
//	Brobnar / Upgrade / Uncommon / 1 Æmber
//	This creature gets +5 power.
var BruteStrength = card.New(
	"Brute Strength",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{PowerBonus: 5}),
)
