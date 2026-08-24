package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// BrobnarGiant deals 2 damage to each enemy creature after its controller forges
// a key.
func TestBrobnarGiant(t *testing.T) {
	g := game.NewGame("A", "B", 1)
	g.AddToBattleline(BrobnarGiant, 0)
	tough := g.AddToBattleline(cardtest.Vanilla("Tough", game.Mars, 10), 1)
	g.AddToBattleline(cardtest.Vanilla("Weakling", game.Mars, 2), 1)

	g.State.Aember[0] = game.KeyCost // exactly enough for one key
	g.BeginTurn(0)                   // forging fires the giant's after-forge ability

	if g.Keys(0) != 1 {
		t.Fatalf("keys = %d, want 1", g.Keys(0))
	}
	if g.Damage(tough) != 2 {
		t.Errorf("tough enemy damage = %d, want 2", g.Damage(tough))
	}
	// The 2-power weakling takes 2 and is destroyed, leaving only the tough one.
	if bl := g.Battleline(1); len(bl) != 1 || bl[0] != tough {
		t.Errorf("enemy battleline = %v, want just the tough creature", bl)
	}
}
