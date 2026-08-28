package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// EMP Blast stuns each Mars creature and each Robot, and destroys every artifact.
func TestEMPBlast(t *testing.T) {
	g := cardtest.Started(t, engine.Mars)
	marsGuy := g.AddToBattleline(cardtest.Vanilla("MarsGuy", engine.Mars, 3), 1)
	robot := g.AddToBattleline(engine.NewCard("Bot", engine.Logos, engine.Creature, engine.Common,
		engine.WithPower(3), engine.WithTraits("Robot")), 1)
	other := g.AddToBattleline(cardtest.Vanilla("Other", engine.Sanctum, 3), 1)
	g.AddArtifact(Cannon, 1)
	g.AddToHand(EMPBlast, 0)

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if !g.State.Cards[marsGuy].Stunned {
		t.Error("Mars creature should be stunned")
	}
	if !g.State.Cards[robot].Stunned {
		t.Error("Robot creature should be stunned")
	}
	if g.State.Cards[other].Stunned {
		t.Error("a non-Mars, non-Robot creature should not be stunned")
	}
	if len(g.Artifacts(1)) != 0 {
		t.Error("each artifact should be destroyed")
	}
}
