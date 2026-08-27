package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Macis Asp is a 3-power Beast with both Skirmish (no retaliation) and Poison
// (any damage it deals destroys the creature).
func TestMacisAsp(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	g.AddToHand(MacisAsp, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 3 {
		t.Errorf("Macis Asp power = %d, want 3", g.Power(id))
	}

	var hasSkirmish, hasPoison bool
	for _, kw := range MacisAsp.Keywords {
		switch kw {
		case engine.Skirmish:
			hasSkirmish = true
		case engine.Poison:
			hasPoison = true
		}
	}
	if !hasSkirmish || !hasPoison {
		t.Errorf("Macis Asp keywords = %v, want Skirmish and Poison", MacisAsp.Keywords)
	}
}
