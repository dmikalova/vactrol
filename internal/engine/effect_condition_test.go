package engine

import "testing"

func TestConditionalEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	atLeast := Conditional{Cond: OpponentAemberAtLeast{Amount: 7}, Then: LoseAember{Player: Opponent, Amount: 4}}
	if atLeast.Text() != "if your opponent has 7 Æmber or more, your opponent loses 4 Æmber" {
		t.Errorf("at-least text = %q", atLeast.Text())
	}
	g.State.Aember[1] = 5
	atLeast.Resolve(ctx) // condition not met (5 < 7)
	if g.State.Aember[1] != 5 {
		t.Errorf("unmet condition should be a no-op; opp = %d", g.State.Aember[1])
	}
	g.State.Aember[1] = 8
	atLeast.Resolve(ctx) // met
	if g.State.Aember[1] != 4 {
		t.Errorf("met condition should apply; opp = %d, want 4", g.State.Aember[1])
	}

	exact := Conditional{Cond: OpponentAemberExactly{Amount: 1}, Then: StealAember{Amount: 1}}
	if exact.Text() != "if your opponent has exactly 1 Æmber, steal 1 Æmber" {
		t.Errorf("exact text = %q", exact.Text())
	}
	g.State.Aember[0], g.State.Aember[1] = 0, 2
	exact.Resolve(ctx) // 2 != 1, no-op
	if g.State.Aember[0] != 0 {
		t.Errorf("exact unmet should be a no-op; you = %d", g.State.Aember[0])
	}
	g.State.Aember[1] = 1
	exact.Resolve(ctx) // met
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 0 {
		t.Errorf("exact met: you=%d opp=%d, want 1/0", g.State.Aember[0], g.State.Aember[1])
	}
}
