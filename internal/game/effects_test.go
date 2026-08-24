package game

import "testing"

func TestTargetSelect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	myArt := g.AddArtifact(NewCard("myrelic", Brobnar, Artifact, Rare), 0)
	enemyArt := g.AddArtifact(NewCard("enemyrelic", Brobnar, Artifact, Rare), 1)
	ctx := &EffectContext{Game: g, Source: src, Controller: 0}

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
}

func TestReturnToDeckEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	myArt := g.AddArtifact(NewCard("myrelic", Brobnar, Artifact, Rare), 0)
	enemyArt := g.AddArtifact(NewCard("enemyrelic", Brobnar, Artifact, Rare), 1)
	g.State.Cards[myArt].Exhausted = true
	ctx := &EffectContext{Game: g, Source: src, Controller: 0}

	e := ReturnToDeck{Target: Target{Kind: TargetEachArtifact}}
	if e.Text() != "put each artifact on top of its owner's deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Artifacts[0].Count != 0 || g.State.Artifacts[1].Count != 0 {
		t.Errorf("artifact rows not cleared: %d %d", g.State.Artifacts[0].Count, g.State.Artifacts[1].Count)
	}
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != myArt {
		t.Errorf("player 0 deck top = %v, want %d", g.State.Deck[0].IDs[0], myArt)
	}
	if g.State.Deck[1].Count != 1 || g.State.Deck[1].IDs[0] != enemyArt {
		t.Errorf("player 1 deck top = %v, want %d", g.State.Deck[1].IDs[0], enemyArt)
	}
	if g.State.Cards[myArt].Exhausted {
		t.Errorf("returned artifact should be readied")
	}
}

func TestGainAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Game: g, Source: src, Controller: 0}
	e := GainAember{Amount: 2}
	if e.Text() != "gain 2 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 2 {
		t.Errorf("aember = %d, want 2", g.State.Aember[0])
	}
}

func TestSequenceEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Game: g, Source: src, Controller: 0}
	seq := Sequence{Effects: []Effect{GainAember{Amount: 1}, GainAember{Amount: 2}}}
	if seq.Text() != "gain 1 Æmber, and gain 2 Æmber" {
		t.Errorf("sequence text = %q", seq.Text())
	}
	seq.Resolve(ctx)
	if g.State.Aember[0] != 3 {
		t.Errorf("aember = %d, want 3", g.State.Aember[0])
	}
}

func TestOnChosenCreatureEnemyAndNoTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	g.State.Cards[enemy].Exhausted = true
	ctx := &EffectContext{Game: g, Source: src, Controller: 0}

	onEnemy := OnChosenCreature{Controller: Enemy, Verbs: []CreatureVerb{ReadyVerb{}}}
	if onEnemy.Text() != "ready an enemy creature" {
		t.Errorf("text = %q", onEnemy.Text())
	}
	onEnemy.Resolve(ctx)
	if g.State.Cards[enemy].Exhausted {
		t.Error("enemy should have been readied")
	}

	// No candidates: remove the enemy and resolve again (logs, no panic).
	g.State.Battleline[1].remove(enemy)
	onEnemy.Resolve(ctx)
}

func TestExaltEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	ctx := &EffectContext{Game: g, Source: src, Controller: 0}

	if got := (Exalt{Controller: Friendly, Times: 1}).Text(); got != "exalt a friendly creature" {
		t.Errorf("single exalt text = %q", got)
	}
	e := Exalt{Controller: Enemy, Times: 2}
	if e.Text() != "exalt an enemy creature 2 times" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Cards[enemy].Amber != 2 {
		t.Errorf("amber on enemy = %d, want 2", g.State.Cards[enemy].Amber)
	}

	// No candidates: remove the enemy and resolve again (logs, no panic).
	g.State.Battleline[1].remove(enemy)
	e.Resolve(ctx)
}

func TestFightVerbNoEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 2), 0)
	ctx := &EffectContext{Game: g, Source: src, Controller: 0}
	FightVerb{}.Apply(ctx, src) // no enemies -> logs and returns
	if g.State.Cards[src].Exhausted {
		t.Error("no fight should have occurred")
	}
}

// orderLastChooser always picks the last candidate, reversing an ordering.
type orderLastChooser struct{}

func (orderLastChooser) ChooseCreature(_ string, c []LocalID) (LocalID, bool) {
	return c[len(c)-1], true
}

// orderRejectChooser refuses to pick, so ordering falls back to the given order.
type orderRejectChooser struct{}

func (orderRejectChooser) ChooseCreature(string, []LocalID) (LocalID, bool) { return 0, false }

func TestOrderByChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	ids := []LocalID{10, 20, 30}
	eq := func(a, b []LocalID) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	// Nothing to order for 0 or 1 ids.
	if got := g.orderByChoice(0, "p", nil); len(got) != 0 {
		t.Errorf("nil order = %v", got)
	}
	if got := g.orderByChoice(0, "p", []LocalID{7}); !eq(got, []LocalID{7}) {
		t.Errorf("single order = %v", got)
	}
	// FirstChooser (default) keeps the original order.
	if got := g.orderByChoice(0, "p", ids); !eq(got, []LocalID{10, 20, 30}) {
		t.Errorf("first-chooser order = %v", got)
	}
	// Picking the last each time reverses the order.
	g.SetChooser(0, orderLastChooser{})
	if got := g.orderByChoice(0, "p", ids); !eq(got, []LocalID{30, 20, 10}) {
		t.Errorf("last-chooser order = %v", got)
	}
	// A rejected pick falls back to the given order.
	g.SetChooser(0, orderRejectChooser{})
	if got := g.orderByChoice(0, "p", ids); !eq(got, []LocalID{10, 20, 30}) {
		t.Errorf("reject order = %v", got)
	}
}
