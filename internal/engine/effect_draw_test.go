package engine

import "testing"

func TestDrawEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	g.AddToDeck(testCreature("d1", 1), 0)
	g.AddToDeck(testCreature("d2", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if (Draw{Amount: 1}).Text() != "draw a card" {
		t.Errorf("single text = %q", (Draw{Amount: 1}).Text())
	}
	e := Draw{Amount: 2}
	if e.Text() != "draw 2 cards" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Hand[0].Count != 2 {
		t.Errorf("hand = %d, want 2", g.State.Hand[0].Count)
	}
	(Draw{Amount: 3}).Resolve(ctx) // deck now empty; draws nothing
	if g.State.Hand[0].Count != 2 {
		t.Errorf("hand after empty-deck draw = %d, want 2", g.State.Hand[0].Count)
	}
}
