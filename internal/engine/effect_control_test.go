package engine

import "testing"

func TestTakeControl(t *testing.T) {
	if got := (TakeControl{}).Text(); got != "take control of this creature until {upgrade} leaves play" {
		t.Fatalf("TakeControl text = %q", got)
	}

	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 1)

	(TakeControl{}).Resolve(&EffectContext{Resolver: g, Source: host, Controller: 0})

	if g.owner(host) != 1 {
		t.Fatalf("owner = %d, want unchanged P2", g.owner(host))
	}
	if g.controller(host) != 0 {
		t.Fatalf("controller = %d, want P1", g.controller(host))
	}
	if !g.State.Battleline[0].contains(host) || g.State.Battleline[1].contains(host) {
		t.Fatalf("battlelines = %v/%v, want host only under P1", g.Battleline(0), g.Battleline(1))
	}
}

func TestControlRevertsWhenTakingUpgradeLeaves(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 1)
	collar := g.Register(NewCard("collar", Dis, Upgrade, Rare, WithStatic(StaticModifier{TakesControl: true})), 0)
	core := &g.State.Cards[host]
	core.Upgrades[0] = collar
	core.UpgradeCount = 1
	g.takeControl(host, 0)

	g.discardAttachedUpgradeAt(host, 0)

	if g.controller(host) != 1 {
		t.Fatalf("controller = %d, want owner P2", g.controller(host))
	}
	if !g.State.Battleline[1].contains(host) || g.State.Battleline[0].contains(host) {
		t.Fatalf("battlelines = %v/%v, want host only under P2", g.Battleline(0), g.Battleline(1))
	}
	if !g.State.Discard[0].contains(collar) {
		t.Fatalf("collar discard = %v, want collar in owner discard", g.Discard(0))
	}
}

func TestControlledCreatureLeavesForOwner(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 1)
	collar := g.Register(NewCard("collar", Dis, Upgrade, Rare, WithStatic(StaticModifier{TakesControl: true})), 0)
	core := &g.State.Cards[host]
	core.Upgrades[0] = collar
	core.UpgradeCount = 1
	g.takeControl(host, 0)

	g.DealDamage(0, []DamageTarget{{ID: host, Amount: 3}})

	if !g.State.Discard[1].contains(host) {
		t.Fatalf("host discard = %v/%v, want owner P2 discard", g.Discard(0), g.Discard(1))
	}
	if !g.State.Discard[0].contains(collar) {
		t.Fatalf("collar discard = %v, want owner P1 discard", g.Discard(0))
	}
	if g.State.Cards[host].ControlPlus != 0 {
		t.Fatalf("ControlPlus after leaving play = %d, want 0", g.State.Cards[host].ControlPlus)
	}
}
