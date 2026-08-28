package engine

import "testing"

func TestReturnFromDiscardToHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.Register(testCreature("c", 3), 0)
	g.State.Discard[0].add(c)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := ReturnFromDiscard{}
	if e.Text() != "put a card from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != c {
		t.Errorf("hand = %v, want [%d]", g.Hand(0), c)
	}
	if len(g.Discard(0)) != 0 {
		t.Error("the card should have left the discard pile")
	}

	// Empty discard: no candidate, so nothing happens.
	e.Resolve(ctx)
	if len(g.Hand(0)) != 1 {
		t.Error("an empty discard should return nothing")
	}
}

func TestReturnCreatureFromDiscardToDeck(t *testing.T) {
	g := NewGame("A", "B", 1)
	act := g.Register(NewCard("act", Brobnar, Action, Common), 0) // non-creature
	g.State.Discard[0].add(act)
	crea := g.Register(testCreature("crea", 3), 0)
	g.State.Discard[0].add(crea)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := ReturnFromDiscard{CreaturesOnly: true, ToDeck: true}
	if e.Text() != "put a creature from your discard pile on top of your deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	// The action is skipped (not a creature); the creature goes to the deck top.
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != crea {
		t.Errorf("deck top = %v, want the creature %d", g.State.Deck[0].IDs[:g.State.Deck[0].Count], crea)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != act {
		t.Errorf("discard = %v, want just the action %d", d, act)
	}
}
