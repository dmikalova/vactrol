package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Hand of Dis destroys a chosen creature that is not on a flank, so only the
// interior of a three-creature line is a legal target.
func TestHandOfDis(t *testing.T) {
	g := cardtest.Started(t, engine.Dis)
	left := g.AddToBattleline(cardtest.Vanilla("Left", engine.Mars, 3), 1)
	mid := g.AddToBattleline(cardtest.Vanilla("Mid", engine.Mars, 3), 1)
	right := g.AddToBattleline(cardtest.Vanilla("Right", engine.Mars, 3), 1)
	g.AddToHand(HandOfDis, 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	bl := g.Battleline(1)
	if len(bl) != 2 || bl[0] != left || bl[1] != right {
		t.Errorf("battleline = %v, want the two flanks [%d %d] (mid %d destroyed)", bl, left, right, mid)
	}
}
