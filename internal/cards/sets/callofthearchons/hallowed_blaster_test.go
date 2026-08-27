package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Hallowed Blaster's Action heals 3 damage from a creature the controller chooses.
func TestHallowedBlaster(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	art := g.AddArtifact(HallowedBlaster, 0)
	c := g.AddToBattleline(cardtest.Vanilla("Wounded", engine.Sanctum, 6), 0)
	g.State.Cards[c].Damage = 4

	if err := g.UseAction(0, art); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Damage(c) != 1 {
		t.Errorf("damage = %d, want 1 (4 - 3)", g.Damage(c))
	}
}
