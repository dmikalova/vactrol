package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ritual of Balance is an artifact whose Action steals 1 Æmber, but only while
// the opponent is hoarding 6 or more.
func TestRitualOfBalance(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	art := g.AddArtifact(RitualOfBalance, 0)

	// Below the threshold: the action does nothing.
	g.State.Aember[1] = 5
	if err := g.UseAction(0, art); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Aember(0) != 0 || g.Aember(1) != 5 {
		t.Errorf("below threshold: pools = %d/%d, want 0/5", g.Aember(0), g.Aember(1))
	}

	// Ready it again and push the opponent to the threshold: now it steals 1.
	g.State.Cards[art].Exhausted = false
	g.State.Aember[1] = 6
	if err := g.UseAction(0, art); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Aember(0) != 1 || g.Aember(1) != 5 {
		t.Errorf("at threshold: pools = %d/%d, want 1/5", g.Aember(0), g.Aember(1))
	}
}
