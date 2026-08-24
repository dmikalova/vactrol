package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// Snufflegator is a vanilla Untamed creature — a plain 5-power body.
func TestSnufflegator(t *testing.T) {
	g := cardtest.Started(t, game.Untamed)
	g.AddToHand(Snufflegator, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 5 {
		t.Errorf("Snufflegator power = %d, want 5", g.Power(id))
	}
	if len(Snufflegator.Abilities) != 0 {
		t.Errorf("Snufflegator should have no abilities, got %v", Snufflegator.Abilities)
	}
}
