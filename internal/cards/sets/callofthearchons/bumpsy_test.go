package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Bumpsy makes the opponent lose 1 Æmber when played.
func TestBumpsy(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	g.State.Aember[1] = 3
	g.AddToHand(Bumpsy, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 5 {
		t.Errorf("Bumpsy power = %d, want 5", g.Power(id))
	}
	if g.Aember(1) != 2 {
		t.Errorf("opponent Æmber = %d, want 2", g.Aember(1))
	}
}
