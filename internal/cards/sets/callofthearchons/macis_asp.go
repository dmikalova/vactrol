package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Macis Asp
//
//	Shadows / Creature / Uncommon / 3 Power / Beast / Skirmish / Poison
//	Skirmish.
//	Poison.
var MacisAsp = card.New(
	"Macis Asp",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 301),
	card.WithPower(3),
	card.WithTraits("Beast"),
	card.WithKeywords(card.Keyword.Skirmish, card.Keyword.Poison),
)
