package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// DustImp gains 2 Æmber for its controller when it is destroyed.
func TestDustImp(t *testing.T) {
	g := cardtest.Started(t, engine.Dis)
	id := g.AddToBattleline(DustImp, 0)

	g.DestroyEach(0, []engine.LocalID{id})
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (Destroyed: gain 2 Æmber)", g.Aember(0))
	}
}
