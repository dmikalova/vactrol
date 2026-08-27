package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Cannon's Action deals 2 damage to a creature the controller chooses.
func TestCannon(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	art := g.AddArtifact(Cannon, 0)
	target := g.AddToBattleline(cardtest.Vanilla("Target", engine.Untamed, 5), 1)

	if err := g.UseAction(0, art); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Damage(target) != 2 {
		t.Errorf("damage = %d, want 2", g.Damage(target))
	}
}
