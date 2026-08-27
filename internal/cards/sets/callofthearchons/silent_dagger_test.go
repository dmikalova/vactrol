package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Silent Dagger grants its host "Reap: Deal 4 damage to a flank creature", so
// reaping with the host deals 4 to a chosen enemy on a flank.
func TestSilentDagger(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	host := g.AddToBattleline(cardtest.Vanilla("Host", engine.Shadows, 3), 0)
	g.AddToHand(SilentDagger, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil { // only candidate is the friendly host
		t.Fatalf("PlayUpgrade: %v", err)
	}
	foe := g.AddToBattleline(cardtest.Vanilla("Foe", engine.Mars, 5), 1)
	g.SetChooser(0, flankPicker{id: foe}) // a lone creature is on a flank

	if err := g.Reap(0, host); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.Damage(foe) != 4 {
		t.Errorf("foe damage = %d, want 4", g.Damage(foe))
	}
}
