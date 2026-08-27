package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Clear Mind removes the stun from every friendly creature, leaving the
// opponent's stunned creatures untouched.
func TestClearMind(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	ally1 := g.AddToBattleline(cardtest.Vanilla("ally1", engine.Sanctum, 3), 0)
	ally2 := g.AddToBattleline(cardtest.Vanilla("ally2", engine.Sanctum, 3), 0)
	foe := g.AddToBattleline(cardtest.Vanilla("foe", engine.Dis, 3), 1)
	g.State.Cards[ally1].Stunned = true
	g.State.Cards[ally2].Stunned = true
	g.State.Cards[foe].Stunned = true

	g.AddToHand(ClearMind, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}

	if g.State.Cards[ally1].Stunned || g.State.Cards[ally2].Stunned {
		t.Error("Clear Mind should unstun each friendly creature")
	}
	if !g.State.Cards[foe].Stunned {
		t.Error("Clear Mind should not affect enemy creatures")
	}
}
