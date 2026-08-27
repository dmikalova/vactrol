package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// DustImp gains an extra Æmber from its Reap ability, on top of the +1 for
// reaping.
func TestDustImp(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	id := g.AddToBattleline(DustImp, 0)

	if err := g.Reap(0, id); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (+1 reap, +1 ability)", g.Aember(0))
	}
	if !g.Exhausted(id) {
		t.Error("Dust Imp should be exhausted after reaping")
	}
}
