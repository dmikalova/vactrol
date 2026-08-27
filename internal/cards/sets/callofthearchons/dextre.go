package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dextre
//
//	Logos / Creature / Common / 3 Power / Human / Scientist
//	Play: Dextre captures 1 Æmber.
//	Destroyed: Put Dextre on top of its owner's deck.
var Dextre = card.New(
	"Dextre",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 138),
	card.WithPower(3),
	card.WithTraits("Human", "Scientist"),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{Amount: 1}),
	card.WithAbility(
		card.Trigger.Destroyed, card.ReturnToDeck{Target: card.Target.This}),
)
