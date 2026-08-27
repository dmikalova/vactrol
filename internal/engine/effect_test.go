package engine

import "testing"

func TestEffectValidation(t *testing.T) {
	bad := Heal{Fully: true, Amount: 1, Target: Target{Kind: TargetThisCreature}}
	good := Heal{Fully: true, Target: Target{Kind: TargetThisCreature}}

	if err := validateEffect(GainAember{Amount: 1}); err != nil {
		t.Errorf("non-validating effect should be nil, got %v", err)
	}
	if err := (Sequence{Effects: []Effect{good, bad}}).validate(); err == nil {
		t.Error("sequence should surface a bad child")
	}
	if err := (Sequence{Effects: []Effect{good, GainAember{Amount: 1}}}).validate(); err != nil {
		t.Errorf("sequence of valid effects should pass, got %v", err)
	}
	if err := (Conditional{Then: bad}).validate(); err == nil {
		t.Error("conditional should surface a bad gated effect")
	}
	if err := (Conditional{Then: GainAember{Amount: 1}}).validate(); err != nil {
		t.Errorf("conditional with a valid effect should pass, got %v", err)
	}
}

func TestNewCardRejectsConflictingHeal(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should panic on a Heal with both Amount and Fully")
		}
	}()
	NewCard("bad", Sanctum, Creature, Common,
		WithAbility(TriggerAfterPlay, Heal{Amount: 2, Fully: true, Target: Target{Kind: TargetThisCreature}}))
}
