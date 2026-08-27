package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Way of the Bear grants its host +2 assault, so the host deals 2 extra damage
// to whatever it attacks, before fight damage.
func TestWayOfTheBear(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	host := g.AddToBattleline(cardtest.Vanilla("Host", engine.Untamed, 3), 0)
	g.AddToHand(WayOfTheBear, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}
	foe := g.AddToBattleline(cardtest.Vanilla("Foe", engine.Mars, 10), 1)

	if err := g.Fight(0, host, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	// Host (3 power) with Assault 2: foe takes 2 assault + 3 fight = 5.
	if g.Damage(foe) != 5 {
		t.Errorf("foe damage = %d, want 5 (2 assault + 3 fight)", g.Damage(foe))
	}
}
