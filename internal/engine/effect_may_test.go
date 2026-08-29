package engine

import "testing"

func TestMayText(t *testing.T) {
	e := May{Do: GainAember{Player: Controller, Amount: 1}}
	if got := e.Text(); got != "you may gain 1 Æmber" {
		t.Errorf("text = %q", got)
	}
}

func TestMayResolvesWhenAccepted(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 0}) // "Yes"
	ctx := &EffectContext{Resolver: g, Controller: 0}

	May{Do: GainAember{Player: Controller, Amount: 2}}.Resolve(ctx)
	if g.Aember(0) != 2 {
		t.Errorf("accepted: aember = %d, want 2", g.Aember(0))
	}
}

func TestMayDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 1}) // "No"
	ctx := &EffectContext{Resolver: g, Controller: 0}

	May{Do: GainAember{Player: Controller, Amount: 2}}.Resolve(ctx)
	if g.Aember(0) != 0 {
		t.Errorf("declined: aember = %d, want 0", g.Aember(0))
	}
}

func TestMayValidate(t *testing.T) {
	bad := May{Do: Heal{Fully: true, Amount: 1, Target: Target{Kind: TargetThisCreature}}}
	if validateEffect(bad) == nil {
		t.Error("May wrapping an invalid effect should fail validation")
	}
	good := May{Do: GainAember{Player: Controller, Amount: 1}}
	if validateEffect(good) != nil {
		t.Error("May wrapping a valid effect should pass validation")
	}
}
