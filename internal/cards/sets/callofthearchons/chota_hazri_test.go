package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Chota Hazri loses 1 Æmber to forge a key at its current cost.
func TestChotaHazri(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	g.State.Aember[0] = engine.KeyCost + 1 // enough to pay the 1 and the key cost
	g.AddToHand(ChotaHazri, 0)

	if _, err := g.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1 (forged)", g.Keys(0))
	}
	if g.Aember(0) != 0 {
		t.Errorf("aember = %d, want 0 (1 lost + key cost paid)", g.Aember(0))
	}
}
