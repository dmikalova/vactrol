package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ganger Chieftain is a vanilla Brobnar creature — a plain 5-power body.
func TestGangerChieftain(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	g.AddToHand(GangerChieftain, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 5 {
		t.Errorf("Ganger Chieftain power = %d, want 5", g.Power(id))
	}
	if len(GangerChieftain.Abilities) != 0 {
		t.Errorf("Ganger Chieftain should have no abilities, got %v", GangerChieftain.Abilities)
	}
}
