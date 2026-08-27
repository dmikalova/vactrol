package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ghostly Hand always yields its 2 Æmber bonus and, when the opponent holds
// exactly 1 Æmber, steals it for a third.
func TestGhostlyHand(t *testing.T) {
	// Condition met: opponent at exactly 1 — bonus 2 plus 1 stolen makes 3.
	g := cardtest.Started(t, engine.Shadows)
	g.State.Aember[1] = 1
	g.AddToHand(GhostlyHand, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.Aember(0) != 3 {
		t.Errorf("controller Æmber = %d, want 3", g.Aember(0))
	}
	if g.Aember(1) != 0 {
		t.Errorf("opponent Æmber = %d, want 0", g.Aember(1))
	}

	// Condition not met: opponent at 2 — only the 2 Æmber bonus applies.
	g = cardtest.Started(t, engine.Shadows)
	g.State.Aember[1] = 2
	g.AddToHand(GhostlyHand, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("controller Æmber = %d, want 2", g.Aember(0))
	}
	if g.Aember(1) != 2 {
		t.Errorf("opponent Æmber = %d, want 2", g.Aember(1))
	}
}
