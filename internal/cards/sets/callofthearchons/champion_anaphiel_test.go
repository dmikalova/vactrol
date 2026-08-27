package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Champion Anaphiel is a 6-power Knight with Taunt, protecting its smaller
// neighbors from being fought directly.
func TestChampionAnaphiel(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	g.AddToHand(ChampionAnaphiel, 0)
	id, err := g.PlayCreature(0, 0, false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Power(id) != 6 {
		t.Errorf("Champion Anaphiel power = %d, want 6", g.Power(id))
	}

	hasTaunt := false
	for _, kw := range ChampionAnaphiel.Keywords {
		if kw == engine.Taunt {
			hasTaunt = true
		}
	}
	if !hasTaunt {
		t.Errorf("Champion Anaphiel keywords = %v, want Taunt", ChampionAnaphiel.Keywords)
	}
}
