package engine

import "testing"

func TestReturnToDeckEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	myArt := g.AddArtifact(NewCard("myrelic", Brobnar, Artifact, Rare), 0)
	enemyArt := g.AddArtifact(NewCard("enemyrelic", Brobnar, Artifact, Rare), 1)
	g.State.Cards[myArt].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := ReturnToDeck{Target: Target{Kind: TargetEachArtifact}}
	if e.Text() != "put each artifact on top of its owner's deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Artifacts[0].Count != 0 || g.State.Artifacts[1].Count != 0 {
		t.Errorf("artifact rows not cleared: %d %d", g.State.Artifacts[0].Count, g.State.Artifacts[1].Count)
	}
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != myArt {
		t.Errorf("player 0 deck top = %v, want %d", g.State.Deck[0].IDs[0], myArt)
	}
	if g.State.Deck[1].Count != 1 || g.State.Deck[1].IDs[0] != enemyArt {
		t.Errorf("player 1 deck top = %v, want %d", g.State.Deck[1].IDs[0], enemyArt)
	}
	if g.State.Cards[myArt].Exhausted {
		t.Errorf("returned artifact should be readied")
	}
}

func TestReturnToHandEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.State.Cards[src].Damage = 2
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := ReturnToHand{Target: Target{Kind: TargetThisCreature}}
	if e.Text() != "put "+SelfName+" into its owner's hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Battleline[0].Count != 0 {
		t.Errorf("battleline not cleared: %d", g.State.Battleline[0].Count)
	}
	if g.State.Hand[0].Count != 1 || g.State.Hand[0].IDs[0] != src {
		t.Errorf("hand = count %d id %v, want 1 / %d", g.State.Hand[0].Count, g.State.Hand[0].IDs[0], src)
	}
	// Returning clears the per-match state the card accrued in play.
	if g.State.Cards[src].Damage != 0 {
		t.Errorf("damage after return = %d, want 0", g.State.Cards[src].Damage)
	}
}

func TestReturnToArchivesEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	attachUpgrade(g, src, NewCard("plating", Mars, Upgrade, Common))
	g.State.Cards[src].Damage = 1
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := ReturnToArchives{Target: Target{Kind: TargetThisCreature}}
	if e.Text() != "put "+SelfName+" into its owner's archives" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.inPlay(src) {
		t.Error("the creature should have left play")
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != src {
		t.Errorf("creature should be archived, got %v", g.State.Archives[0].IDs[:g.State.Archives[0].Count])
	}
	if len(g.Discard(0)) != 1 { // the upgrade sheds to the discard pile
		t.Errorf("upgrade should be discarded; discard = %v", g.Discard(0))
	}
	if g.State.Cards[src].Damage != 0 {
		t.Error("archiving should clear the creature's in-play state")
	}
}
