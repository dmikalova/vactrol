package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dew Faerie
//
//	Untamed / Creature / Common / 2 Power / Faerie / Elusive
//	Reap: Gain 1 Æmber.
var DewFaerie = card.New(
	"Dew Faerie",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 350),
	card.WithPower(2),
	card.WithTraits("Faerie"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{Amount: 1}),
)
