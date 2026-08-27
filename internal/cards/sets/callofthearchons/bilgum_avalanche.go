package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bilgum Avalanche
//
//	Brobnar / Creature / Rare / 5 Power / Giant
//	After you forge a key, deal 2 Damage to each enemy creature.
var BilgumAvalanche = card.New(
	"Bilgum Avalanche",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 28),
	card.WithPower(5),
	card.WithTraits("Giant"),
	card.WithAbility(
		card.Trigger.AfterForgeKey, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachEnemyCreature,
		}),
)
