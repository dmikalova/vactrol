package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Dr. Escotera gains 1 Æmber for each key the opponent has forged.
func TestDrEscotera(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	g.State.Keys[1] = 2 // opponent has forged 2 keys
	g.AddToHand(DrEscotera, 0)

	if _, err := g.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (1 per each of 2 forged keys)", g.Aember(0))
	}
}
