package engine

import "testing"

func TestThen(t *testing.T) {
	// A gate whose action happens runs the Result.
	g := NewGame("A", "B", 1)
	self := g.AddToBattleline(testCreature("self", 3), 0)
	foe := g.AddToBattleline(testCreature("foe", 1), 1)
	ctx := &EffectContext{Resolver: g, Source: self, Controller: 0}

	then := Then{
		First:  Destroy{Target: Target{Kind: TargetEachEnemyCreature}},
		Result: AddPowerCounter{Target: Target{Kind: TargetThisCreature}, Amount: 1},
	}
	if then.Text() != "destroy each enemy creature -> give {self} a +1 power counter" {
		t.Errorf("text = %q", then.Text())
	}

	then.Resolve(ctx)
	if g.inPlay(foe) {
		t.Error("the gate should destroy the enemy creature")
	}
	if g.Power(self) != 4 {
		t.Errorf("the result should run: power = %d, want 4", g.Power(self))
	}

	// A gate whose action does nothing skips the Result.
	then.Resolve(ctx) // no enemy creatures remain
	if g.Power(self) != 4 {
		t.Errorf("the result should not run again: power = %d, want 4", g.Power(self))
	}
}
