package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Barehanded returns every artifact in play — both players' — to the top of its
// owner's deck, and grants its controller 1 Æmber from its bonus pip.
func TestBarehanded(t *testing.T) {
	g := cardtest.Started(t, engine.Brobnar)

	mine := engine.NewCard("My Relic", engine.Brobnar, engine.Artifact, engine.Rare)
	theirs := engine.NewCard("Their Relic", engine.Untamed, engine.Artifact, engine.Rare)
	myArtifact := g.AddArtifact(mine, 0)
	theirArtifact := g.AddArtifact(theirs, 1)

	g.AddToHand(Barehanded, 0)
	if err := g.PlayAction(0, 0); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}

	// Both artifact rows are cleared.
	if got := g.Artifacts(0); len(got) != 0 {
		t.Errorf("friendly artifacts = %v, want empty", got)
	}
	if got := g.Artifacts(1); len(got) != 0 {
		t.Errorf("enemy artifacts = %v, want empty", got)
	}

	// Each artifact sits on top of its own owner's deck.
	if deck := g.State.Deck[0]; deck.Count != 1 || deck.IDs[0] != myArtifact {
		t.Errorf("player 0 deck top = %v (count %d), want %d", deck.IDs[0], deck.Count, myArtifact)
	}
	if deck := g.State.Deck[1]; deck.Count != 1 || deck.IDs[0] != theirArtifact {
		t.Errorf("player 1 deck top = %v (count %d), want %d", deck.IDs[0], deck.Count, theirArtifact)
	}

	// The Æmber bonus pip resolves on play.
	if g.Aember(0) != 1 {
		t.Errorf("controller Æmber = %d, want 1", g.Aember(0))
	}
}
