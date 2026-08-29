package engine

import "testing"

func TestGainAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	e := GainAember{Player: Controller, Amount: 2}
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

func TestGainAemberPerCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Keys[1] = 2 // opponent has forged 2 keys
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := GainAember{Player: Controller, Amount: 1, Per: OpponentForgedKeys{}}
	if e.Text() != "for each key your opponent has forged, gain 1 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Aember(0) != 2 { // 1 per each of the 2 forged keys
		t.Errorf("aember = %d, want 2", g.Aember(0))
	}
}

func TestGainAemberPerArchivedCards(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < 3; i++ {
		g.State.Archives[0].add(g.Register(testCreature("a", 1), 0))
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := GainAember{Player: Controller, Amount: 1, Per: CardsInArchives{Player: Controller}}
	if e.Text() != "for each card in your archives, gain 1 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Aember(0) != 3 { // 1 per each of the 3 archived cards
		t.Errorf("aember = %d, want 3", g.Aember(0))
	}
}

func TestLoseAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[0] = 3
	g.State.Aember[1] = 1

	self := LoseAember{Player: Controller, Amount: 2}
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

func TestCaptureAllAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("drumble", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 7

	e := CaptureAember{All: true}
	if e.Text() != "{self} captures all your opponent's Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.AmberOn(src) != 7 {
		t.Errorf("captured = %d, want 7 (all of the opponent's pool)", g.AmberOn(src))
	}
	if g.Aember(1) != 0 {
		t.Errorf("opponent aember = %d, want 0", g.Aember(1))
	}
}

func TestExaltEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if got := (Exalt{Target: Target{Kind: TargetChosenFriendlyCreature}, Times: 1}).Text(); got != "exalt a friendly creature" {
		t.Errorf("single exalt text = %q", got)
	}
	e := Exalt{Target: Target{Kind: TargetChosenEnemyCreature}, Times: 2}
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

func TestEachPlayerLosesAllBut(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 8 // over the cap
	g.State.Aember[1] = 3 // at or below the cap
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := EachPlayerLosesAllBut{Keep: 5}
	if e.Text() != "each player with 6 Æmber or more loses all but 5 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 5 {
		t.Errorf("player 0 aember = %d, want 5 (reduced)", g.State.Aember[0])
	}
	if g.State.Aember[1] != 3 {
		t.Errorf("player 1 aember = %d, want 3 (unchanged)", g.State.Aember[1])
	}
}

func TestEachPlayerLosesHalfAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 5 // loses 2 (floor 5/2), keeps 3
	g.State.Aember[1] = 0 // loses nothing (skip branch)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := EachPlayerLosesHalfAember{}
	if e.Text() != "each player loses half of their Æmber, rounded down" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 3 {
		t.Errorf("player 0 aember = %d, want 3", g.State.Aember[0])
	}
	if g.State.Aember[1] != 0 {
		t.Errorf("player 1 aember = %d, want 0", g.State.Aember[1])
	}
}
