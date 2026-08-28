package engine

import "testing"

func TestDestroyEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("shaker", 7), 0)
	weakFriendly := g.AddToBattleline(testCreature("weak", 3), 0)
	strongEnemy := g.AddToBattleline(testCreature("strong", 5), 1)
	weakEnemy := g.AddToBattleline(testCreature("weakfoe", 2), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	byPower := Destroy{Target: Target{Kind: TargetEachCreature}.PowerAtMost(3)}
	if byPower.Text() != "destroy each creature with power 3 or lower" {
		t.Errorf("power text = %q", byPower.Text())
	}
	byPower.Resolve(ctx)
	if g.inPlay(weakFriendly) || g.inPlay(weakEnemy) {
		t.Error("power<=3 creatures should be destroyed")
	}
	if !g.inPlay(src) || !g.inPlay(strongEnemy) {
		t.Error("power>3 creatures should survive")
	}

	sci := g.AddToBattleline(NewCard("sci", Logos, Creature, Common, WithPower(6), WithTraits("Scientist")), 1)
	byTrait := Destroy{Target: Target{Kind: TargetEachCreature}.WithTrait("Scientist")}
	if byTrait.Text() != "destroy each Scientist trait creature" {
		t.Errorf("trait text = %q", byTrait.Text())
	}
	byTrait.Resolve(ctx)
	if g.inPlay(sci) {
		t.Error("Scientist creature should be destroyed")
	}
	if !g.inPlay(strongEnemy) {
		t.Error("non-Scientist creature should survive")
	}
}

func TestDestroyChosenArtifact(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddArtifact(exAutocannon(), 0)
	theirs := g.AddArtifact(exAutocannon(), 1)
	ctx := &EffectContext{Resolver: g, Source: mine, Controller: 0}

	e := Destroy{Target: Target{Kind: TargetChosenArtifact}}
	if e.Text() != "destroy an artifact" {
		t.Errorf("text = %q", e.Text())
	}
	// The default chooser picks the first candidate (the controller's artifact).
	e.Resolve(ctx)
	if g.inPlay(mine) {
		t.Error("the chosen artifact should be destroyed and removed from play")
	}
	if !g.inPlay(theirs) {
		t.Error("the other artifact should be untouched")
	}
}
