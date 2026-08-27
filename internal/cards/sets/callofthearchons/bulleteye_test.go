package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// flankPicker is a chooser that always picks a fixed creature.
type flankPicker struct{ id engine.LocalID }

func (f flankPicker) ChooseCreature(string, []engine.LocalID) (engine.LocalID, bool) {
	return f.id, true
}

// Bulleteye's Reap destroys a chosen flank creature.
func TestBulleteye(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	eye := g.AddToBattleline(Bulleteye, 0)
	foe := g.AddToBattleline(cardtest.Vanilla("Foe", engine.Mars, 3), 1)
	g.SetChooser(0, flankPicker{id: foe}) // a lone creature is on a flank

	if err := g.Reap(0, eye); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("the flank creature should be destroyed")
	}
}
