package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Library of the Damned archives a card from hand via its Action.
func TestLibraryOfTheDamned(t *testing.T) {
	g := cardtest.Started(t, engine.Dis)
	g.AddToHand(cardtest.Vanilla("Spare", engine.Dis, 2), 0)
	lib := g.AddArtifact(LibraryOfTheDamned, 0)

	if err := g.UseAction(0, lib); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.State.Archives[0].Count != 1 {
		t.Errorf("archives count = %d, want 1", g.State.Archives[0].Count)
	}
}
