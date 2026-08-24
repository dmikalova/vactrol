package game

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

// RenderCardText renders the full multi-line printed text of a card, matching the
// layout used on physical cards: house, type, rarity, stats, bonuses, traits, and
// ability lines.
func RenderCardText(def *CardDefinition) string {
	var lines []string
	lines = append(lines, def.House.String(), string(def.Type), string(def.Rarity))

	if def.Type == Creature {
		lines = append(lines, fmt.Sprintf("%d Power", def.Power))
		lines = append(lines, fmt.Sprintf("%d Armor", def.Armor))
	}

	if def.AemberBonus > 0 {
		lines = append(lines, fmt.Sprintf("%d Æmber", def.AemberBonus))
	}

	if len(def.Traits) > 0 {
		traits := make([]string, len(def.Traits))
		for i, t := range def.Traits {
			traits[i] = string(t)
		}
		lines = append(lines, strings.Join(traits, " • "))
	}

	if s := staticText(def.Static); s != "" {
		lines = append(lines, s)
	}

	for _, ab := range def.Abilities {
		lines = append(lines, RenderAbility(ab))
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
	return "This creature gets " + strings.Join(parts, " and ") + "."
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
