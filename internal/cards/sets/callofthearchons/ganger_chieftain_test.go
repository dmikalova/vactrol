package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ganger Chieftain's Play readies and fights with a neighboring creature.
func TestGangerChieftain(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	// A friendly creature already in play becomes Ganger's neighbor when Ganger
	// is played onto the flank beside it.
	neighbor := g.AddToBattleline(cardtest.Vanilla("Neighbor", engine.Brobnar, 4), 0)
	g.AddToBattleline(cardtest.Vanilla("Foe", engine.Mars, 1), 1)
	g.AddToHand(GangerChieftain, 0)

	if _, err := g.PlayCreature(0, 0, false); err != nil { // right flank, beside Neighbor
		t.Fatalf("PlayCreature: %v", err)
	}
	// The neighbor is readied and made to fight the enemy, destroying the 1-power foe.
	if len(g.Battleline(1)) != 0 {
		t.Error("the neighbor should have fought and destroyed the enemy")
	}
	if !g.State.Cards[neighbor].Exhausted {
		t.Error("the neighbor should be exhausted after fighting")
	}
}
