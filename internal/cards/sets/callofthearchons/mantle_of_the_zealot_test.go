package callofthearchons

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Mantle of the Zealot grants versatile, which lets its host be used as if it
// belonged to the active house: a creature of another house can reap once it
// wears the Mantle, even though the active house does not match its own.
func TestMantleOfTheZealot(t *testing.T) {
	if got := engine.RenderCardText(&MantleOfTheZealot); !strings.Contains(got, "This creature gains versatile.") {
		t.Errorf("Mantle text = %q, want it to grant versatile", got)
	}

	g := cardtest.Started(t, engine.Sanctum)
	off := g.AddToBattleline(cardtest.Vanilla("Offhouse", engine.Dis, 3), 0)
	g.AddToHand(MantleOfTheZealot, 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}
	if err := g.Reap(0, off); err != nil {
		t.Fatalf("reap of Versatile-granted creature out of house: %v", err)
	}
}
