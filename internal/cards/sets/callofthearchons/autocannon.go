package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Autocannon
//
//	Brobnar / Artifact / Rare / 1 Æmber / Weapon
//	After a creature enters play, deal 1 damage to it.
var Autocannon = card.New(
	"Autocannon",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 19),
	card.WithAemberBonus(1),
	card.WithTraits("Weapon"),
	card.WithAbility(
		card.Trigger.AfterCreatureEnters, card.DealDamage{
			Amount: 1,
			Target: card.Target.Triggering,
		}),
)
