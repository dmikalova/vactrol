package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

func contains(ids []engine.LocalID, id engine.LocalID) bool {
	for _, other := range ids {
		if other == id {
			return true
		}
	}
	return false
}

// Earthshaker destroys each creature with power 3 or lower on Play, sparing
// stronger creatures and itself.
func TestEarthshaker(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	weak := g.AddToBattleline(cardtest.Vanilla("weak", engine.Brobnar, 2), 0)
	strong := g.AddToBattleline(cardtest.Vanilla("strong", engine.Brobnar, 5), 1)
	g.AddToHand(Earthshaker, 0)

	shaker, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}

	if contains(g.Battleline(0), weak) {
		t.Error("weak creature should be destroyed")
	}
	if !contains(g.Battleline(1), strong) {
		t.Error("strong creature should survive")
	}
	if !contains(g.Battleline(0), shaker) {
		t.Error("Earthshaker (7 power) should survive its own Play")
	}
}
