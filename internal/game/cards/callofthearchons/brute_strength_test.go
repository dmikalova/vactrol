package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// BruteStrength gives its host creature +5 power.
func TestBruteStrength(t *testing.T) {
	g := cardtest.Started(t, game.Brobnar)
	host := g.AddToBattleline(cardtest.Vanilla("Host", game.Brobnar, 3), 0)

	g.AddToHand(BruteStrength, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil { // default chooser attaches to the only creature
		t.Fatalf("PlayUpgrade: %v", err)
	}
	if g.Power(host) != 8 {
		t.Errorf("host power = %d, want 8 (3 + 5)", g.Power(host))
	}
}
