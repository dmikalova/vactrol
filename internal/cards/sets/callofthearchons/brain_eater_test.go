package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Brain Eater draws a card after a creature is destroyed fighting it.
func TestBrainEater(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	g.AddToDeck(cardtest.Vanilla("d", engine.Logos, 1), 0)
	eater := g.AddToBattleline(BrainEater, 0)
	prey := g.AddToBattleline(cardtest.Vanilla("Prey", engine.Mars, 1), 1)

	before := len(g.Hand(0))
	if err := g.Fight(0, eater, prey); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("the 1-power prey should be destroyed")
	}
	if got := len(g.Hand(0)); got != before+1 {
		t.Errorf("hand = %d, want %d (drew a card)", got, before+1)
	}
}
