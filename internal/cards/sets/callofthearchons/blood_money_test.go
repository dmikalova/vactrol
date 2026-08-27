package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// BloodMoney exalts a chosen enemy creature twice, leaving 2 Æmber on it.
func TestBloodMoney(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	foe := g.AddToBattleline(cardtest.Vanilla("Foe", engine.Mars, 4), 1)

	g.AddToHand(BloodMoney, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if got := g.State.Cards[foe].Amber; got != 2 {
		t.Errorf("Æmber on foe = %d, want 2", got)
	}
}
