package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// Batdrone is a vanilla Logos creature — it just enters play as a 2-power body.
func TestBatdrone(t *testing.T) {
	g := cardtest.Started(t, game.Logos)
	g.AddToHand(Batdrone, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 2 {
		t.Errorf("Batdrone power = %d, want 2", g.Power(id))
	}
	if len(Batdrone.Abilities) != 0 {
		t.Errorf("Batdrone should have no abilities, got %v", Batdrone.Abilities)
	}
}
