package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Fear puts a chosen enemy creature back into its owner's hand.
func TestFear(t *testing.T) {
	g := cardtest.Started(t, engine.Dis)
	g.AddToBattleline(cardtest.Vanilla("Enemy", engine.Mars, 3), 1)
	g.AddToHand(Fear, 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("enemy creature should have left play")
	}
	if g.State.Hand[1].Count != 1 {
		t.Errorf("enemy hand = %d, want 1 (creature returned)", g.State.Hand[1].Count)
	}
}
