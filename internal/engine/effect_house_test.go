package engine

import "testing"

func TestReadyAndBelongToHouseAfterYouPlayCreature(t *testing.T) {
	e := ReadyAndBelongToHouseAfterYouPlayCreature{House: Mars}
	if got := e.Text(); got != "after you play a Mars creature, ready {self} and for the remainder of the turn it belongs to house Mars" {
		t.Errorf("text = %q", got)
	}
	if got := RenderAbility(Ability{Trigger: TriggerAfterCreatureEnters, Effect: e}); got != "After you play a Mars creature, ready {self} and for the remainder of the turn it belongs to house Mars." {
		t.Errorf("ability text = %q", got)
	}
	if got := RenderAbility(Ability{Trigger: TriggerAfterReap, Effect: e}); got != "Reap: After you play a Mars creature, ready {self} and for the remainder of the turn it belongs to house Mars." {
		t.Errorf("non-matching trigger text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Fatalf("validate = %v", err)
	}
	if err := (ReadyAndBelongToHouseAfterYouPlayCreature{}).validate(); err == nil {
		t.Fatal("validate without house = nil, want error")
	}

	g := NewGame("Alice", "Bob", 1)
	host := g.AddToBattleline(NewCard("Host", Brobnar, Creature, Common, WithPower(4)), 0)
	friendlyMars := g.AddToBattleline(NewCard("Martian", Mars, Creature, Common, WithPower(2)), 0)
	friendlyLogos := g.AddToBattleline(NewCard("Thinker", Logos, Creature, Common, WithPower(2)), 0)
	enemyMars := g.AddToBattleline(NewCard("Enemy Martian", Mars, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Source: host, Controller: 0}

	g.State.Cards[host].Exhausted = true
	e.Resolve(ctx)
	if !g.Exhausted(host) || g.House(host) != Brobnar {
		t.Fatalf("missing trigger should do nothing: exhausted=%v house=%s", g.Exhausted(host), g.House(host))
	}

	ctx.HasIt = true
	ctx.It = friendlyLogos
	e.Resolve(ctx)
	if !g.Exhausted(host) || g.House(host) != Brobnar {
		t.Fatalf("non-Mars trigger should do nothing: exhausted=%v house=%s", g.Exhausted(host), g.House(host))
	}

	ctx.It = enemyMars
	e.Resolve(ctx)
	if !g.Exhausted(host) || g.House(host) != Brobnar {
		t.Fatalf("enemy Mars trigger should do nothing: exhausted=%v house=%s", g.Exhausted(host), g.House(host))
	}

	ctx.It = friendlyMars
	e.Resolve(ctx)
	if g.Exhausted(host) || g.House(host) != Mars {
		t.Fatalf("friendly Mars trigger should ready host and make it Mars: exhausted=%v house=%s", g.Exhausted(host), g.House(host))
	}

	g.EndTurn(0)
	if got := g.House(host); got != Brobnar {
		t.Fatalf("end turn house = %s, want Brobnar", got)
	}

	g.BelongToHouseForRemainderOfTurn(host, Mars)
	g.State.Battleline[0].remove(host)
	if got := g.House(host); got != Brobnar {
		t.Fatalf("out-of-play house = %s, want printed Brobnar", got)
	}
}
