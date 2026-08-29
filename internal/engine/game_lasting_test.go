package engine

import "testing"

func TestFireLastingOrders(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	// Two reactions fire on the same event; the controller orders them.
	g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)
	g.AddLasting(EventCreaturePlayed, actDealDamage, 0, 2)
	g.SetChooser(0, optionPicker{idx: 1}) // resolve "deal damage" first

	g.AddToHand(testCreature("minion", 4), 0)
	before := g.Aember(0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "minion"), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got := g.Aember(0) - before; got != 1 {
		t.Errorf("Æmber gained = %d, want 1", got)
	}
	if g.Damage(foe) != 2 {
		t.Errorf("foe damage = %d, want 2", g.Damage(foe))
	}
}

func TestClearLastingKeepsOtherPlayer(t *testing.T) {
	g := started(t)
	g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)
	g.AddLasting(EventCreaturePlayed, actGainAember, 1, 1) // the opponent's reaction

	g.clearLasting(0)

	if g.State.LastingCount != 1 {
		t.Fatalf("lasting count = %d, want 1 (the opponent's kept)", g.State.LastingCount)
	}
	if g.State.Lasting[0].Controller != 1 {
		t.Errorf("kept controller = %d, want 1", g.State.Lasting[0].Controller)
	}
}

func TestAddLastingCap(t *testing.T) {
	g := started(t)
	for i := 0; i < maxLasting+3; i++ {
		g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)
	}
	if int(g.State.LastingCount) != maxLasting {
		t.Errorf("lasting count = %d, want %d (capped)", g.State.LastingCount, maxLasting)
	}
}
