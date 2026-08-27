package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Champion Tabris captures 1 Æmber from the opponent's pool onto itself each
// time it fights.
func TestChampionTabris(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	tabris := g.AddToBattleline(ChampionTabris, 0)
	enemy := g.AddToBattleline(cardtest.Vanilla("Foe", engine.Brobnar, 3), 1)
	g.State.Aember[1] = 2

	if err := g.Fight(0, tabris, enemy); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.AmberOn(tabris) != 1 {
		t.Errorf("captured Æmber on Tabris = %d, want 1", g.AmberOn(tabris))
	}
	if g.Aember(1) != 1 {
		t.Errorf("opponent Æmber = %d, want 1", g.Aember(1))
	}
}
