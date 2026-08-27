package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Bilgum Avalanche deals 2 damage to each enemy creature every time its
// controller forges a key.
func TestBilgumAvalanche(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	g.AddToBattleline(BilgumAvalanche, 0)
	foe := g.AddToBattleline(cardtest.Vanilla("foe", engine.Brobnar, 5), 1)

	g.State.Aember[0] = engine.KeyCost
	g.BeginTurn(0)

	if g.Damage(foe) != 2 {
		t.Errorf("enemy creature damage = %d, want 2", g.Damage(foe))
	}
}
