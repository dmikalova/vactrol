package game

import "testing"

// ---- playing cards ----

func TestPlayCreatureFlankLeft(t *testing.T) {
	g := started(t)
	first := g.AddToHand(testCreature("first", 2), 0)
	second := g.AddToHand(testCreature("second", 2), 0)
	if _, err := g.PlayCreature(0, 0, false); err != nil { // right flank: [first]
		t.Fatal(err)
	}
	if _, err := g.PlayCreature(0, 0, true); err != nil { // left flank: [second, first]
		t.Fatal(err)
	}
	bl := g.Battleline(0)
	if bl[0] != second || bl[1] != first {
		t.Errorf("battleline order = %v, want [%d %d]", bl, second, first)
	}
}

func TestPlayTypeMismatchErrors(t *testing.T) {
	g := started(t)
	creatureID := g.AddToHand(testCreature("c", 2), 0)
	_ = creatureID
	// Playing a creature as an artifact/action/upgrade should fail on type.
	if _, err := g.PlayArtifact(0, 0); err != ErrWrongType {
		t.Errorf("PlayArtifact err = %v, want ErrWrongType", err)
	}
	if err := g.PlayAction(0, 0); err != ErrWrongType {
		t.Errorf("PlayAction err = %v, want ErrWrongType", err)
	}
	if _, err := g.PlayUpgrade(0, 0); err != ErrWrongType {
		t.Errorf("PlayUpgrade err = %v, want ErrWrongType", err)
	}
	// Now play it correctly.
	if _, err := g.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
}

func TestTakeFromHandErrors(t *testing.T) {
	g := started(t)
	// Not active player.
	if _, err := g.PlayCreature(1, 0, false); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Bad hand index.
	if _, err := g.PlayCreature(0, 5, false); err != ErrCardNotInHand {
		t.Errorf("err = %v, want ErrCardNotInHand", err)
	}
	// Wrong house.
	dis := NewCard("Dissonant", Dis, Creature, Common, WithPower(2))
	g.AddToHand(dis, 0)
	if _, err := g.PlayCreature(0, 0, false); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// Game over.
	g.State.Winner = 0
	if _, err := g.PlayCreature(0, 0, false); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

func TestPlayArtifactAndAction(t *testing.T) {
	g := started(t)
	g.AddToHand(exAutocannon(), 0)
	if _, err := g.PlayArtifact(0, 0); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}
	if g.Aember(0) != 1 {
		t.Errorf("aember after artifact = %d, want 1", g.Aember(0))
	}

	// Action with a friendly creature to ready-and-fight.
	g.AddToBattleline(testCreature("friend", 4), 0)
	g.AddToBattleline(testCreature("foe", 2), 1)
	g.AddToHand(exBattleFury(), 0)
	if err := g.PlayAction(0, handIdx(g, 0, "Battle Fury")); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("foe should be destroyed by the ready-and-fight")
	}
	if len(g.Discard(0)) == 0 {
		t.Error("action card should be in discard")
	}
}

func TestPlayUpgradeErrorsAndSuccess(t *testing.T) {
	g := started(t)
	// No creature to attach to.
	g.AddToHand(exBruteStrength(), 0)
	if _, err := g.PlayUpgrade(0, 0); err != ErrNoTarget {
		t.Errorf("err = %v, want ErrNoTarget", err)
	}
	// Not active player.
	if _, err := g.PlayUpgrade(1, 0); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Bad index.
	if _, err := g.PlayUpgrade(0, 9); err != ErrCardNotInHand {
		t.Errorf("err = %v, want ErrCardNotInHand", err)
	}
	// Wrong house.
	disUp := NewCard("Dis Upgrade", Dis, Upgrade, Common, WithStatic(StaticModifier{PowerBonus: 1}))
	g.AddToHand(disUp, 0)
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Dis Upgrade")); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// Success: attach +5 power upgrade to a creature.
	host := g.AddToBattleline(testCreature("host", 3), 0)
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Brute Strength")); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}
	if g.Power(host) != 8 {
		t.Errorf("host power = %d, want 8", g.Power(host))
	}
	if ups := g.Upgrades(host); len(ups) != 1 {
		t.Errorf("Upgrades(host) = %v, want 1 upgrade", ups)
	}
	// Game over branch.
	g.State.Winner = 0
	if _, err := g.PlayUpgrade(0, 0); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

func TestDiscardFromHand(t *testing.T) {
	g := started(t)
	// Wrong house.
	dis := NewCard("Dis", Dis, Creature, Common, WithPower(2))
	g.AddToHand(dis, 0)
	if err := g.DiscardFromHand(0, 0); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// Not the active player.
	if err := g.DiscardFromHand(1, 0); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Bad index.
	if err := g.DiscardFromHand(0, 9); err != ErrCardNotInHand {
		t.Errorf("err = %v, want ErrCardNotInHand", err)
	}
	// Success: discard a Brobnar (active-house) card.
	brob := g.AddToHand(testCreature("brob", 3), 0)
	if err := g.DiscardFromHand(0, handIdxByID(g, 0, brob)); err != nil {
		t.Fatalf("DiscardFromHand: %v", err)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != brob {
		t.Errorf("discard = %v, want [%d]", d, brob)
	}
	// Game over.
	g.State.Winner = 0
	if err := g.DiscardFromHand(0, 0); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

// ---- turn flow ----

func TestForgeKeyWins(t *testing.T) {
	g := started(t)
	// Giant on board so the after-forge trigger fires; a tough enemy survives.
	g.AddToBattleline(exGiant(), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 10), 1)
	g.State.Aember[0] = 3 * KeyCost
	g.forgeKeys(0)
	if g.Keys(0) != KeysToWin {
		t.Errorf("keys = %d, want %d", g.Keys(0), KeysToWin)
	}
	if g.Winner() != 0 {
		t.Errorf("winner = %d, want 0", g.Winner())
	}
	// The giant deals 2 damage per forged key: 3 keys => 6 damage.
	if g.Damage(enemy) != 6 {
		t.Errorf("enemy damage = %d, want 6", g.Damage(enemy))
	}
	// BeginTurn after a win is a no-op.
	turn := g.State.Turn
	g.BeginTurn(1)
	if g.State.Turn != turn {
		t.Error("BeginTurn should be a no-op once the game is over")
	}
}

func TestChooseHouseWrongPlayer(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.BeginTurn(0)
	if err := g.ChooseHouse(1, Brobnar); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
}

func TestEndTurnReadyDrawAndArmorRefresh(t *testing.T) {
	g := started(t)
	// A creature with an armor upgrade, exhausted and damaged; plus an artifact.
	host := g.AddToBattleline(testCreature("host", 5, WithArmor(1)), 0)
	g.AddToHand(NewCard("Plating", Brobnar, Upgrade, Common, WithStatic(StaticModifier{ArmorBonus: 2})), 0)
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Plating")); err != nil {
		t.Fatal(err)
	}
	g.State.Cards[host].Exhausted = true
	g.State.Cards[host].ArmorRemaining = 0
	g.AddArtifact(exAutocannon(), 0)

	// Deck to draw from.
	for i := 0; i < 3; i++ {
		g.AddToDeck(testCreature("deckcard", 1), 0)
	}
	g.Shuffle(0)

	g.EndTurn(0)

	if g.Exhausted(host) {
		t.Error("host should be readied")
	}
	if g.State.Cards[host].ArmorRemaining != 3 { // base 1 + upgrade 2
		t.Errorf("armor refreshed to %d, want 3", g.State.Cards[host].ArmorRemaining)
	}
	if len(g.Hand(0)) != 3 {
		t.Errorf("hand size = %d, want 3 drawn", len(g.Hand(0)))
	}
}

// ---- accessors / choosers / logging ----

func TestAccessorsAndChooser(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	if g.PlayerName(0) != "Alice" || g.PlayerName(1) != "Bob" {
		t.Error("player names wrong")
	}
	id := g.AddToBattleline(testCreature("c", 4), 0)
	if g.Def(id).Name != "c" {
		t.Errorf("Def name = %q", g.Def(id).Name)
	}
	if g.Name(id) != "c" {
		t.Errorf("Name = %q", g.Name(id))
	}

	// Artifacts accessor returns the artifact row.
	art := g.AddArtifact(exAutocannon(), 0)
	if arts := g.Artifacts(0); len(arts) != 1 || arts[0] != art {
		t.Errorf("Artifacts = %v, want [%d]", arts, art)
	}

	// Nil chooser falls back to the default (FirstChooser).
	g.SetChooser(0, nil)
	got, ok := g.chooserFor(0).ChooseCreature("p", []LocalID{id})
	if !ok || got != id {
		t.Errorf("default chooser = (%d,%v)", got, ok)
	}
	// Empty candidate list.
	if _, ok := g.chooserFor(0).ChooseCreature("p", nil); ok {
		t.Error("empty candidates should return ok=false")
	}
}

func TestVerboseLogging(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.Verbose = true
	g.logf("hello %d", 1)
	if len(g.Log) != 1 || g.Log[0] != "hello 1" {
		t.Errorf("log = %v", g.Log)
	}
}
