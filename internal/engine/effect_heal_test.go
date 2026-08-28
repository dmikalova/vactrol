package engine

import "testing"

func TestHealEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 5), 0)
	other := g.AddToBattleline(testCreature("other", 5), 0)
	g.State.Cards[src].Damage = 3
	g.State.Cards[other].Damage = 4
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	partial := Heal{Amount: 2, Target: Target{Kind: TargetThisCreature}}
	if partial.Text() != "heal 2 damage from "+SelfName {
		t.Errorf("partial text = %q", partial.Text())
	}
	partial.Resolve(ctx)
	if g.State.Cards[src].Damage != 1 {
		t.Errorf("partial heal: src damage = %d, want 1", g.State.Cards[src].Damage)
	}
	(Heal{Amount: 5, Target: Target{Kind: TargetThisCreature}}).Resolve(ctx) // over-heal floors at 0
	if g.State.Cards[src].Damage != 0 {
		t.Errorf("over-heal should floor at 0, got %d", g.State.Cards[src].Damage)
	}

	full := Heal{Fully: true, Target: Target{Kind: TargetEachOtherFriendlyCreature}}
	if full.Text() != "fully heal each other friendly creature" {
		t.Errorf("full text = %q", full.Text())
	}
	full.Resolve(ctx)
	if g.State.Cards[other].Damage != 0 {
		t.Errorf("full heal: other damage = %d, want 0", g.State.Cards[other].Damage)
	}
}

func TestHealValidate(t *testing.T) {
	if err := (Heal{Fully: true, Amount: 2}).validate(); err == nil {
		t.Error("Heal with both Amount and Fully should be invalid")
	}
	if err := (Heal{Fully: true}).validate(); err != nil {
		t.Errorf("full heal should be valid, got %v", err)
	}
	if err := (Heal{Amount: 2}).validate(); err != nil {
		t.Errorf("fixed heal should be valid, got %v", err)
	}
}
