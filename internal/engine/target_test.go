package engine

import "testing"

func TestTargetSelect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	myArt := g.AddArtifact(NewCard("myrelic", Brobnar, Artifact, Rare), 0)
	enemyArt := g.AddArtifact(NewCard("enemyrelic", Brobnar, Artifact, Rare), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if ids := (Target{Kind: TargetThisCreature}).Select(ctx); len(ids) != 1 || ids[0] != src {
		t.Errorf("this-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetTriggeringCreature}).Select(ctx); ids != nil {
		t.Errorf("triggering-creature without It should be nil, got %v", ids)
	}
	ctx.HasIt, ctx.It = true, enemy
	if ids := (Target{Kind: TargetTriggeringCreature}).Select(ctx); len(ids) != 1 || ids[0] != enemy {
		t.Errorf("triggering-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachEnemyCreature}).Select(ctx); len(ids) != 1 || ids[0] != enemy {
		t.Errorf("each-enemy select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachCreature}).Select(ctx); len(ids) != 2 || ids[0] != src || ids[1] != enemy {
		t.Errorf("each-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachArtifact}).Select(ctx); len(ids) != 2 || ids[0] != myArt || ids[1] != enemyArt {
		t.Errorf("each-artifact select = %v", ids)
	}
	if ids := (Target{Kind: TargetKind(99)}).Select(ctx); ids != nil {
		t.Errorf("default select should be nil, got %v", ids)
	}

	// TargetChosenCreature asks the chooser to pick one creature from either side.
	if ids := (Target{Kind: TargetChosenCreature}).Select(ctx); len(ids) != 1 || ids[0] != src {
		t.Errorf("chosen-creature (first chooser) = %v, want [%d]", ids, src)
	}
	g.SetChooser(0, orderRejectChooser{})
	if ids := (Target{Kind: TargetChosenCreature}).Select(ctx); ids != nil {
		t.Errorf("chosen-creature (reject) = %v, want nil", ids)
	}
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := (Target{Kind: TargetChosenCreature}).Select(empty); ids != nil {
		t.Errorf("chosen-creature (no candidates) = %v, want nil", ids)
	}

	// TargetChosenEnemyCreature only offers enemy creatures.
	g.SetChooser(0, nil) // back to the default FirstChooser
	if ids := (Target{Kind: TargetChosenEnemyCreature}).Select(ctx); len(ids) != 1 || ids[0] != enemy {
		t.Errorf("chosen-enemy = %v, want [%d]", ids, enemy)
	}
	g.SetChooser(0, orderRejectChooser{})
	if ids := (Target{Kind: TargetChosenEnemyCreature}).Select(ctx); ids != nil {
		t.Errorf("chosen-enemy (reject) = %v, want nil", ids)
	}
	g.SetChooser(0, nil)

	// Damaged filter keeps only creatures with damage on them.
	g.State.Cards[src].Damage = 1
	if ids := (Target{Kind: TargetEachCreature}).Damaged().Select(ctx); len(ids) != 1 || ids[0] != src {
		t.Errorf("damaged filter = %v, want [%d]", ids, src)
	}
	g.State.Cards[src].Damage = 0

	// OnFlank keeps only the leftmost/rightmost creatures of a battleline.
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	right := g.AddToBattleline(testCreature("right", 1), 0)
	// Player 0's battleline is now [src, mid, right]; only src and right are flanks.
	if ids := (Target{Kind: TargetEachFriendlyCreature}).OnFlank().Select(ctx); len(ids) != 2 || ids[0] != src || ids[1] != right {
		t.Errorf("flank filter = %v, want [%d %d]", ids, src, right)
	}
	// NotOnFlank keeps only the interior creatures (here, just mid).
	if ids := (Target{Kind: TargetEachFriendlyCreature}).NotOnFlank().Select(ctx); len(ids) != 1 || ids[0] != mid {
		t.Errorf("not-on-flank filter = %v, want [%d]", ids, mid)
	}
}

func TestTargetText(t *testing.T) {
	cases := map[TargetKind]string{
		TargetThisCreature:         "this creature",
		TargetTriggeringCreature:   "it",
		TargetEachCreature:         "each creature",
		TargetEachFriendlyCreature: "each friendly creature",
		TargetEachEnemyCreature:    "each enemy creature",
		TargetEachArtifact:         "each artifact",
		TargetKind(99):             "a creature",
	}
	for kind, want := range cases {
		if got := (Target{Kind: kind}).Text(); got != want {
			t.Errorf("Text(%d) = %q, want %q", kind, got, want)
		}
	}
	if got := (Target{Kind: TargetChosenEnemyCreature}).Text(); got != "an enemy creature" {
		t.Errorf("chosen-enemy text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).Damaged().Text(); got != "a damaged creature" {
		t.Errorf("damaged text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).OnFlank().Text(); got != "a flank creature" {
		t.Errorf("flank text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).NotOnFlank().Text(); got != "a creature that is not on a flank" {
		t.Errorf("not-on-flank text = %q", got)
	}
	if got := (Target{Kind: TargetEachEnemyCreature}).NotOnFlank().Text(); got != "each enemy creature that is not on a flank" {
		t.Errorf("not-on-flank each text = %q", got)
	}
}
