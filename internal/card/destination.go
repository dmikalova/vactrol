package card

import "github.com/dmikalova/vactrol/internal/engine"

// To groups the destinations an effect can put a card, e.g. card.To.TopOfDeck. An
// effect that moves a card takes one — see card.ReturnFromDiscard's Destination.
// It mirrors the engine's destination.go.
var To = destinations{
	Hand:         engine.ToHand,
	TopOfDeck:    engine.ToTopOfDeck,
	BottomOfDeck: engine.ToBottomOfDeck,
	DeckShuffled: engine.ToDeckShuffled,
}

type destinations struct {
	Hand,
	TopOfDeck,
	BottomOfDeck,
	DeckShuffled engine.Destination
}

// Destination names where an effect puts a card it moves (see card.To).
type Destination = engine.Destination
