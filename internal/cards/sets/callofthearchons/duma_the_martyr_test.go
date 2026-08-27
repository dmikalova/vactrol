package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Duma the Martyr fully heals each other friendly creature and draws 2 cards
// when it is destroyed.
func TestDumaTheMartyr(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	duma := g.AddToBattleline(DumaTheMartyr, 0)
	other := g.AddToBattleline(cardtest.Vanilla("ally", engine.Sanctum, 4), 0)
	g.State.Cards[other].Damage = 2
	g.AddToDeck(cardtest.Vanilla("d1", engine.Sanctum, 1), 0)
	g.AddToDeck(cardtest.Vanilla("d2", engine.Sanctum, 1), 0)
	enemy := g.AddToBattleline(cardtest.Vanilla("foe", engine.Sanctum, 3), 1)

	handBefore := len(g.Hand(0))
	if err := g.Fight(0, duma, enemy); err != nil {
		t.Fatalf("Fight: %v", err)
	}

	if contains(g.Battleline(0), duma) {
		t.Error("Duma (3 power) should die to 3 return damage")
	}
	if g.Damage(other) != 0 {
		t.Errorf("ally damage = %d, want 0 (fully healed)", g.Damage(other))
	}
	if got := len(g.Hand(0)); got != handBefore+2 {
		t.Errorf("hand = %d, want %d (drew 2)", got, handBefore+2)
	}
}
