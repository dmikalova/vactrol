package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ganymede Archivist archives a card from hand when it reaps.
func TestGanymedeArchivist(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	g.AddToHand(cardtest.Vanilla("Spare", engine.Logos, 2), 0)
	arch := g.AddToBattleline(GanymedeArchivist, 0)

	if err := g.Reap(0, arch); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.State.Archives[0].Count != 1 {
		t.Errorf("archives count = %d, want 1", g.State.Archives[0].Count)
	}
}
