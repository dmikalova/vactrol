package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Nepenthe Seed sacrifices itself and returns a card from the discard pile to hand.
func TestNepentheSeed(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	seed := g.AddArtifact(NepentheSeed, 0)
	ghost := g.AddToBattleline(cardtest.Vanilla("Ghost", engine.Untamed, 4), 0)
	g.DestroyEach(0, []engine.LocalID{ghost}) // a card in the discard to retrieve
	before := len(g.Hand(0))

	if err := g.UseAction(0, seed); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if len(g.Artifacts(0)) != 0 {
		t.Error("Nepenthe Seed should sacrifice itself")
	}
	if len(g.Hand(0)) != before+1 {
		t.Errorf("hand = %d, want %d (a card returned from discard)", len(g.Hand(0)), before+1)
	}
}
