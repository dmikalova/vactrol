package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Hidden Stash archives a card from hand.
func TestHiddenStash(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	g.AddToHand(HiddenStash, 0)                                  // index 0: the action to play
	g.AddToHand(cardtest.Vanilla("Spare", engine.Shadows, 2), 0) // a card to archive

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.State.Archives[0].Count != 1 {
		t.Errorf("archives count = %d, want 1", g.State.Archives[0].Count)
	}
}
