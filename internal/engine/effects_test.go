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
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

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

func TestReturnToHandEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.State.Cards[src].Damage = 2
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := ReturnToHand{Target: Target{Kind: TargetThisCreature}}
	if e.Text() != "put this creature into its owner's hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Battleline[0].Count != 0 {
		t.Errorf("battleline not cleared: %d", g.State.Battleline[0].Count)
	}
	if g.State.Hand[0].Count != 1 || g.State.Hand[0].IDs[0] != src {
		t.Errorf("hand = count %d id %v, want 1 / %d", g.State.Hand[0].Count, g.State.Hand[0].IDs[0], src)
	}
	// Returning clears the per-match state the card accrued in play.
	if g.State.Cards[src].Damage != 0 {
		t.Errorf("damage after return = %d, want 0", g.State.Cards[src].Damage)
	}
}

func TestGainAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	e := GainAember{Amount: 2}
	if e.Text() != "gain 2 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 2 {
		t.Errorf("aember = %d, want 2", g.State.Aember[0])
	}

	foe := GainAember{Player: Opponent, Amount: 1}
	if foe.Text() != "your opponent gains 1 Æmber" {
		t.Errorf("enemy text = %q", foe.Text())
	}
	foe.Resolve(ctx)
	if g.State.Aember[1] != 1 {
		t.Errorf("opponent aember = %d, want 1", g.State.Aember[1])
	}
}

func TestLoseAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[0] = 3
	g.State.Aember[1] = 1

	self := LoseAember{Amount: 2}
	if self.Text() != "lose 2 Æmber" {
		t.Errorf("self text = %q", self.Text())
	}
	self.Resolve(ctx)
	if g.State.Aember[0] != 1 {
		t.Errorf("self aember = %d, want 1", g.State.Aember[0])
	}

	foe := LoseAember{Player: Opponent, Amount: 4}
	if foe.Text() != "your opponent loses 4 Æmber" {
		t.Errorf("enemy text = %q", foe.Text())
	}
	foe.Resolve(ctx) // opponent has only 1; floors at 0
	if g.State.Aember[1] != 0 {
		t.Errorf("opponent aember = %d, want 0", g.State.Aember[1])
	}
}

func TestStealAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 1

	e := StealAember{Amount: 3}
	if e.Text() != "steal 3 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx) // opponent has only 1
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 0 {
		t.Errorf("after steal: you=%d opp=%d, want 1/0", g.State.Aember[0], g.State.Aember[1])
	}
}

func TestCaptureAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 2

	e := CaptureAember{Amount: 3}
	if e.Text() != "{self} captures 3 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx) // opponent has only 2
	if g.State.Cards[src].Amber != 2 {
		t.Errorf("captured = %d, want 2", g.State.Cards[src].Amber)
	}
	if g.State.Aember[1] != 0 {
		t.Errorf("opponent aember = %d, want 0", g.State.Aember[1])
	}
}

func TestDrawEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	g.AddToDeck(testCreature("d1", 1), 0)
	g.AddToDeck(testCreature("d2", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if (Draw{Amount: 1}).Text() != "draw a card" {
		t.Errorf("single text = %q", (Draw{Amount: 1}).Text())
	}
	e := Draw{Amount: 2}
	if e.Text() != "draw 2 cards" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Hand[0].Count != 2 {
		t.Errorf("hand = %d, want 2", g.State.Hand[0].Count)
	}
	(Draw{Amount: 3}).Resolve(ctx) // deck now empty; draws nothing
	if g.State.Hand[0].Count != 2 {
		t.Errorf("hand after empty-deck draw = %d, want 2", g.State.Hand[0].Count)
	}
}

func TestHealEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 5), 0)
	other := g.AddToBattleline(testCreature("other", 5), 0)
	g.State.Cards[src].Damage = 3
	g.State.Cards[other].Damage = 4
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	partial := Heal{Amount: 2, Target: Target{Kind: TargetThisCreature}}
	if partial.Text() != "heal 2 damage from this creature" {
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

func TestStunEffects(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	stun := Stun{Target: Target{Kind: TargetEachFriendlyCreature}}
	if stun.Text() != "stun each friendly creature" {
		t.Errorf("stun text = %q", stun.Text())
	}
	stun.Resolve(ctx)
	if !g.State.Cards[src].Stunned || !g.State.Cards[friend].Stunned {
		t.Error("stun should stun each friendly creature")
	}
	if g.State.Cards[foe].Stunned {
		t.Error("stun of friendly creatures should not touch the enemy")
	}

	unstun := Unstun{Target: Target{Kind: TargetEachFriendlyCreature}}
	if unstun.Text() != "unstun each friendly creature" {
		t.Errorf("unstun text = %q", unstun.Text())
	}
	unstun.Resolve(ctx)
	if g.State.Cards[src].Stunned || g.State.Cards[friend].Stunned {
		t.Error("unstun should clear the stun on each friendly creature")
	}
}

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

func TestSequenceEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	seq := Sequence{Effects: []Effect{GainAember{Amount: 1}, GainAember{Amount: 2}}}
	if seq.Text() != "gain 1 Æmber, and gain 2 Æmber" {
		t.Errorf("sequence text = %q", seq.Text())
	}
	seq.Resolve(ctx)
	if g.State.Aember[0] != 3 {
		t.Errorf("aember = %d, want 3", g.State.Aember[0])
	}
}

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
func TestOnChosenCreatureEnemyAndNoTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	g.State.Cards[enemy].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	onEnemy := OnChosenCreature{Player: Opponent, Verbs: []CreatureVerb{ReadyVerb{}}}
	if onEnemy.Text() != "ready an enemy creature" {
		t.Errorf("text = %q", onEnemy.Text())
	}
	onEnemy.Resolve(ctx)
	if g.State.Cards[enemy].Exhausted {
		t.Error("enemy should have been readied")
	}

	// No candidates: remove the enemy and resolve again (logs, no panic).
	g.DestroyEach(0, []LocalID{enemy})
	onEnemy.Resolve(ctx)
}

func TestExaltEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if got := (Exalt{Player: Controller, Times: 1}).Text(); got != "exalt a friendly creature" {
		t.Errorf("single exalt text = %q", got)
	}
	e := Exalt{Player: Opponent, Times: 2}
	if e.Text() != "exalt an enemy creature 2 times" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Cards[enemy].Amber != 2 {
		t.Errorf("amber on enemy = %d, want 2", g.State.Cards[enemy].Amber)
	}

	// No candidates: remove the enemy and resolve again (logs, no panic).
	g.DestroyEach(0, []LocalID{enemy})
	e.Resolve(ctx)
}

func TestFightVerbNoEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
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
	if got := g.OrderByChoice(0, "p", nil); len(got) != 0 {
		t.Errorf("nil order = %v", got)
	}
	if got := g.OrderByChoice(0, "p", []LocalID{7}); !eq(got, []LocalID{7}) {
		t.Errorf("single order = %v", got)
	}
	// FirstChooser (default) keeps the original order.
	if got := g.OrderByChoice(0, "p", ids); !eq(got, []LocalID{10, 20, 30}) {
		t.Errorf("first-chooser order = %v", got)
	}
	// Picking the last each time reverses the order.
	g.SetChooser(0, orderLastChooser{})
	if got := g.OrderByChoice(0, "p", ids); !eq(got, []LocalID{30, 20, 10}) {
		t.Errorf("last-chooser order = %v", got)
	}
	// A rejected pick falls back to the given order.
	g.SetChooser(0, orderRejectChooser{})
	if got := g.OrderByChoice(0, "p", ids); !eq(got, []LocalID{10, 20, 30}) {
		t.Errorf("reject order = %v", got)
	}
}
