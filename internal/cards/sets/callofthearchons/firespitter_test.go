package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Firespitter deals 1 damage to each enemy creature just before it fights.
func TestFirespitter(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	spitter := g.AddToBattleline(Firespitter, 0)
	g.AddToBattleline(cardtest.Vanilla("weak", engine.Brobnar, 1), 1)
	tough := g.AddToBattleline(cardtest.Vanilla("tough", engine.Brobnar, 10), 1)

	if err := g.Fight(0, spitter, tough); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	// The 1-power enemy is destroyed by the before-fight damage, leaving only
	// the defender on the enemy battleline.
	if got := len(g.Battleline(1)); got != 1 {
		t.Errorf("enemy battleline size = %d, want 1 (weak creature destroyed)", got)
	}
	if g.Damage(tough) != 6 { // 1 before-fight + 5 combat
		t.Errorf("defender damage = %d, want 6", g.Damage(tough))
	}
}
