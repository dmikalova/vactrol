package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Blood of Titans is an upgrade that grants its host creature +5 power while
// attached.
func TestBloodOfTitans(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)
	host := g.AddToBattleline(cardtest.Vanilla("Host", engine.Brobnar, 4), 0)
	before := g.Power(host)

	g.AddToHand(BloodOfTitans, 0)
	upHost, err := g.PlayUpgrade(0, 0)
	if err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}
	if upHost != host {
		t.Fatalf("upgrade attached to %v, want host %v", upHost, host)
	}
	if got := g.Power(host); got != before+5 {
		t.Errorf("host power = %d, want %d", got, before+5)
	}
}
