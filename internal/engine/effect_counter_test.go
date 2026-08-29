package engine

import "testing"

func TestAddPowerCounter(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: c, Controller: 0}

	e := AddPowerCounter{Target: Target{Kind: TargetThisCreature}, Amount: 1}
	if e.Text() != "give {self} a +1 power counter" {
		t.Errorf("text = %q", e.Text())
	}

	e.Resolve(ctx)
	if g.Power(c) != 4 {
		t.Errorf("power = %d, want 4 (+1 counter)", g.Power(c))
	}
	// Counters stack.
	e.Resolve(ctx)
	if g.Power(c) != 5 {
		t.Errorf("power = %d, want 5 (two +1 counters)", g.Power(c))
	}
}
