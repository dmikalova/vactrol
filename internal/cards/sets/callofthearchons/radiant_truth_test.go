package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Radiant Truth stuns each enemy creature not on a flank, sparing the flanks.
func TestRadiantTruth(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	left := g.AddToBattleline(cardtest.Vanilla("Left", engine.Mars, 3), 1)
	mid := g.AddToBattleline(cardtest.Vanilla("Mid", engine.Mars, 3), 1)
	right := g.AddToBattleline(cardtest.Vanilla("Right", engine.Mars, 3), 1)
	g.AddToHand(RadiantTruth, 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if !g.State.Cards[mid].Stunned {
		t.Error("the interior enemy creature should be stunned")
	}
	if g.State.Cards[left].Stunned || g.State.Cards[right].Stunned {
		t.Error("flank creatures should not be stunned")
	}
}
