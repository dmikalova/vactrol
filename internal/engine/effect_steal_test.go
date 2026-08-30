package engine

import "testing"

func TestStealAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 1

	e := StealAember{Amount: 3}
	if e.Text() != "steal 3 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx) // opponent has only 1
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 0 {
		t.Errorf("after steal: you=%d opp=%d, want 1/0", g.State.Aember[0], g.State.Aember[1])
	}
}
