package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Curiosity destroys each Scientist trait creature on Play, leaving other
// creatures untouched.
func TestCuriosity(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	sci := g.AddToBattleline(engine.NewCard("sci", engine.Logos, engine.Creature, engine.Common, engine.WithPower(3), engine.WithTraits("Scientist")), 1)
	beast := g.AddToBattleline(engine.NewCard("beast", engine.Untamed, engine.Creature, engine.Common, engine.WithPower(3), engine.WithTraits("Beast")), 1)
	g.AddToHand(Curiosity, 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}

	if contains(g.Battleline(1), sci) {
		t.Error("Scientist trait creature should be destroyed")
	}
	if !contains(g.Battleline(1), beast) {
		t.Error("Beast trait creature should survive")
	}
}
