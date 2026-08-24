package game

import "testing"

func TestTriggerPrefixDefault(t *testing.T) {
	// An unknown trigger renders with no prefix and a capitalized effect.
	got := RenderAbility(Ability{Trigger: Trigger(99), Effect: GainAember{Amount: 1}})
	if got != "Gain 1 Æmber." {
		t.Errorf("RenderAbility(unknown) = %q", got)
	}
}

func TestTargetTextDefault(t *testing.T) {
	cases := map[TargetKind]string{
		TargetThisCreature:       "this creature",
		TargetTriggeringCreature: "it",
		TargetEachEnemyCreature:  "each enemy creature",
		TargetKind(99):           "a creature",
	}
	for kind, want := range cases {
		if got := (Target{Kind: kind}).Text(); got != want {
			t.Errorf("Target{%d}.Text() = %q, want %q", kind, got, want)
		}
	}
}

func TestAllTriggerPrefixes(t *testing.T) {
	cases := map[Trigger]string{
		TriggerPlay:                "Play: Gain 1 Æmber.",
		TriggerReap:                "Reap: Gain 1 Æmber.",
		TriggerFight:               "Fight: Gain 1 Æmber.",
		TriggerAction:              "Action: Gain 1 Æmber.",
		TriggerDestroyed:           "Destroyed: Gain 1 Æmber.",
		TriggerAfterForgeKey:       "After you forge a key, gain 1 Æmber.",
		TriggerAfterCreatureEnters: "After a creature enters play, gain 1 Æmber.",
	}
	for tr, want := range cases {
		got := RenderAbility(Ability{Trigger: tr, Effect: GainAember{Amount: 1}})
		if got != want {
			t.Errorf("trigger %d prefix = %q, want %q", tr, got, want)
		}
	}
}

func TestCapitalizeFirst(t *testing.T) {
	if capitalizeFirst("") != "" {
		t.Error("capitalizeFirst(\"\") should be empty")
	}
	if capitalizeFirst("hello") != "Hello" {
		t.Errorf("capitalizeFirst(hello) = %q", capitalizeFirst("hello"))
	}
}

func TestStaticText(t *testing.T) {
	if staticText(StaticModifier{}) != "" {
		t.Error("empty static modifier should render empty")
	}
	got := staticText(StaticModifier{PowerBonus: 5, ArmorBonus: 2})
	if got != "This creature gets +5 power and +2 armor." {
		t.Errorf("staticText = %q", got)
	}
}

func TestGeneratedCardText(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		{exGiant(), "Brobnar\nCreature\nRare\n5 Power\n0 Armor\nGiant\nAfter you forge a key, deal 2 damage to each enemy creature."},
		{exBruteStrength(), "Brobnar\nUpgrade\nUncommon\n1 Æmber\nThis creature gets +5 power."},
		{exBattleFury(), "Brobnar\nAction\nCommon\n1 Æmber\nPlay: Ready and fight with a friendly creature."},
		{exAutocannon(), "Brobnar\nArtifact\nRare\n1 Æmber\nWeapon\nAfter a creature enters play, deal 1 damage to it."},
	}
	for _, tc := range cases {
		if got := RenderCardText(&tc.def); got != tc.want {
			t.Errorf("%s text mismatch:\n got:\n%s\nwant:\n%s", tc.def.Name, got, tc.want)
		}
	}
}
