package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Champion Tabris
//
//	Sanctum / Creature / Uncommon / 6 Power / 2 Armor / Human / Knight
//	Fight: Champion Tabris captures 1 Æmber.
var ChampionTabris = card.New(
	"Champion Tabris",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 240),
	card.WithPower(6),
	card.WithArmor(2),
	card.WithTraits("Human", "Knight"),
	card.WithAbility(
		card.Trigger.Fight, card.CaptureAember{Amount: 1}),
)
