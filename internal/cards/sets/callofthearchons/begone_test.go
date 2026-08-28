package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Begone!'s first option destroys each Dis creature (the default chooser takes
// the first option).
func TestBegone(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	g.AddToBattleline(cardtest.Vanilla("DisGuy", engine.Dis, 3), 1)
	other := g.AddToBattleline(cardtest.Vanilla("Other", engine.Mars, 3), 1)
	g.AddToHand(Begone, 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	bl := g.Battleline(1)
	if len(bl) != 1 || bl[0] != other {
		t.Errorf("battleline = %v, want only the non-Dis creature [%d]", bl, other)
	}
}
