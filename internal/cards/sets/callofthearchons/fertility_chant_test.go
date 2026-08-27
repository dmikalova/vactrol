package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Fertility Chant gives its controller 4 Æmber (pips) and the opponent 2.
func TestFertilityChant(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	g.AddToHand(FertilityChant, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.Aember(0) != 4 {
		t.Errorf("controller Æmber = %d, want 4", g.Aember(0))
	}
	if g.Aember(1) != 2 {
		t.Errorf("opponent Æmber = %d, want 2", g.Aember(1))
	}
}
