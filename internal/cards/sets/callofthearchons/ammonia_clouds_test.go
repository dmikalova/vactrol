package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ammonia Clouds deals 3 damage to every creature in play — both players' — so
// bodies of 3 power or less are destroyed while tougher ones survive marked.
func TestAmmoniaClouds(t *testing.T) {
	g := cardtest.Started(t, engine.Mars)
	toughAlly := g.AddToBattleline(cardtest.Vanilla("Tough Ally", engine.Mars, 5), 0)
	g.AddToBattleline(cardtest.Vanilla("Weak Ally", engine.Mars, 3), 0)
	toughFoe := g.AddToBattleline(cardtest.Vanilla("Tough Foe", engine.Brobnar, 4), 1)
	g.AddToBattleline(cardtest.Vanilla("Weak Foe", engine.Brobnar, 2), 1)

	g.AddToHand(AmmoniaClouds, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}

	// Weak creatures (power <= 3) die; only the tough bodies remain.
	if got := g.Battleline(0); len(got) != 1 || got[0] != toughAlly {
		t.Errorf("friendly battleline = %v, want just the 5-power ally", got)
	}
	if got := g.Battleline(1); len(got) != 1 || got[0] != toughFoe {
		t.Errorf("enemy battleline = %v, want just the 4-power foe", got)
	}
	if g.Damage(toughAlly) != 3 {
		t.Errorf("tough ally damage = %d, want 3", g.Damage(toughAlly))
	}
	if g.Damage(toughFoe) != 3 {
		t.Errorf("tough foe damage = %d, want 3", g.Damage(toughFoe))
	}
}
