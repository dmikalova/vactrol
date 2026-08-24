package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards/cardtest"
)

// Ganymede Archivist gains an extra Æmber from its Reap ability, on top of the
// +1 for reaping.
func TestGanymedeArchivist(t *testing.T) {
	g := cardtest.Started(t, game.Logos)
	id := g.AddToBattleline(GanymedeArchivist, 0)

	if err := g.Reap(0, id); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (+1 reap, +1 ability)", g.Aember(0))
	}
	if !g.Exhausted(id) {
		t.Error("Ganymede Archivist should be exhausted after reaping")
	}
}
