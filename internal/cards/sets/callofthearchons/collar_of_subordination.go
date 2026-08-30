package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Collar of Subordination
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Rare
//
//	Play: Take control of this creature until Collar of Subordination leaves play.
var CollarOfSubordination = card.New(
	"Collar of Subordination",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 105),
	card.WithStatic(card.StaticModifier{TakesControl: true}),
	card.WithAbility(
		card.Trigger.Play, card.TakeControl{}),
)
