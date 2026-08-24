package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// Nexus is a vanilla Logos creature — a plain 4-power body.
func TestNexus(t *testing.T) {
	g := cardtest.Started(t, game.Logos)
	g.AddToHand(Nexus, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 4 {
		t.Errorf("Nexus power = %d, want 4", g.Power(id))
	}
	if len(Nexus.Abilities) != 0 {
		t.Errorf("Nexus should have no abilities, got %v", Nexus.Abilities)
	}
}
