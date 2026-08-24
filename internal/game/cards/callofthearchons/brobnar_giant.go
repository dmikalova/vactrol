package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Brobnar Giant
//
//	Brobnar / Creature / Rare / 5 Power / 0 Armor / Giant
//	After you forge a key, deal 2 damage to each enemy creature.
var BrobnarGiant = card.New(
	"Brobnar Giant",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.WithPower(5),
	card.WithArmor(0),
	card.WithTraits("Giant"),
	card.WithAbility(
		card.Trigger.AfterForgeKey, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachEnemyCreature,
		}),
)
