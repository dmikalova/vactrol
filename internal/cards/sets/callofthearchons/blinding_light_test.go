package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// housePicker chooses a fixed house option (and never a creature).
type housePicker struct{ idx int }

func (housePicker) ChooseCreature(string, []engine.LocalID) (engine.LocalID, bool) {
	return 0, false
}
func (h housePicker) ChooseOption(string, []string) int { return h.idx }

// Blinding Light stuns each creature of a chosen house, sparing the rest.
func TestBlindingLight(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	marsFoe := g.AddToBattleline(cardtest.Vanilla("MarsFoe", engine.Mars, 3), 1)
	shadowFoe := g.AddToBattleline(cardtest.Vanilla("ShadowFoe", engine.Shadows, 3), 1)
	g.AddToHand(BlindingLight, 0)
	g.SetChooser(0, housePicker{idx: 3}) // Mars is the fourth house

	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if !g.State.Cards[marsFoe].Stunned {
		t.Error("creatures of the chosen house (Mars) should be stunned")
	}
	if g.State.Cards[shadowFoe].Stunned {
		t.Error("creatures of other houses should not be stunned")
	}
}
