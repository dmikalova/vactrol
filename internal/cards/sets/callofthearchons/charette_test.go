package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Charette captures 3 Æmber from the opponent's pool onto itself when played.
func TestCharette(t *testing.T) {
	g := cardtest.Started(t, engine.Dis)
	g.State.Aember[1] = 5
	g.AddToHand(Charette, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 4 {
		t.Errorf("Charette power = %d, want 4", g.Power(id))
	}
	if g.AmberOn(id) != 3 {
		t.Errorf("captured Æmber = %d, want 3", g.AmberOn(id))
	}
	if g.Aember(1) != 2 {
		t.Errorf("opponent Æmber = %d, want 2", g.Aember(1))
	}
}
