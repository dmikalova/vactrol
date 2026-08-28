package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Knowledge is Power's second option gains 1 Æmber for each archived card.
func TestKnowledgeIsPower(t *testing.T) {
	g := cardtest.Started(t, engine.Logos)
	// Two cards already banked in the archives.
	g.State.Archives[0].IDs[0] = g.Register(cardtest.Vanilla("b1", engine.Logos, 1), 0)
	g.State.Archives[0].IDs[1] = g.Register(cardtest.Vanilla("b2", engine.Logos, 1), 0)
	g.State.Archives[0].Count = 2
	g.AddToHand(KnowledgeIsPower, 0)
	g.SetChooser(0, housePicker{idx: 1}) // choose the "gain" option

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (1 per each of 2 archived cards)", g.Aember(0))
	}
}
