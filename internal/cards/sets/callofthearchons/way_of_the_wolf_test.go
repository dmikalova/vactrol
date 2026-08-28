package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Way of the Wolf grants its host skirmish, so the host takes no damage in return
// when it is used to fight.
func TestWayOfTheWolf(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	host := g.AddToBattleline(cardtest.Vanilla("Host", engine.Untamed, 4), 0)
	g.AddToHand(WayOfTheWolf, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil { // only candidate is the friendly host
		t.Fatalf("PlayUpgrade: %v", err)
	}
	wall := g.AddToBattleline(cardtest.Vanilla("Wall", engine.Mars, 10), 1)
	if err := g.Fight(0, host, wall); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if len(g.Battleline(0)) != 1 {
		t.Error("granted skirmish should spare the host from the 10-power wall")
	}
	if g.Damage(host) != 0 {
		t.Errorf("host took %d damage, want 0 (skirmish)", g.Damage(host))
	}
}
