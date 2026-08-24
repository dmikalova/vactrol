package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// Autocannon deals 1 damage to a creature as it enters play.
func TestAutocannon(t *testing.T) {
	g := cardtest.Started(t, game.Brobnar)
	g.AddToHand(Autocannon, 0)
	if _, err := g.PlayArtifact(0, 0); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}

	// Playing a creature triggers Autocannon, zapping the newcomer for 1.
	g.AddToHand(cardtest.Vanilla("Newcomer", game.Brobnar, 3), 0)
	entered, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Damage(entered) != 1 {
		t.Errorf("entering creature damage = %d, want 1", g.Damage(entered))
	}
}
