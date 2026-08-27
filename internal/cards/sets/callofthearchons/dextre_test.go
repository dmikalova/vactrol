package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Dextre captures 1 Æmber from the opponent when played.
func TestDextrePlay(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	g.State.Aember[1] = 3
	g.AddToHand(Dextre, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.AmberOn(id) != 1 {
		t.Errorf("captured Æmber = %d, want 1", g.AmberOn(id))
	}
	if g.Aember(1) != 2 {
		t.Errorf("opponent Æmber = %d, want 2", g.Aember(1))
	}
}

// When destroyed, Dextre returns to the top of its owner's deck instead of the
// discard pile.
func TestDextreDestroyed(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	id := g.AddToBattleline(Dextre, 0)

	g.DestroyEach(0, []engine.LocalID{id})

	if g.State.Deck[0].Count != 1 {
		t.Fatalf("deck count = %d, want 1", g.State.Deck[0].Count)
	}
	if g.State.Deck[0].IDs[0] != id {
		t.Errorf("top of deck = %v, want Dextre (%v)", g.State.Deck[0].IDs[0], id)
	}
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard count = %d, want 0 (Dextre should not be discarded)", g.State.Discard[0].Count)
	}
}
