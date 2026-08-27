package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Anger
//
//	Brobnar / Action / Common / 1 Æmber
//	Play: Ready and fight with a friendly creature.
var Anger = card.New(
	"Anger",
	card.House.Brobnar,
	card.Type.Action,
	card.Rarity.Common,
	card.Provenance(card.CotA, 1),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.OnChosenCreature{
			Player: card.Controller,
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
