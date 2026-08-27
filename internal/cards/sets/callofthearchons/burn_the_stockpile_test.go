package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Burn the Stockpile drains 4 Æmber from an opponent flush with 7 or more, but
// does nothing to a smaller pool.
func TestBurnTheStockpile(t *testing.T) {
	// Condition met: opponent at 7 loses 4, ending at 3.
	g := cardtest.Started(t, engine.Brobnar)
	g.State.Aember[1] = 7
	g.AddToHand(BurnTheStockpile, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.Aember(1) != 3 {
		t.Errorf("opponent Æmber = %d, want 3", g.Aember(1))
	}

	// Condition not met: opponent at 6 keeps every Æmber.
	g = cardtest.Started(t, engine.Brobnar)
	g.State.Aember[1] = 6
	g.AddToHand(BurnTheStockpile, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.Aember(1) != 6 {
		t.Errorf("opponent Æmber = %d, want 6", g.Aember(1))
	}
}
