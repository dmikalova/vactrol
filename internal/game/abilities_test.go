package game

import "testing"

// These tests exercise creature abilities and combat keywords at the game level,
// using the generic example blueprints from helpers_test.go. Tests for specific
// released cards live with those cards in the set packages under cards/.

func TestReapAndActionAbility(t *testing.T) {
	g := started(t)
	reaper := testCreature("reaper", 2, WithAbility(TriggerReap, GainAember{Amount: 1}))
	id := g.AddToBattleline(reaper, 0)
	if err := g.Reap(0, id); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	// +1 from reap, +1 from the reap ability.
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2", g.Aember(0))
	}
	if !g.Exhausted(id) {
		t.Error("reaper should be exhausted")
	}

	// Action ability.
	actor := testCreature("actor", 2, WithAbility(TriggerAction, GainAember{Amount: 3}))
	aid := g.AddToBattleline(actor, 0)
	if err := g.UseAction(0, aid); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Aember(0) != 5 {
		t.Errorf("aember = %d, want 5", g.Aember(0))
	}
	// Re-using an exhausted creature fails via canUse.
	if err := g.UseAction(0, aid); err != ErrCardExhausted {
		t.Errorf("second UseAction err = %v, want ErrCardExhausted", err)
	}
	// A creature without an action ability cannot use one.
	noAction := g.AddToBattleline(testCreature("plain", 2), 0)
	if err := g.UseAction(0, noAction); err != ErrWrongType {
		t.Errorf("UseAction err = %v, want ErrWrongType", err)
	}
}

func TestCanUseErrors(t *testing.T) {
	g := started(t)
	// Enemy creature (owner mismatch).
	enemy := g.AddToBattleline(testCreature("enemy", 2), 1)
	if err := g.Reap(0, enemy); err != ErrWrongType {
		t.Errorf("reap enemy err = %v, want ErrWrongType", err)
	}
	// Exhausted.
	own := g.AddToBattleline(testCreature("own", 2), 0)
	g.State.Cards[own].Exhausted = true
	if err := g.CanUse(0, own); err != ErrCardExhausted {
		t.Errorf("CanUse exhausted = %v, want ErrCardExhausted", err)
	}
	if err := g.Reap(0, own); err != ErrCardExhausted {
		t.Errorf("err = %v, want ErrCardExhausted", err)
	}
	g.State.Cards[own].Exhausted = false
	// A creature played this turn enters exhausted, so it cannot be used yet.
	sick := g.AddToHand(testCreature("sick", 2), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, sick), false); err != nil {
		t.Fatal(err)
	}
	if !g.Exhausted(sick) {
		t.Error("a freshly played creature should be exhausted")
	}
	if err := g.Reap(0, sick); err != ErrCardExhausted {
		t.Errorf("err = %v, want ErrCardExhausted", err)
	}
	// Wrong house.
	disC := NewCard("Dis Creature", Dis, Creature, Common, WithPower(2))
	disID := g.AddToBattleline(disC, 0)
	if err := g.Reap(0, disID); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// Not active player.
	if err := g.Reap(1, own); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Game over.
	g.State.Winner = 0
	if err := g.Reap(0, own); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

func TestFight(t *testing.T) {
	g := started(t)
	attacker := g.AddToBattleline(testCreature("att", 5), 0)
	defender := g.AddToBattleline(testCreature("def", 4), 1)
	if err := g.Fight(0, attacker, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("defender should be destroyed")
	}
	if g.Damage(attacker) != 4 {
		t.Errorf("attacker damage = %d, want 4", g.Damage(attacker))
	}

	// Invalid target: a friendly creature.
	friend := g.AddToBattleline(testCreature("friend", 2), 0)
	att2 := g.AddToBattleline(testCreature("att2", 2), 0)
	if err := g.Fight(0, att2, friend); err != ErrNoTarget {
		t.Errorf("err = %v, want ErrNoTarget", err)
	}
	// Invalid attacker: exhausted (fails via canUse before target checks).
	g.State.Cards[att2].Exhausted = true
	if err := g.Fight(0, att2, friend); err != ErrCardExhausted {
		t.Errorf("err = %v, want ErrCardExhausted", err)
	}
}

func TestSkirmishAndArmorAndPoison(t *testing.T) {
	g := started(t)
	// Skirmish attacker takes no retaliation.
	skirm := g.AddToBattleline(testCreature("skirm", 3, WithKeywords(Skirmish)), 0)
	def := g.AddToBattleline(testCreature("def", 2), 1)
	if err := g.Fight(0, skirm, def); err != nil {
		t.Fatal(err)
	}
	if g.Damage(skirm) != 0 {
		t.Errorf("skirmisher took %d damage, want 0", g.Damage(skirm))
	}

	// Armor absorbs damage.
	armored := g.AddToBattleline(testCreature("armored", 5, WithArmor(2)), 0)
	g.applyDamage(armored, 3) // 2 absorbed, 1 through
	if g.Damage(armored) != 1 {
		t.Errorf("armored damage = %d, want 1", g.Damage(armored))
	}
	g.applyDamage(armored, 0) // no-op path

	// Poison: any damage destroys.
	poison := g.AddToBattleline(testCreature("poison", 5, WithKeywords(Poison)), 1)
	g.applyDamage(poison, 1)
	g.checkDestroyed(poison)
	if g.inPlay(poison) {
		t.Error("poisoned creature should be destroyed")
	}
	// checkDestroyed on a non-in-play card is a no-op.
	g.checkDestroyed(poison)
}

func TestDestroyMovesUpgradesToDiscard(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 1), 0)
	g.AddToHand(exBruteStrength(), 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil {
		t.Fatal(err)
	}
	// Host now has 6 power (1 + 5). Kill it with enough damage.
	g.applyDamage(host, 6)
	g.checkDestroyed(host)
	if g.inPlay(host) {
		t.Error("host should be destroyed")
	}
	if len(g.Discard(0)) != 2 { // host + upgrade
		t.Errorf("discard size = %d, want 2", len(g.Discard(0)))
	}
}

func TestDestroyGivesAmberToOpponent(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0) // owned by player 0
	g.State.Cards[c].Amber = 2
	if g.AmberOn(c) != 2 {
		t.Fatalf("AmberOn = %d, want 2", g.AmberOn(c))
	}
	g.destroy(c)
	// The Æmber on the creature goes to its owner's opponent (player 1).
	if g.Aember(1) != 2 {
		t.Errorf("opponent aember = %d, want 2", g.Aember(1))
	}
	if g.AmberOn(c) != 0 {
		t.Errorf("amber not cleared: %d", g.AmberOn(c))
	}
}
