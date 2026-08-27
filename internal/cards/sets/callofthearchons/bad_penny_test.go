package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// When destroyed, Bad Penny returns to its owner's hand instead of the discard.
func TestBadPenny(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	id := g.AddToBattleline(BadPenny, 0)

	g.DestroyEach(0, []engine.LocalID{id})

	if g.State.Hand[0].Count != 1 {
		t.Fatalf("hand count = %d, want 1", g.State.Hand[0].Count)
	}
	if g.State.Hand[0].IDs[0] != id {
		t.Errorf("hand card = %v, want Bad Penny (%v)", g.State.Hand[0].IDs[0], id)
	}
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard count = %d, want 0", g.State.Discard[0].Count)
	}
	if len(g.Battleline(0)) != 0 {
		t.Errorf("battleline = %d creatures, want 0", len(g.Battleline(0)))
	}
}
