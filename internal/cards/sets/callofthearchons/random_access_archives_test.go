package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Random Access Archives archives the top card of the controller's deck.
func TestRandomAccessArchives(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	top := g.AddToDeck(cardtest.Vanilla("Top", engine.Logos, 2), 0)
	g.AddToHand(RandomAccessArchives, 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != top {
		t.Errorf("archives = %v, want the top deck card [%d]", g.State.Archives[0].IDs[:g.State.Archives[0].Count], top)
	}
}
