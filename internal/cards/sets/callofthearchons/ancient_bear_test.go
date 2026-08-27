package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ancient Bear has Assault 2: it deals 2 damage to the creature it attacks
// before fight damage is exchanged.
func TestAncientBear(t *testing.T) {
	g := cardtest.Started(t, engine.Untamed)
	bear := g.AddToBattleline(AncientBear, 0) // 5 power, Assault 2
	foe := g.AddToBattleline(cardtest.Vanilla("Foe", engine.Mars, 10), 1)

	if err := g.Fight(0, bear, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	// Foe (10 power) survives, taking Assault 2 + 5 fight damage.
	if g.Damage(foe) != 7 {
		t.Errorf("foe damage = %d, want 7 (2 assault + 5 fight)", g.Damage(foe))
	}
	// The bear no longer has Skirmish, so it takes the 10-power foe's retaliation.
	if len(g.Battleline(0)) != 0 {
		t.Error("bear should be destroyed by the 10-power foe's retaliation")
	}
}
