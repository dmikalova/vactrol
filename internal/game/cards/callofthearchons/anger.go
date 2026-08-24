package callofthearchons

import "github.com/dmikalova/vactrol/internal/game/card"

// Anger
//
//	Brobnar / Action / Common / 1 Æmber
//	Play: Ready and fight with a friendly creature.
var Anger = card.New(
	"Anger",
	card.House.Brobnar,
	card.Type.Action,
	card.Rarity.Common,
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.OnChosenCreature{
			Controller: card.Controller.Friendly,
			Verbs:      []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
