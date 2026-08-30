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

func TestLoseAemberAllBut(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 8 // over the cap: loses 3, left with 5
	g.State.Aember[1] = 3 // at or below the cap: loses 0
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := LoseAember{Player: EachPlayer, By: AllBut(5)}
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

func TestLoseAemberHalfAndEachPlayer(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 5 // loses 2 (floor 5/2), keeps 3
	g.State.Aember[1] = 4 // loses 2, keeps 2
	ctx := &EffectContext{Resolver: g, Controller: 0}

	each := LoseAember{Player: EachPlayer, By: Half}
	if each.Text() != "each player loses half of their Æmber, rounded down" {
		t.Errorf("each text = %q", each.Text())
	}
	each.Resolve(ctx)
	if g.State.Aember[0] != 3 || g.State.Aember[1] != 2 {
		t.Errorf("aember after each-half = %d/%d, want 3/2", g.State.Aember[0], g.State.Aember[1])
	}

	// A controller losing half uses the "your" possessive.
	me := LoseAember{Player: Controller, By: Half}
	if me.Text() != "lose half of your Æmber, rounded down" {
		t.Errorf("controller-half text = %q", me.Text())
	}
	me.Resolve(ctx) // player 0 has 3 -> loses 1 -> 2
	if g.State.Aember[0] != 2 {
		t.Errorf("controller-half aember = %d, want 2", g.State.Aember[0])
	}
}

func TestLoseAemberValidate(t *testing.T) {
	if err := (LoseAember{Player: Controller, Amount: 2, By: Half}).validate(); err == nil {
		t.Error("setting both Amount and By should be rejected")
	}
	if err := (LoseAember{Player: Controller, By: Half}).validate(); err != nil {
		t.Errorf("By alone should be valid, got %v", err)
	}
	if err := (LoseAember{Player: Controller, Amount: 2}).validate(); err != nil {
		t.Errorf("Amount alone should be valid, got %v", err)
	}
}
