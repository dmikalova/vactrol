package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Irradiated Aember deals 3 damage to each enemy creature only when the opponent
// hoards 6 Æmber or more; otherwise enemy creatures are untouched.
func TestIrradiatedAember(t *testing.T) {
	// Condition met: opponent at 6 — weak foe dies, tough foe survives marked.
	g := cardtest.Started(t, engine.Mars)
	g.State.Aember[1] = 6
	toughFoe := g.AddToBattleline(cardtest.Vanilla("Tough Foe", engine.Brobnar, 5), 1)
	g.AddToBattleline(cardtest.Vanilla("Weak Foe", engine.Brobnar, 2), 1)
	g.AddToHand(IrradiatedAember, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if got := g.Battleline(1); len(got) != 1 || got[0] != toughFoe {
		t.Errorf("enemy battleline = %v, want just the 5-power foe", got)
	}
	if g.Damage(toughFoe) != 3 {
		t.Errorf("tough foe damage = %d, want 3", g.Damage(toughFoe))
	}

	// Condition not met: opponent at 5 — every enemy creature is unharmed.
	g = cardtest.Started(t, engine.Mars)
	g.State.Aember[1] = 5
	survivor := g.AddToBattleline(cardtest.Vanilla("Weak Foe", engine.Brobnar, 2), 1)
	g.AddToHand(IrradiatedAember, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if got := g.Battleline(1); len(got) != 1 || got[0] != survivor {
		t.Errorf("enemy battleline = %v, want the untouched foe", got)
	}
	if g.Damage(survivor) != 0 {
		t.Errorf("foe damage = %d, want 0", g.Damage(survivor))
	}
}
