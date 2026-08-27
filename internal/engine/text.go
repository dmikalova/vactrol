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
	return prefix + body + "."
}

// renderAbilityLine renders an ability with its source card's self-references
// (the SelfName placeholder) resolved to the card's name.
func renderAbilityLine(def *CardDefinition, a Ability) string {
	return strings.ReplaceAll(RenderAbility(a), SelfName, def.Name)
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
		fields = append(fields,
			field{"Power", fmt.Sprintf("%d", def.Power)},
			field{"Armor", fmt.Sprintf("%d", def.Armor)},
		)
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
	var rules []string
	if s := keywordText(def.Keywords); s != "" {
		rules = append(rules, s)
	}
	if s := staticText(def.Static); s != "" {
		rules = append(rules, s)
	}
	for _, ab := range def.Abilities {
		rules = append(rules, renderAbilityLine(def, ab))
	}
	if len(rules) > 0 {
		lines = append(lines, "")
		lines = append(lines, rules...)
	}

	return strings.Join(lines, "\n")
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
	if len(parts) == 0 {
		return ""
	}
	return "This creature gains " + strings.Join(parts, " and ") + "."
}

// keywordText renders a card's keywords as a single leading line, e.g.
// "Skirmish. Poison.". Returns "" when the card has no keywords.
func keywordText(kws []Keyword) string {
	if len(kws) == 0 {
		return ""
	}
	parts := make([]string, len(kws))
	for i, k := range kws {
		parts[i] = string(k) + "."
	}
	return strings.Join(parts, " ")
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
