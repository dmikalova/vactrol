package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ancient Bear
//
//	Untamed / Creature / Common / 5 Power / Skirmish
var AncientBear = card.New(
	"Ancient Bear",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 345),
	card.WithPower(5),
	card.WithKeywords(card.Keyword.Skirmish),
)
