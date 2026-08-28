package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Labwork archives a card from hand.
func TestLabwork(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	g.AddToHand(Labwork, 0)
	g.AddToHand(cardtest.Vanilla("Spare", engine.Logos, 2), 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.State.Archives[0].Count != 1 {
		t.Errorf("archives count = %d, want 1", g.State.Archives[0].Count)
	}
}
