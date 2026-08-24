package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// Anger readies a friendly creature and fights with it, so a creature that is
// already exhausted can still attack.
func TestAnger(t *testing.T) {
	g := cardtest.Started(t, game.Brobnar)
	friend := g.AddToBattleline(cardtest.Vanilla("Friend", game.Brobnar, 4), 0)
	g.State.Cards[friend].Exhausted = true // Anger should ready it before fighting
	g.AddToBattleline(cardtest.Vanilla("Foe", game.Mars, 2), 1)

	g.AddToHand(Anger, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}

	if len(g.Battleline(1)) != 0 {
		t.Errorf("foe should be destroyed by the ready-and-fight, battleline = %v", g.Battleline(1))
	}
	if g.Damage(friend) != 2 {
		t.Errorf("friend took %d retaliation damage, want 2", g.Damage(friend))
	}
}
