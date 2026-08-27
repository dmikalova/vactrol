package engine

import "testing"

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

func TestFightVerbNoEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	FightVerb{}.Apply(ctx, src) // no enemies -> logs and returns
	if g.State.Cards[src].Exhausted {
		t.Error("no fight should have occurred")
	}
}
