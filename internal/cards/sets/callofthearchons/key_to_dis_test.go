package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Key to Dis sacrifices itself to destroy every creature in play.
func TestKeyToDis(t *testing.T) {
	g := cardtest.Started(t, engine.Dis)
	key := g.AddArtifact(KeyToDis, 0)
	g.AddToBattleline(cardtest.Vanilla("Friend", engine.Dis, 3), 0)
	g.AddToBattleline(cardtest.Vanilla("Foe", engine.Mars, 5), 1)

	if err := g.UseAction(0, key); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if len(g.Battleline(0)) != 0 || len(g.Battleline(1)) != 0 {
		t.Error("Key to Dis should destroy every creature")
	}
	if len(g.Artifacts(0)) != 0 {
		t.Error("Key to Dis should sacrifice itself")
	}
}
