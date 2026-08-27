package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Champion Anaphiel
//
//	Sanctum / Creature / Common / 6 Power / 1 Armor / Knight / Spirit / Taunt
//	Taunt.
var ChampionAnaphiel = card.New(
	"Champion Anaphiel",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 239),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits("Knight", "Spirit"),
	card.WithKeywords(card.Keyword.Taunt),
)
