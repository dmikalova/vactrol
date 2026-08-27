package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// DocBookton draws a card when it reaps.
func TestDocBookton(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	g.AddToDeck(cardtest.Vanilla("d", engine.Logos, 1), 0)
	id := g.AddToBattleline(DocBookton, 0)

	before := len(g.Hand(0))
	if err := g.Reap(0, id); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if got := len(g.Hand(0)); got != before+1 {
		t.Errorf("hand = %d, want %d (drew a card)", got, before+1)
	}
	if g.Power(id) != 5 {
		t.Errorf("power = %d, want 5", g.Power(id))
	}
}
