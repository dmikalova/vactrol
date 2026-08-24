package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// Bumpsy is a vanilla Brobnar creature — a plain 6-power body.
func TestBumpsy(t *testing.T) {
	g := cardtest.Started(t, game.Brobnar)
	g.AddToHand(Bumpsy, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 6 {
		t.Errorf("Bumpsy power = %d, want 6", g.Power(id))
	}
	if len(Bumpsy.Abilities) != 0 {
		t.Errorf("Bumpsy should have no abilities, got %v", Bumpsy.Abilities)
	}
}
