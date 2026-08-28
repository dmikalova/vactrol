package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ring of Invisibility grants elusive and skirmish; the skirmish spares its host
// from retaliation when it fights.
func TestRingOfInvisibility(t *testing.T) {
	g := cardtest.Started(t, engine.Shadows)
	host := g.AddToBattleline(cardtest.Vanilla("Host", engine.Shadows, 4), 0)
	g.AddToHand(RingOfInvisibility, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}
	wall := g.AddToBattleline(cardtest.Vanilla("Wall", engine.Mars, 10), 1)
	if err := g.Fight(0, host, wall); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.Damage(host) != 0 {
		t.Errorf("host took %d damage, want 0 (granted skirmish)", g.Damage(host))
	}
}
