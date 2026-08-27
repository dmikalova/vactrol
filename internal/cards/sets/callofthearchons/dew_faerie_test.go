package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Dew Faerie reaps for the usual 1 Æmber plus 1 more from its ability.
func TestDewFaerie(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	id := g.AddToBattleline(DewFaerie, 0)
	if g.Power(id) != 2 {
		t.Errorf("Dew Faerie power = %d, want 2", g.Power(id))
	}
	if err := g.Reap(0, id); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("Æmber after reap = %d, want 2 (1 reap + 1 ability)", g.Aember(0))
	}
}
