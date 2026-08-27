package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Francus captures 1 Æmber after an enemy creature is destroyed fighting it. Its
// armor lets it survive the retaliating fight damage.
func TestFrancus(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	francus := g.AddToBattleline(Francus, 0)
	prey := g.AddToBattleline(cardtest.Vanilla("Prey", engine.Mars, 1), 1)
	g.State.Aember[1] = 3 // opponent has Æmber to capture

	if err := g.Fight(0, francus, prey); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("the 1-power prey should be destroyed")
	}
	if g.AmberOn(francus) != 1 {
		t.Errorf("captured aember = %d, want 1", g.AmberOn(francus))
	}
	if g.Aember(1) != 2 {
		t.Errorf("opponent aember = %d, want 2 (3 - 1 captured)", g.Aember(1))
	}
}
