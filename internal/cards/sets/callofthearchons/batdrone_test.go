package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Batdrone steals 1 Æmber when it fights (and Skirmish spares it return damage).
func TestBatdrone(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	g.State.Aember[1] = 2
	id := g.AddToBattleline(Batdrone, 0)
	enemy := g.AddToBattleline(cardtest.Vanilla(" Box", engine.Logos, 3), 1)
	if g.Power(id) != 2 {
		t.Errorf("Batdrone power = %d, want 2", g.Power(id))
	}
	if err := g.Fight(0, id, enemy); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.Aember(0) != 1 || g.Aember(1) != 1 {
		t.Errorf("after steal: you=%d opp=%d, want 1/1", g.Aember(0), g.Aember(1))
	}
	if g.Damage(id) != 0 {
		t.Errorf("Skirmish should spare Batdrone; damage = %d", g.Damage(id))
	}
}
