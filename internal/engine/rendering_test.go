package engine

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
		TargetThisCreature:       SelfName,
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
		TriggerAfterPlay:              "Play: Gain 1 Æmber.",
		TriggerAfterReap:              "Reap: Gain 1 Æmber.",
		TriggerAfterFight:             "Fight: Gain 1 Æmber.",
		TriggerBeforeFight:            "Before Fight: Gain 1 Æmber.",
		TriggerAction:                 "Action: Gain 1 Æmber.",
		TriggerDestroyed:              "Destroyed: Gain 1 Æmber.",
		TriggerAfterForgeKey:          "After you forge a key, gain 1 Æmber.",
		TriggerAfterCreatureEnters:    "After a creature enters play, gain 1 Æmber.",
		TriggerAfterDestroyedFighting: "After a creature is destroyed fighting {self}, gain 1 Æmber.",
	}
	for tr, want := range cases {
		got := RenderAbility(Ability{Trigger: tr, Effect: GainAember{Amount: 1}})
		if got != want {
			t.Errorf("trigger %d prefix = %q, want %q", tr, got, want)
		}
	}
}

func TestGeneratedCardText(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		{exGiant(), "House:  Brobnar\nType:   Creature\nRarity: Rare\nPower:  5\nTraits: Giant\n\nAfter you forge a key, deal 2 damage to each enemy creature."},
		{exBruteStrength(), "House:  Brobnar\nType:   Upgrade\nRarity: Uncommon\nÆmber:  1\n\nThis creature gains +5 power."},
		{exBattleFury(), "House:  Brobnar\nType:   Action\nRarity: Common\nÆmber:  1\n\nPlay: Ready and fight with a friendly creature."},
		{exAutocannon(), "House:  Brobnar\nType:   Artifact\nRarity: Rare\nÆmber:  1\nTraits: Weapon\n\nAfter a creature enters play, deal 1 damage to it."},
		{NewCard("Asp", Shadows, Creature, Uncommon, WithPower(3), WithKeywords(Skirmish, Poison)), "House:  Shadows\nType:   Creature\nRarity: Uncommon\nPower:  3\n\nSkirmish. Poison."},
		{NewCard("Anaphiel", Sanctum, Creature, Common, WithPower(6), WithArmor(1), WithTraits("Knight"), WithKeywords(Taunt)), "House:  Sanctum\nType:   Creature\nRarity: Common\nPower:  6\nArmor:  1\nTraits: Knight\n\nTaunt."},
		{NewCard("Tabris", Sanctum, Creature, Uncommon, WithPower(6), WithAbility(TriggerAfterFight, CaptureAember{Amount: 1})), "House:  Sanctum\nType:   Creature\nRarity: Uncommon\nPower:  6\n\nFight: Tabris captures 1 Æmber."},
		{NewCard("Bear", Untamed, Creature, Common, WithPower(5), WithTraits("Beast"), WithAssault(2)), "House:  Untamed\nType:   Creature\nRarity: Common\nPower:  5\nTraits: Beast\n\nAssault 2."},
		{NewCard("Grub", Untamed, Creature, Rare, WithPower(2), WithHazardous(5)), "House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  2\n\nHazardous 5."},
		{NewCard("Runner", Shadows, Upgrade, Uncommon, WithStatic(StaticModifier{Granted: []Ability{{Trigger: TriggerAfterReap, Effect: StealAember{Amount: 1}}}})), "House:  Shadows\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains, \"Reap: Steal 1 Æmber.\""},
	}
	for _, tc := range cases {
		if got := RenderCardText(&tc.def); got != tc.want {
			t.Errorf("%s text mismatch:\n got:\n%s\nwant:\n%s", tc.def.Name, got, tc.want)
		}
	}
}

func TestCardDocComment(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		{
			exGiant(),
			"// Brobnar Giant\n//\n//\tHouse:  Brobnar\n//\tType:   Creature\n//\tRarity: Rare\n//\tPower:  5\n//\tTraits: Giant\n//\n//\tAfter you forge a key, deal 2 damage to each enemy creature.",
		},
		{
			NewCard("Asp", Shadows, Creature, Uncommon, WithPower(3), WithKeywords(Skirmish, Poison)),
			"// Asp\n//\n//\tHouse:  Shadows\n//\tType:   Creature\n//\tRarity: Uncommon\n//\tPower:  3\n//\n//\tSkirmish. Poison.",
		},
	}
	for _, tc := range cases {
		if got := CardDocComment(&tc.def); got != tc.want {
			t.Errorf("%s doc comment mismatch:\n got:\n%s\nwant:\n%s", tc.def.Name, got, tc.want)
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
	if got != "This creature gains +5 power and +2 armor." {
		t.Errorf("staticText = %q", got)
	}
	if got := staticText(StaticModifier{AssaultBonus: 2}); got != "This creature gains +2 assault." {
		t.Errorf("assault staticText = %q", got)
	}
	if got := staticText(StaticModifier{PowerBonus: 2, HazardousBonus: 2}); got != "This creature gains +2 power and +2 hazardous." {
		t.Errorf("hazardous staticText = %q", got)
	}
	if got := staticText(StaticModifier{Keywords: []Keyword{Skirmish}}); got != "This creature gains skirmish." {
		t.Errorf("keyword staticText = %q", got)
	}
	if got := staticText(StaticModifier{Keywords: []Keyword{Elusive, Skirmish}}); got != "This creature gains elusive and skirmish." {
		t.Errorf("two-keyword staticText = %q", got)
	}
	if got := staticText(StaticModifier{ArmorBonus: 1, Keywords: []Keyword{Taunt}}); got != "This creature gains +1 armor and taunt." {
		t.Errorf("armor+keyword staticText = %q", got)
	}
}
