package engine

import "testing"

func TestDealDamagePerCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Three friendly creatures multiply the 1 base damage to 3.
	g.AddToBattleline(testCreature("f1", 5), 0)
	g.AddToBattleline(testCreature("f2", 5), 0)
	g.AddToBattleline(testCreature("f3", 5), 0)
	victim := g.AddToBattleline(testCreature("victim", 10), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{Amount: 1, Per: FriendlyCreaturesInPlay{}, Target: Target{Kind: TargetEachEnemyCreature}}
	if e.Text() != "for each friendly creature in play, deal 1 damage to each enemy creature" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Damage(victim) != 3 {
		t.Errorf("victim damage = %d, want 3 (1 × 3 friendly creatures)", g.Damage(victim))
	}
}

func TestSplashDamage(t *testing.T) {
	g := NewGame("A", "B", 1)
	// left and right are on flanks; only mid is a legal (non-flank) target.
	left := g.AddToBattleline(testCreature("left", 10), 1)
	mid := g.AddToBattleline(testCreature("mid", 10), 1)
	right := g.AddToBattleline(testCreature("right", 10), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := SplashDamage{Amount: 4, Splash: 2}
	if e.Text() != "deal 4 damage to a creature that is not on a flank and 2 damage to each of its neighbors" {
		t.Errorf("text = %q", e.Text())
	}
	// The default chooser picks the only non-flank creature, mid.
	e.Resolve(ctx)
	if g.Damage(mid) != 4 {
		t.Errorf("chosen damage = %d, want 4", g.Damage(mid))
	}
	if g.Damage(left) != 2 || g.Damage(right) != 2 {
		t.Errorf("neighbor damage = %d/%d, want 2/2", g.Damage(left), g.Damage(right))
	}

	// No non-flank creature (two creatures, both flanks): nothing happens.
	g2 := NewGame("A", "B", 1)
	a := g2.AddToBattleline(testCreature("a", 5), 1)
	g2.AddToBattleline(testCreature("b", 5), 1)
	SplashDamage{Amount: 4, Splash: 2}.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Damage(a) != 0 {
		t.Errorf("no legal target should deal no damage, got %d", g2.Damage(a))
	}
}
