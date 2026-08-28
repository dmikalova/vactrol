package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Biomatrix Backup relocates its destroyed host into its owner's archives instead
// of the discard pile; the upgrade itself is discarded.
func TestBiomatrixBackup(t *testing.T) {
	g := cardtest.Started(t, engine.Mars)
	host := g.AddToBattleline(cardtest.Vanilla("Host", engine.Mars, 3), 0)
	g.AddToHand(BiomatrixBackup, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil { // attaches to the friendly host
		t.Fatalf("PlayUpgrade: %v", err)
	}

	g.DestroyEach(0, []engine.LocalID{host})
	if len(g.Battleline(0)) != 0 {
		t.Error("host should have left the battleline")
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != host {
		t.Errorf("host should be archived, got %v", g.State.Archives[0].IDs[:g.State.Archives[0].Count])
	}
	if len(g.Discard(0)) != 1 { // the upgrade is discarded, the host is not
		t.Errorf("discard = %v, want just the upgrade", g.Discard(0))
	}
}
