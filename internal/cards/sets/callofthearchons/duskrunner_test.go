package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Duskrunner grants its host "Reap: Steal 1 Æmber", so reaping with the host
// steals on top of the Æmber the reap itself yields.
func TestDuskrunner(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	host := g.AddToBattleline(cardtest.Vanilla("Host", engine.Shadows, 3), 0)
	g.AddToHand(Duskrunner, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil { // only candidate is the friendly host
		t.Fatalf("PlayUpgrade: %v", err)
	}
	g.State.Aember[1] = 2 // opponent has Æmber to steal

	if err := g.Reap(0, host); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("controller aember = %d, want 2 (1 reap + 1 stolen)", g.Aember(0))
	}
	if g.Aember(1) != 1 {
		t.Errorf("opponent aember = %d, want 1", g.Aember(1))
	}
}
