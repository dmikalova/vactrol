package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Champion Tabris
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  2
//	Traits: Human • Knight
//
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
