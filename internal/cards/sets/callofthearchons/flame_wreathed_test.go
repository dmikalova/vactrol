package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Flame-Wreathed grants its host +2 power and +2 hazardous. Here it is attached
// to an enemy creature, so an attacker takes the hazardous damage.
func TestFlameWreathed(t *testing.T) {
	g := cardtest.Started(t, engine.Dis)
	enemy := g.AddToBattleline(cardtest.Vanilla("Enemy", engine.Mars, 4), 1)
	g.AddToHand(FlameWreathed, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil { // only candidate is the enemy
		t.Fatalf("PlayUpgrade: %v", err)
	}
	if g.Power(enemy) != 6 {
		t.Errorf("enemy power = %d, want 6 (4 + 2)", g.Power(enemy))
	}

	attacker := g.AddToBattleline(cardtest.Vanilla("Attacker", engine.Dis, 1), 0)
	if err := g.Fight(0, attacker, enemy); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	// The 1-power attacker is destroyed by the granted Hazardous 2 before combat.
	if len(g.Battleline(0)) != 0 {
		t.Error("attacker should be destroyed by the granted Hazardous 2")
	}
}
