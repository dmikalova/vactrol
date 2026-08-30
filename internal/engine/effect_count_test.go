package engine

import "testing"

func TestInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("m1", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("b1", Brobnar, Creature, Common, WithPower(2)), 0)
	g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	g.AddToBattleline(NewCard("foe", Shadows, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Value across type and house filters, and the opposing side.
	values := []struct {
		name string
		in   InPlay
		want int
	}{
		{"friendly creatures", InPlay{Player: Controller, Type: Creature}, 3},
		{"friendly Mars creatures", InPlay{Player: Controller, Type: Creature, House: Mars}, 2},
		{"friendly artifacts", InPlay{Player: Controller, Type: Artifact}, 1},
		{"friendly cards, any type", InPlay{Player: Controller}, 4},
		{"enemy creatures", InPlay{Player: Opponent, Type: Creature}, 1},
	}
	for _, tc := range values {
		if got := tc.in.Value(ctx); got != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.name, got, tc.want)
		}
	}

	// CountText — the singular "for each" noun.
	texts := []struct {
		in   InPlay
		want string
	}{
		{InPlay{Player: Controller, Type: Creature}, "friendly creature in play"},
		{InPlay{Player: Controller, Type: Creature, House: Mars}, "friendly Mars creature"},
		{InPlay{Player: Opponent, Type: Creature}, "enemy creature in play"},
		{InPlay{Player: Controller, Type: Artifact}, "friendly artifact in play"},
		{InPlay{Player: Controller}, "friendly card in play"},
	}
	for _, tc := range texts {
		if got := tc.in.CountText(); got != tc.want {
			t.Errorf("CountText = %q, want %q", got, tc.want)
		}
	}

	// CondText — singular and plural.
	if got := (InPlay{Player: Controller, Type: Creature}).CondText(); got != "if there is a friendly creature in play" {
		t.Errorf("singular CondText = %q", got)
	}
	if got := (InPlay{Player: Controller, Type: Creature, Amount: 2}).CondText(); got != "if there are 2 friendly creatures in play" {
		t.Errorf("plural CondText = %q", got)
	}

	// Met — Amount defaults to one; a higher threshold may not be reached.
	if !(InPlay{Player: Controller, Type: Creature}).Met(ctx) {
		t.Error("default threshold should be met with 3 creatures")
	}
	if !(InPlay{Player: Controller, Type: Creature, Amount: 3}).Met(ctx) {
		t.Error("threshold 3 should be met with 3 creatures")
	}
	if (InPlay{Player: Controller, Type: Creature, Amount: 4}).Met(ctx) {
		t.Error("threshold 4 should not be met with 3 creatures")
	}
}
