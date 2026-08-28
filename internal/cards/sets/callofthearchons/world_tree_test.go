package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// World Tree returns a creature from the discard pile to the top of the deck.
func TestWorldTree(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	tree := g.AddArtifact(WorldTree, 0)
	ghost := g.AddToBattleline(cardtest.Vanilla("Ghost", engine.Untamed, 4), 0)
	g.DestroyEach(0, []engine.LocalID{ghost}) // Ghost -> discard

	if err := g.UseAction(0, tree); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.State.Deck[0].Count == 0 || g.State.Deck[0].IDs[0] != ghost {
		t.Errorf("deck top = %v, want Ghost (%d)", g.State.Deck[0].IDs[:g.State.Deck[0].Count], ghost)
	}
}
