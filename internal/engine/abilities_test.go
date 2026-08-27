package engine

import "testing"

// These tests exercise creature abilities and combat keywords at the game level,
// using the generic example blueprints from helpers_test.go. Tests for specific
// released cards live with those cards in the set packages under cards/.

func TestReapAndActionAbility(t *testing.T) {
	g := started(t)
	reaper := testCreature("reaper", 2, WithAbility(TriggerAfterReap, GainAember{Amount: 1}))
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
	// An artifact is in play but is not a creature, so it cannot reap or fight.
	relic := g.AddArtifact(NewCard("Relic", Brobnar, Artifact, Common), 0)
	if err := g.Reap(0, relic); err != ErrWrongType {
		t.Errorf("reap artifact err = %v, want ErrWrongType", err)
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

func TestFightSimultaneousDestruction(t *testing.T) {
	g := started(t)
	// Equal-power creatures with no armor deal lethal damage to each other; both
	// are judged destroyed from the same combat and removed together.
	att := g.AddToBattleline(testCreature("att", 3), 0)
	def := g.AddToBattleline(testCreature("def", 3), 1)
	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if len(g.Battleline(0)) != 0 || len(g.Battleline(1)) != 0 {
		t.Errorf("both fighters should be destroyed: friendly=%v enemy=%v", g.Battleline(0), g.Battleline(1))
	}
}

// recordEnemyCount is a test effect that records how many enemy creatures are in
// play at the moment it resolves.
type recordEnemyCount struct{ got *int }

func (recordEnemyCount) Text() string { return "record enemy count" }
func (e recordEnemyCount) Resolve(ctx *EffectContext) {
	*e.got = len(ctx.Resolver.Battleline(1 - ctx.Controller))
}

func TestDestroyTogetherResolvesBeforeDiscard(t *testing.T) {
	g := started(t)
	// A's "Destroyed:" ability records how many enemy creatures are in play when
	// it fires. KeyForge tags both for destruction and resolves the Destroyed
	// abilities before moving anything to the discard, so B is still present.
	var enemiesWhenADied int
	a := g.AddToBattleline(testCreature("a", 3, WithAbility(TriggerDestroyed, recordEnemyCount{got: &enemiesWhenADied})), 0)
	b := g.AddToBattleline(testCreature("b", 3), 1)

	g.DestroyEach(0, []LocalID{a, b})

	if enemiesWhenADied != 1 {
		t.Errorf("A's Destroyed ability saw %d enemy creatures; want 1 (still in play until discard)", enemiesWhenADied)
	}
	if g.inPlay(a) || g.inPlay(b) {
		t.Error("both creatures should be in the discard after the event")
	}
}

func TestDestroyedRelocationSkipsDiscard(t *testing.T) {
	g := started(t)
	// A creature whose "Destroyed:" ability returns it to the top of its deck
	// leaves play during the event, so it is not also moved to the discard.
	c := g.AddToBattleline(testCreature("wanderer", 3, WithAbility(TriggerDestroyed, ReturnToDeck{Target: Target{Kind: TargetThisCreature}})), 0)
	g.DestroyEach(0, []LocalID{c})

	if g.inPlay(c) {
		t.Error("creature should have left play")
	}
	if len(g.Discard(0)) != 0 {
		t.Errorf("relocated creature should not be discarded; discard size = %d", len(g.Discard(0)))
	}
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != c {
		t.Error("creature should be on top of its owner's deck")
	}
}

func TestBeforeFightTrigger(t *testing.T) {
	g := started(t)
	spit := DealDamage{Amount: 1, Target: Target{Kind: TargetEachEnemyCreature}}

	// Defender survives: before-fight damages every enemy, then combat proceeds.
	att := g.AddToBattleline(testCreature("spitter", 3, WithAbility(TriggerBeforeFight, spit)), 0)
	weak := g.AddToBattleline(testCreature("weak", 1), 1)
	tough := g.AddToBattleline(testCreature("tough", 10), 1)
	if err := g.Fight(0, att, tough); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.inPlay(weak) {
		t.Error("before-fight damage should destroy the 1-power enemy")
	}
	if g.Damage(tough) != 4 { // 1 before-fight + 3 combat
		t.Errorf("tough damage = %d, want 4", g.Damage(tough))
	}
	if g.inPlay(att) {
		t.Error("attacker should die to the 10-power defender's retaliation")
	}

	// Defender destroyed before combat: the attack still happens but deals no
	// combat damage, so the attacker takes no retaliation.
	att2 := g.AddToBattleline(testCreature("spitter2", 3, WithAbility(TriggerBeforeFight, spit)), 0)
	frail := g.AddToBattleline(testCreature("frail", 1), 1)
	if err := g.Fight(0, att2, frail); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.inPlay(frail) {
		t.Error("frail defender should be destroyed before combat")
	}
	if g.Damage(att2) != 0 {
		t.Errorf("attacker took %d damage, want 0 (no combat)", g.Damage(att2))
	}
	if !g.State.Cards[att2].Exhausted {
		t.Error("attacker should be exhausted from attacking")
	}
}

func TestAfterDestroyedFightingTrigger(t *testing.T) {
	gain := GainAember{Amount: 1}

	// The attacker survives and destroys the defender: its ability fires.
	g := started(t)
	hunter := g.AddToBattleline(testCreature("hunter", 6, WithAbility(TriggerAfterDestroyedFighting, gain)), 0)
	prey := g.AddToBattleline(testCreature("prey", 1), 1)
	before := g.Aember(0)
	if err := g.Fight(0, hunter, prey); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.inPlay(prey) {
		t.Fatal("prey should be destroyed by the fight")
	}
	if g.Aember(0) != before+1 {
		t.Errorf("survivor aember = %d, want %d", g.Aember(0), before+1)
	}

	// The defender survives while its Hazardous destroys the attacker: the
	// defender's ability fires (for its own controller).
	g2 := started(t)
	weak := g2.AddToBattleline(testCreature("weak", 1), 0)
	guard := g2.AddToBattleline(NewCard("guard", Sanctum, Creature, Common,
		WithPower(6), WithHazardous(5), WithAbility(TriggerAfterDestroyedFighting, gain)), 1)
	before2 := g2.Aember(1)
	if err := g2.Fight(0, weak, guard); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g2.inPlay(weak) {
		t.Fatal("weak attacker should be destroyed by Hazardous")
	}
	if g2.Aember(1) != before2+1 {
		t.Errorf("defender-controller aember = %d, want %d", g2.Aember(1), before2+1)
	}

	// Both combatants survive (Skirmish attacker, tougher defender): no ability.
	g3 := started(t)
	brawler := g3.AddToBattleline(testCreature("brawler", 3, WithKeywords(Skirmish), WithAbility(TriggerAfterDestroyedFighting, gain)), 0)
	wall := g3.AddToBattleline(testCreature("wall", 10), 1)
	before3 := g3.Aember(0)
	if err := g3.Fight(0, brawler, wall); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if !g3.inPlay(brawler) || !g3.inPlay(wall) {
		t.Fatal("both combatants should survive")
	}
	if g3.Aember(0) != before3 {
		t.Errorf("no destruction: aember = %d, want %d (unchanged)", g3.Aember(0), before3)
	}
}

func TestStunBehavior(t *testing.T) {
	g := started(t)

	// Reap: a stunned creature recovers instead of reaping (no Æmber gained).
	r := g.AddToBattleline(testCreature("stunreap", 2), 0)
	g.State.Cards[r].Stunned = true
	before := g.State.Aember[0]
	if err := g.Reap(0, r); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.State.Cards[r].Stunned {
		t.Error("reaping a stunned creature should clear its stun")
	}
	if g.State.Aember[0] != before {
		t.Error("reaping a stunned creature should not gain Æmber")
	}
	if !g.State.Cards[r].Exhausted {
		t.Error("recovering from stun should exhaust the creature")
	}

	// Fight: a stunned attacker recovers instead of fighting (defender unharmed).
	a := g.AddToBattleline(testCreature("stunfight", 5), 0)
	d := g.AddToBattleline(testCreature("victim", 4), 1)
	g.State.Cards[a].Stunned = true
	if err := g.Fight(0, a, d); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.State.Cards[a].Stunned {
		t.Error("fighting with a stunned creature should clear its stun")
	}
	if g.Damage(d) != 0 {
		t.Errorf("stunned attacker should deal no damage; defender took %d", g.Damage(d))
	}

	// Action: a stunned creature recovers instead of using its action ability.
	act := g.AddToBattleline(testCreature("stunact", 3, WithAbility(TriggerAction, GainAember{Amount: 5})), 0)
	g.State.Cards[act].Stunned = true
	before = g.State.Aember[0]
	if err := g.UseAction(0, act); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.State.Cards[act].Stunned {
		t.Error("using a stunned creature's action should clear its stun")
	}
	if g.State.Aember[0] != before {
		t.Error("a stunned action ability should not resolve")
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

	// Armor absorbs damage (raw application, no destruction).
	armored := g.AddToBattleline(testCreature("armored", 5, WithArmor(2)), 0)
	g.DealDamage(0, []DamageTarget{{ID: armored, Amount: 3}}) // 2 absorbed, 1 through
	if g.Damage(armored) != 1 {
		t.Errorf("armored damage = %d, want 1", g.Damage(armored))
	}
	g.DealDamage(0, []DamageTarget{{ID: armored, Amount: 0}}) // no-op path

	// Poison: dealing any damage destroys the creature as part of the process.
	poison := g.AddToBattleline(testCreature("poison", 5, WithKeywords(Poison)), 1)
	g.DealDamage(0, []DamageTarget{{ID: poison, Amount: 1}})
	if g.inPlay(poison) {
		t.Error("poisoned creature should be destroyed")
	}
}

// attachUpgrade registers an upgrade and attaches it to a host creature.
func attachUpgrade(g *Game, host LocalID, def CardDefinition) {
	up := g.Register(def, g.owner(host))
	core := &g.State.Cards[host]
	core.Upgrades[core.UpgradeCount] = up
	core.UpgradeCount++
}

func TestAssaultAndHazardous(t *testing.T) {
	// Assault: the attacker deals its Assault to the defender before fight damage.
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(NewCard("assailant", Brobnar, Creature, Common, WithPower(3), WithAssault(2)), 0)
	def := g.AddToBattleline(testCreature("beefy", 10), 1)
	g.fight(att, def)
	if g.Damage(def) != 5 { // 2 assault + 3 fight
		t.Errorf("assault: defender damage = %d, want 5 (2 assault + 3 fight)", g.Damage(def))
	}

	// Hazardous can destroy the attacker before any combat damage is exchanged.
	g2 := NewGame("A", "B", 1)
	frail := g2.AddToBattleline(testCreature("frail", 4), 0)
	thorns := g2.AddToBattleline(NewCard("thorns", Untamed, Creature, Common, WithPower(6), WithHazardous(5)), 1)
	g2.fight(frail, thorns)
	if g2.inPlay(frail) {
		t.Error("hazardous should destroy the frail attacker before combat")
	}
	if g2.Damage(thorns) != 0 {
		t.Errorf("no combat should occur; thorns damage = %d, want 0", g2.Damage(thorns))
	}

	// Assault granted by an upgrade stacks (the accessor sums upgrade bonuses).
	g3 := NewGame("A", "B", 1)
	host := g3.AddToBattleline(testCreature("host", 3), 0)
	wall := g3.AddToBattleline(testCreature("wall", 10), 1)
	attachUpgrade(g3, host, NewCard("bearway", Untamed, Upgrade, Common, WithStatic(StaticModifier{AssaultBonus: 2})))
	g3.fight(host, wall)
	if g3.Damage(wall) != 5 {
		t.Errorf("upgraded assault: wall damage = %d, want 5", g3.Damage(wall))
	}

	// Hazardous granted by an upgrade stacks too.
	g4 := NewGame("A", "B", 1)
	weak := g4.AddToBattleline(testCreature("weak", 1), 0)
	guard := g4.AddToBattleline(testCreature("guard", 6), 1)
	attachUpgrade(g4, guard, NewCard("flames", Dis, Upgrade, Common, WithStatic(StaticModifier{HazardousBonus: 2})))
	g4.fight(weak, guard)
	if g4.inPlay(weak) {
		t.Error("upgraded hazardous should destroy the weak attacker")
	}
}

func TestDestroyMovesUpgradesToDiscard(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 1), 0)
	g.AddToHand(exBruteStrength(), 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil {
		t.Fatal(err)
	}
	// Host now has 6 power (1 + 5). Kill it with enough damage.
	g.DealDamage(0, []DamageTarget{{ID: host, Amount: 6}})
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
	g.DestroyEach(0, []LocalID{c})
	// The Æmber on the creature goes to its owner's opponent (player 1).
	if g.Aember(1) != 2 {
		t.Errorf("opponent aember = %d, want 2", g.Aember(1))
	}
	if g.AmberOn(c) != 0 {
		t.Errorf("amber not cleared: %d", g.AmberOn(c))
	}
}

func TestGrantedAbilitiesFireFromUpgrades(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 5), 0)
	// An upgrade granting its host "Reap: Gain 1 Æmber".
	attachUpgrade(g, host, NewCard("charm", Shadows, Upgrade, Common,
		WithStatic(StaticModifier{Granted: []Ability{
			{Trigger: TriggerAfterReap, Effect: GainAember{Amount: 1}},
		}})))

	// A non-matching trigger fires nothing granted.
	g.triggerAbilities(host, TriggerAfterFight, 0, false)
	if g.Aember(0) != 0 {
		t.Fatalf("granted Reap ability fired on a Fight trigger; aember = %d", g.Aember(0))
	}
	// The matching trigger fires the granted ability with the host as source.
	g.triggerAbilities(host, TriggerAfterReap, 0, false)
	if g.Aember(0) != 1 {
		t.Errorf("granted Reap ability aember = %d, want 1", g.Aember(0))
	}
}
