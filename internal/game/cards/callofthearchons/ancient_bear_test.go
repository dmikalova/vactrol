package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// AncientBear has Skirmish, so it takes no retaliation damage when it fights.
func TestAncientBear(t *testing.T) {
	g := cardtest.Started(t, game.Untamed)
	bear := g.AddToBattleline(AncientBear, 0) // 5 power, Skirmish
	foe := g.AddToBattleline(cardtest.Vanilla("Foe", game.Mars, 3), 1)

	if err := g.Fight(0, bear, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.Damage(bear) != 0 {
		t.Errorf("skirmisher took %d damage, want 0", g.Damage(bear))
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("foe should be destroyed by the 5-power bear")
	}
}
