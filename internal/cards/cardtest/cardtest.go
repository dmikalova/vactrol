// Package cardtest is a shared test harness for card-set packages. Each set's
// _test.go files import it so the same game-setup helpers work across every set
// (callofthearchons, and future sets) instead of being copied per package.
package cardtest

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Started returns a new game with player 0 active and the given house chosen,
// ready to play that player's cards of that house.
func Started(t *testing.T, house engine.House) *engine.Game {
	t.Helper()
	g := engine.NewGame("A", "B", 1)
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, house); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	return g
}

// Vanilla builds a plain creature (no abilities) of the given house and power,
// used as a supporting body or opponent in card scenarios.
func Vanilla(name string, house engine.House, power int) engine.CardDefinition {
	return engine.NewCard(name, house, engine.Creature, engine.Common, engine.WithPower(power))
}
