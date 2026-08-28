package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Longfused Mines sacrifices itself to deal 3 damage to each enemy creature not
// on a flank, sparing the flanks.
func TestLongfusedMines(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	mines := g.AddArtifact(LongfusedMines, 0)
	left := g.AddToBattleline(cardtest.Vanilla("Left", engine.Mars, 5), 1)
	mid := g.AddToBattleline(cardtest.Vanilla("Mid", engine.Mars, 5), 1)
	right := g.AddToBattleline(cardtest.Vanilla("Right", engine.Mars, 5), 1)

	if err := g.UseAction(0, mines); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Damage(mid) != 3 {
		t.Errorf("interior enemy damage = %d, want 3", g.Damage(mid))
	}
	if g.Damage(left) != 0 || g.Damage(right) != 0 {
		t.Error("flank enemies should be untouched")
	}
	if len(g.Artifacts(0)) != 0 {
		t.Error("Longfused Mines should sacrifice itself")
	}
}
