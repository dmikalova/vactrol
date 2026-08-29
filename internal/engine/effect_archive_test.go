package engine

import "testing"

func TestArchiveEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	c1 := g.AddToHand(testCreature("c1", 1), 0)
	c2 := g.AddToHand(testCreature("c2", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (ArchiveFromHand{Count: 1}).Text() != "archive a card from your hand" {
		t.Errorf("archive text = %q", (ArchiveFromHand{Count: 1}).Text())
	}
	if (ArchiveFromHand{Count: 2}).Text() != "archive 2 cards from your hand" {
		t.Errorf("archive plural text = %q", (ArchiveFromHand{Count: 2}).Text())
	}

	// The default chooser archives the first hand card (c1).
	(ArchiveFromHand{Count: 1}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != c1 {
		t.Errorf("archives = %v, want [%d]", g.State.Archives[0].slice(), c1)
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != c2 {
		t.Errorf("hand = %v, want [%d]", g.Hand(0), c2)
	}

	// Archiving more than the hand holds stops when the hand empties.
	(ArchiveFromHand{Count: 5}).Resolve(ctx)
	if len(g.Hand(0)) != 0 {
		t.Errorf("hand should be empty, got %v", g.Hand(0))
	}
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
}

func TestArchiveEffectDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Two cards, so declining is a real choice (a sole card would be auto-archived).
	g.AddToHand(testCreature("c", 1), 0)
	g.AddToHand(testCreature("d", 1), 0)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Controller: 0}
	(ArchiveFromHand{Count: 1}).Resolve(ctx)
	if g.State.Archives[0].Count != 0 {
		t.Error("a declined archive choice should archive nothing")
	}
}

func TestArchiveTopOfDeckEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	top := g.AddToDeck(testCreature("top", 1), 0)
	g.AddToDeck(testCreature("next", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (ArchiveTopOfDeck{Count: 1}).Text() != "archive the top card of your deck" {
		t.Errorf("text = %q", (ArchiveTopOfDeck{Count: 1}).Text())
	}
	if (ArchiveTopOfDeck{Count: 2}).Text() != "archive the top 2 cards of your deck" {
		t.Errorf("plural text = %q", (ArchiveTopOfDeck{Count: 2}).Text())
	}

	(ArchiveTopOfDeck{Count: 1}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != top {
		t.Errorf("archived %v, want the top card %d", g.State.Archives[0].slice(), top)
	}

	// Archiving more than the deck holds stops when the deck empties.
	(ArchiveTopOfDeck{Count: 5}).Resolve(ctx)
	if g.State.Deck[0].Count != 0 {
		t.Errorf("deck should be empty, got %d", g.State.Deck[0].Count)
	}
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
}

func TestArchivesOfferedOnChooseHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToHand(testCreature("a", 1), 0)
	(ArchiveFromHand{Count: 1}).Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.Archives[0].Count != 1 {
		t.Fatalf("setup: archives count = %d, want 1", g.State.Archives[0].Count)
	}

	// The archives are NOT taken at the start of the turn.
	g.BeginTurn(0)
	if g.State.Archives[0].Count != 1 {
		t.Error("archives should not return before the house is chosen")
	}

	// Choosing a house offers them; the default chooser accepts (takes them).
	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if g.State.Archives[0].Count != 0 {
		t.Error("accepting the offer should empty the archives")
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != a {
		t.Errorf("hand = %v, want the returned card [%d]", g.Hand(0), a)
	}
}

func TestArchivesDeclinedOnChooseHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToHand(testCreature("a", 1), 0)
	(ArchiveFromHand{Count: 1}).Resolve(&EffectContext{Resolver: g, Controller: 0})
	g.SetChooser(0, optionPicker{idx: 1}) // decline the offer

	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if g.State.Archives[0].Count != 1 {
		t.Error("declining the offer should leave the cards archived")
	}
	if len(g.Hand(0)) != 0 {
		t.Errorf("hand = %v, want empty (offer declined)", g.Hand(0))
	}
}
