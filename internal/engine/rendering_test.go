package engine

import "testing"

func TestTriggerPrefixDefault(t *testing.T) {
	// An unknown trigger renders with no prefix and a capitalized effect.
	got := RenderAbility(Ability{Trigger: Trigger(99), Effect: GainAember{Player: Controller, Amount: 1}})
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
		TriggerAfterArtifactPlayed:    "After you play an artifact, gain 1 Æmber.",
	}
	for tr, want := range cases {
		got := RenderAbility(Ability{Trigger: tr, Effect: GainAember{Player: Controller, Amount: 1}})
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
		{NewCard("Valdr", Brobnar, Creature, Common, WithPower(6), WithTraits("Giant"), WithAttackDamage(AttackDamage{Amount: 2, FlankOnly: true})), "House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  6\nTraits: Giant\n\nValdr deals +2 Damage while attacking an enemy creature on the flank."},
		{NewCard("Spider", Mars, Creature, Common, WithPower(7), WithAttackDamage(AttackDamage{Fixed: true, Amount: 0})), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  7\n\nSpider deals no damage when fighting."},
		{NewCard("Bruiser", Brobnar, Creature, Common, WithPower(8), WithAttackDamage(AttackDamage{Fixed: true, Amount: 5})), "House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  8\n\nBruiser deals 5 damage when fighting."},
		{NewCard("Basher", Brobnar, Creature, Common, WithPower(4), WithAttackDamage(AttackDamage{Amount: 2})), "House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  4\n\nBasher deals +2 Damage when fighting."},
		{NewCard("Runner", Shadows, Upgrade, Uncommon, WithStatic(StaticModifier{Granted: []Ability{{Trigger: TriggerAfterReap, Effect: StealAember{Amount: 1}}}})), "House:  Shadows\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains, \"Reap: Steal 1 Æmber.\""},
		{NewCard("Boots", Logos, Upgrade, Uncommon, WithStatic(StaticModifier{Keywords: []Keyword{Versatile}}), WithAbility(TriggerAfterPlay, Sequence{Effects: []Effect{Stun{Target: Target{Kind: TargetThisCreature}}, Exhaust{Target: Target{Kind: TargetThisCreature}}}})), "House:  Logos\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains versatile.\nPlay: Stun and exhaust this creature."},
		{NewCard("Jammer", Mars, Creature, Common, WithPower(4), WithArmor(1), WithTraits("Robot"), WithKeyCost(NewKeyCostChange(Opponent, 1)), WithAbility(TriggerAfterReap, CaptureAember{Amount: 1}), WithAbility(TriggerAfterFight, CaptureAember{Amount: 1})), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  4\nArmor:  1\nTraits: Robot\n\nYour opponent's keys cost +1 Æmber.\nFight/Reap: Jammer captures 1 Æmber."},
		{NewCard("Pack", Mars, Upgrade, Uncommon, WithAemberBonus(1), WithStatic(StaticModifier{KeyCostChange: NewKeyCostChange(Opponent, 2)})), "House:  Mars\nType:   Upgrade\nRarity: Uncommon\nÆmber:  1\n\nThis creature gains, \"Your opponent's keys cost +2 Æmber.\""},
		{NewCard("SelfTax", Mars, Creature, Common, WithPower(3), WithKeyCost(NewKeyCostChange(Controller, 1))), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  3\n\nYour keys cost +1 Æmber."},
		{NewCard("Tax", Mars, Creature, Common, WithPower(3), WithKeyCost(NewKeyCostChange(EachPlayer, 1))), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  3\n\nEach player's keys cost +1 Æmber."},
		{NewCard("Imp", Dis, Creature, Common, WithPower(2), WithTraits("Imp"), WithRestrictions(Restrictions{PlayCardLimit: PlayCardLimit{Player: Opponent, Amount: 2}})), "House:  Dis\nType:   Creature\nRarity: Common\nPower:  2\nTraits: Imp\n\nYour opponent cannot play more than 2 cards each turn."},
		{NewCard("Twig", Untamed, Creature, Common, WithPower(7), WithTraits("Beast"), WithFightRestriction(Target{Kind: TargetEachCreature}.Stunned())), "House:  Untamed\nType:   Creature\nRarity: Common\nPower:  7\nTraits: Beast\n\nTwig can only fight stunned creatures."},
		{NewCard("Ritual", Dis, Artifact, Rare, WithTraits("Power"), WithConstantAbility(ConstantAbility{Target: Target{Kind: TargetEachCreature}, Granted: []Ability{{Trigger: TriggerDestroyed, Effect: PurgeCreature{Target: Target{Kind: TargetThisCreature}}}}})), "House:  Dis\nType:   Artifact\nRarity: Rare\nTraits: Power\n\nEach creature gains, \"Destroyed: Purge this creature.\""},
	}
	for _, tc := range cases {
		if got := RenderCardText(&tc.def); got != tc.want {
			t.Errorf("%s text mismatch:\n got:\n%s\nwant:\n%s", tc.def.Name, got, tc.want)
		}
	}
}

func TestRenderCardRules(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		// Vanilla creature: no rules text.
		{NewCard("Grunt", Brobnar, Creature, Common, WithPower(4)), ""},
		// Keyword line only.
		{NewCard("Asp", Shadows, Creature, Uncommon, WithPower(3), WithKeywords(Skirmish, Poison)), "Skirmish. Poison."},
		// Upgrade static modifier (no own ability): the face still shows its text.
		{exBruteStrength(), "This creature gains +5 power."},
		// Triggered ability with a self-reference resolved to the card's name.
		{NewCard("Tabris", Sanctum, Creature, Uncommon, WithPower(6), WithAbility(TriggerAfterFight, CaptureAember{Amount: 1})), "Fight: Tabris captures 1 Æmber."},
	}
	for _, tc := range cases {
		if got := RenderCardRules(&tc.def); got != tc.want {
			t.Errorf("%s rules mismatch:\n got:  %q\n want: %q", tc.def.Name, got, tc.want)
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

func TestPunctuate(t *testing.T) {
	cases := map[string]string{
		`gains, "Play: x."`: `gains, "Play: x."`, // quoted, already has inner period
		`gains, "Play: x"`:  `gains, "Play: x."`, // quoted, period tucked inside
		"deal damage.":      "deal damage.",      // already ends with a period
		"deal damage":       "deal damage.",      // period appended
	}
	for in, want := range cases {
		if got := punctuate(in); got != want {
			t.Errorf("punctuate(%q) = %q, want %q", in, got, want)
		}
	}
}
