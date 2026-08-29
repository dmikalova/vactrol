package engine

import (
	"fmt"
	"strings"
	"unicode"
)

// RenderAbility renders a single triggered ability to its printed card line,
// e.g. "After you forge a key, deal 2 damage to each enemy creature."
func RenderAbility(a Ability) string {
	prefix, capitalize := a.Trigger.prefix()
	body := a.Effect.Text()
	if capitalize {
		body = capitalizeFirst(body)
	}
	return prefix + punctuate(body)
}

// punctuate ends an ability body with a period. A body that already ends in a
// period is left alone; one that ends in a closing quote (an embedded ability
// such as Charge!'s granted "Play: ...") takes its period inside the quote, so
// the line reads `... an enemy creature."` rather than doubling or misplacing it.
func punctuate(body string) string {
	if strings.HasSuffix(body, `"`) {
		if inner := body[:len(body)-1]; !strings.HasSuffix(inner, ".") {
			return inner + `."`
		}
		return body
	}
	if strings.HasSuffix(body, ".") {
		return body
	}
	return body + "."
}

// renderAbilityLine renders an ability with its source card's self-references
// (the SelfName placeholder) resolved to the card's name.
func renderAbilityLine(def *CardDefinition, a Ability) string {
	return strings.ReplaceAll(RenderAbility(a), SelfName, def.Name)
}

// abilityLines renders a card's triggered abilities, one printed line each,
// resolving self-references to the card's name (or "this creature" for an
// upgrade's own abilities). A creature's adjacent Reap and Fight abilities with
// the same effect — the pair card.WithFightOrReap adds — merge into one
// "Fight/Reap:" line, the KeyForge shorthand.
func abilityLines(def *CardDefinition) []string {
	self := def.Name
	if def.Type == Upgrade {
		self = "this creature"
	}
	abs := def.Abilities
	lines := make([]string, 0, len(abs))
	for i := 0; i < len(abs); i++ {
		ab := abs[i]
		if ab.Trigger == TriggerAfterReap && i+1 < len(abs) &&
			abs[i+1].Trigger == TriggerAfterFight &&
			abs[i+1].Effect.Text() == ab.Effect.Text() {
			body := capitalizeFirst(strings.ReplaceAll(ab.Effect.Text(), SelfName, self))
			lines = append(lines, "Fight/Reap: "+body+".")
			i++ // the Fight partner prints as part of this line
			continue
		}
		lines = append(lines, strings.ReplaceAll(RenderAbility(ab), SelfName, self))
	}
	return lines
}

// RenderCardText renders a card's details as labeled, colon-aligned lines
// (House, Type, Rarity, stats, Æmber, Traits), followed by the card's rules text
// (keywords, upgrade modifier, and ability lines). Labels are padded by rune
// width so the multi-byte "Æmber" label still aligns.
func RenderCardText(def *CardDefinition) string {
	type field struct{ label, value string }
	fields := []field{
		{"House", def.House.String()},
		{"Type", string(def.Type)},
		{"Rarity", string(def.Rarity)},
	}
	if def.Type == Creature {
		fields = append(fields, field{"Power", fmt.Sprintf("%d", def.Power)})
		if def.Armor > 0 {
			fields = append(fields, field{"Armor", fmt.Sprintf("%d", def.Armor)})
		}
	}
	if def.AemberBonus > 0 {
		fields = append(fields, field{"Æmber", fmt.Sprintf("%d", def.AemberBonus)})
	}
	if len(def.Traits) > 0 {
		traits := make([]string, len(def.Traits))
		for i, t := range def.Traits {
			traits[i] = string(t)
		}
		fields = append(fields, field{"Traits", strings.Join(traits, " • ")})
	}

	// Widest label, measured in runes so "Æmber" (multi-byte Æ) aligns visually.
	width := 0
	for _, f := range fields {
		if n := len([]rune(f.label)); n > width {
			width = n
		}
	}

	var lines []string
	for _, f := range fields {
		pad := strings.Repeat(" ", width-len([]rune(f.label))+1)
		lines = append(lines, f.label+":"+pad+f.value)
	}

	// Rules text (keywords, upgrade modifier, abilities) follows the labeled
	// header, separated by a blank line.
	if rules := cardRules(def); len(rules) > 0 {
		lines = append(lines, "")
		lines = append(lines, rules...)
	}

	return strings.Join(lines, "\n")
}

// RenderCardRules renders just a card's rules text — the keyword, restriction,
// key-cost, static-modifier, constant-ability, granted-ability, and triggered-
// ability lines that RenderCardText shows below the labeled header — joined one
// per line with no header. It is the text drawn on the compact card face, so an
// upgrade shows its granted keywords and abilities there too.
func RenderCardRules(def *CardDefinition) string {
	return strings.Join(cardRules(def), "\n")
}

// cardRules assembles a card's rules lines in printed order: keywords, "cannot"
// restrictions, key-cost, an upgrade's static modifier, a constant ability, the
// abilities an upgrade grants its host, and finally the card's own triggered
// abilities. An upgrade's own abilities name their host "this creature" since it
// is unknown at print time.
func cardRules(def *CardDefinition) []string {
	var rules []string
	if s := keywordText(def); s != "" {
		rules = append(rules, s)
	}
	if s := attackDamageText(def); s != "" {
		rules = append(rules, s)
	}
	if fr := def.FightRestriction; fr != (Target{}) {
		rules = append(rules, def.Name+" can only fight "+singularNoun(fr.Text())+"s.")
	}
	rules = append(rules, restrictionText(def.Restricts)...)
	if s := keyCostText(def.KeyCostChange); s != "" {
		rules = append(rules, s)
	}
	if s := staticText(def.Static); s != "" {
		rules = append(rules, s)
	}
	if s := constantText(def); s != "" {
		rules = append(rules, s)
	}
	rules = append(rules, constantGrantedText(def)...)
	rules = append(rules, grantedText(def.Static)...)
	rules = append(rules, abilityLines(def)...)
	return rules
}

// CardDocComment renders a card's details as a Go doc comment block, the form
// `mage generateComments` writes above each card's declaration. It is the card
// name, a blank separator, then RenderCardText's lines, each turned into a
// comment line: the title as "// <Name>", blanks as "//", and detail lines as
// tab-indented "//\t..." so godoc renders the labeled block preformatted. The
// result has no trailing newline.
func CardDocComment(def *CardDefinition) string {
	var b strings.Builder
	b.WriteString("// " + def.Name + "\n//\n")
	for _, line := range strings.Split(RenderCardText(def), "\n") {
		if line == "" {
			b.WriteString("//\n")
		} else {
			b.WriteString("//\t" + line + "\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// staticText renders an Upgrade's continuous modifier, e.g.
// "This creature gets +5 power."
func staticText(m StaticModifier) string {
	var parts []string
	if m.PowerBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d power", m.PowerBonus))
	}
	if m.ArmorBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d armor", m.ArmorBonus))
	}
	if m.AssaultBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d assault", m.AssaultBonus))
	}
	if m.HazardousBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d hazardous", m.HazardousBonus))
	}
	for _, kw := range m.Keywords {
		parts = append(parts, strings.ToLower(string(kw)))
	}
	if len(parts) == 0 {
		return ""
	}
	return "This creature gains " + strings.Join(parts, " and ") + "."
}

// grantedText renders the triggered abilities an Upgrade grants its host, one
// line each, e.g. `This creature gains, "Reap: Steal 1 Æmber."`. Self-references
// resolve to "this creature" since the host is unknown when the Upgrade prints.
func grantedText(m StaticModifier) []string {
	lines := make([]string, 0, len(m.Granted))
	for _, ab := range m.Granted {
		body := strings.ReplaceAll(RenderAbility(ab), SelfName, "this creature")
		lines = append(lines, `This creature gains, "`+body+`"`)
	}
	if s := keyCostText(m.KeyCostChange); s != "" {
		lines = append(lines, `This creature gains, "`+s+`"`)
	}
	return lines
}

// constantText renders a card's constant ability, e.g. "Each friendly creature
// gains +1 power." or "Each neighboring creature gains +2 armor." The subject is
// the constant ability's Target ("each creature" when unset). Returns "" when the
// card has no constant ability.
func constantText(def *CardDefinition) string {
	c := def.Constant
	var parts []string
	if c.PowerBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d power", c.PowerBonus))
	}
	if c.ArmorBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d armor", c.ArmorBonus))
	}
	if len(parts) == 0 {
		return ""
	}
	who := capitalizeFirst(c.target().Text())
	line := who + " gains " + strings.Join(parts, " and ") + "."
	return strings.ReplaceAll(line, SelfName, def.Name)
}

// constantGrantedText renders the triggered abilities a card's constant ability
// grants the creatures it reaches, one line each, e.g.
// `Each creature gains, "Destroyed: purge this creature."`. Self-references in the
// granted ability resolve to "this creature", the creature that gains it.
func constantGrantedText(def *CardDefinition) []string {
	c := def.Constant
	if len(c.Granted) == 0 {
		return nil
	}
	subject := capitalizeFirst(c.target().Text())
	lines := make([]string, 0, len(c.Granted))
	for _, ab := range c.Granted {
		body := strings.ReplaceAll(RenderAbility(ab), SelfName, "this creature")
		lines = append(lines, subject+` gains, "`+body+`"`)
	}
	return lines
}

// restrictionText renders a card's constant "cannot" rules, one line each, e.g.
// "You cannot play creatures." Returns nil when the card imposes none.
func restrictionText(r Restrictions) []string {
	var lines []string
	if r.Fighting {
		lines = append(lines, "You cannot use creatures to fight.")
	}
	if r.CannotPlay != "" {
		lines = append(lines, "You cannot play "+strings.ToLower(string(r.CannotPlay))+"s.")
	}
	if l := r.PlayCardLimit; l.Amount > 0 {
		who := "You"
		switch l.Player {
		case Opponent:
			who = "Your opponent"
		case EachPlayer:
			who = "Each player"
		}
		lines = append(lines, fmt.Sprintf("%s cannot play more than %d cards each turn.", who, l.Amount))
	}
	return lines
}

// keyCostText renders a card's key-cost change, e.g. "Your opponent's keys cost +1
// Æmber." (or "Your keys cost…" / "Each player's keys cost…"). Returns "" when the
// change is zero. Both a card that prints the rule and an Upgrade that grants it
// use this same sentence.
func keyCostText(kc KeyCostChange) string {
	if kc.amount == 0 {
		return ""
	}
	whose := "Your"
	switch kc.player {
	case Opponent:
		whose = "Your opponent's"
	case EachPlayer:
		whose = "Each player's"
	}
	return fmt.Sprintf("%s keys cost %+d Æmber.", whose, kc.amount)
}

// keywordText renders a card's keywords as a single leading line, e.g.
// "Skirmish. Poison." or "Assault 2.". Returns "" when the card has none.
func keywordText(def *CardDefinition) string {
	var parts []string
	for _, k := range def.Keywords {
		parts = append(parts, string(k)+".")
	}
	if def.Assault > 0 {
		parts = append(parts, fmt.Sprintf("Assault %d.", def.Assault))
	}
	if def.Hazardous > 0 {
		parts = append(parts, fmt.Sprintf("Hazardous %d.", def.Hazardous))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

// attackDamageText renders a creature's custom fight damage, e.g. "Valdr deals +2
// Damage while attacking an enemy creature on the flank." or "Ether Spider deals
// no damage when fighting." Returns "" when the creature deals its plain power.
func attackDamageText(def *CardDefinition) string {
	ad := def.AttackDamage
	switch {
	case ad.Fixed && ad.Amount == 0:
		return def.Name + " deals no damage when fighting."
	case ad.Fixed:
		return fmt.Sprintf("%s deals %d damage when fighting.", def.Name, ad.Amount)
	case ad.Amount != 0 && ad.FlankOnly:
		return fmt.Sprintf("%s deals +%d Damage while attacking an enemy creature on the flank.", def.Name, ad.Amount)
	case ad.Amount != 0:
		return fmt.Sprintf("%s deals +%d Damage when fighting.", def.Name, ad.Amount)
	default:
		return ""
	}
}

// capitalizeFirst upper-cases the first rune of s.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
