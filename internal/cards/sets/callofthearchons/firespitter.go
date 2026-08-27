package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Firespitter
//
//	Brobnar / Creature / Common / 5 Power / 1 Armor / Giant
//	Before Fight: Deal 1 Damage to each enemy creature.
var Firespitter = card.New(
	"Firespitter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 32),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits("Giant"),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			Amount: 1,
			Target: card.Target.EachEnemyCreature,
		}),
)
