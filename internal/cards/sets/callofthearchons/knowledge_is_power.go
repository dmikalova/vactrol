package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Knowledge is Power
//
//	House:  Logos
//	Type:   Action
//	Rarity: Rare
//
//	Play: Choose one:
//	- Archive a card from your hand
//	- Gain 1 Æmber for each card in your archives.
var KnowledgeIsPower = card.New(
	"Knowledge is Power",
	card.House.Logos,
	card.Type.Action,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 113),
	card.WithAbility(card.Trigger.Play, card.ChooseOne{Options: []card.Effect{
		card.ArchiveFromHand{Count: 1},
		card.GainAember{Amount: 1, Per: card.CardsInArchives{}},
	}}),
)
