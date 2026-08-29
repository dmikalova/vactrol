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

func TestRepeatWhile(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0], g.State.Aember[1] = 0, 5 // opponent leads
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := RepeatWhile{Cond: OpponentAemberMoreThanYou{}, Do: StealAember{Amount: 1}}
	if e.Text() != "if your opponent has more Æmber than you, steal 1 Æmber -> repeat this effect" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	// 5/0 -> 4/1 -> 3/2 -> 2/3 (opponent no longer leads): 3 steals.
	if g.Aember(0) != 3 || g.Aember(1) != 2 {
		t.Errorf("after repeat: you=%d opp=%d, want 3/2", g.Aember(0), g.Aember(1))
	}

	// Condition false from the start: the loop never runs.
	g2 := NewGame("A", "B", 1)
	g2.State.Aember[0], g2.State.Aember[1] = 3, 3
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Aember(0) != 3 || g2.Aember(1) != 3 {
		t.Errorf("equal pools should not steal: you=%d opp=%d", g2.Aember(0), g2.Aember(1))
	}

	if err := validateEffect(RepeatWhile{Do: StealAember{Amount: 1}}); err != nil {
		t.Errorf("validate = %v", err)
	}
}
