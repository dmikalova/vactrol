package engine

import "fmt"

// CardDefinition is the immutable blueprint for a card. Definitions are shared,
// read-only data held in the match Catalog; all mutable per-match state lives in
// the flat GameState. This split keeps the runtime state a plain value that can
// be copied cheaply (see GameState.FastCopy).
type CardDefinition struct {
	Name   string
	House  House
	Type   CardType
	Rarity Rarity
	Traits []Trait

	// Creature stats. Zero for non-creatures.
	Power int
	Armor int

	Keywords []Keyword

	// Assault deals this much damage to the creature this one attacks, before
	// fight damage; Hazardous deals this much to a creature that attacks this one.
	// Zero means the creature does not have that numeric keyword.
	Assault   int
	Hazardous int

	// AemberBonus is the number of Æmber "pips" printed on the card; the
	// controller gains this much Æmber when the card is played.
	AemberBonus int

	// Static is a continuous modifier an Upgrade applies to its host creature.
	Static StaticModifier

	// Abilities are the triggered abilities on the card.
	Abilities []Ability
}

// StaticModifier is a continuous change applied by an Upgrade to the creature it
// is attached to.
type StaticModifier struct {
	PowerBonus     int
	ArmorBonus     int
	AssaultBonus   int
	HazardousBonus int

	// Granted are triggered abilities the Upgrade grants its host creature. The
	// host fires them as if they were printed on it (see Game.triggerAbilities).
	Granted []Ability
}

// Ability pairs a trigger with the effect that resolves when it fires.
type Ability struct {
	Trigger Trigger
	Effect  Effect
}

// hasKeyword reports whether the definition has the given keyword.
func (d *CardDefinition) hasKeyword(k Keyword) bool {
	for _, kw := range d.Keywords {
		if kw == k {
			return true
		}
	}
	return false
}

// hasTrait reports whether the definition has the given trait.
func (d *CardDefinition) hasTrait(t Trait) bool {
	for _, tr := range d.Traits {
		if tr == t {
			return true
		}
	}
	return false
}

// hasTrigger reports whether the definition has an ability with the trigger.
func (d *CardDefinition) hasTrigger(t Trigger) bool {
	for _, ab := range d.Abilities {
		if ab.Trigger == t {
			return true
		}
	}
	return false
}

// CardOption configures a CardDefinition. Definitions use the functional options
// pattern so optional fields read clearly and defaults are centralized.
type CardOption func(*CardDefinition)

// NewCard builds a CardDefinition. Required fields (including rarity) are
// positional; everything optional is supplied via options.
func NewCard(name string, house House, ct CardType, rarity Rarity, opts ...CardOption) CardDefinition {
	c := CardDefinition{
		Name:   name,
		House:  house,
		Type:   ct,
		Rarity: rarity,
	}
	for _, opt := range opts {
		opt(&c)
	}
	for _, ab := range c.Abilities {
		if err := validateEffect(ab.Effect); err != nil {
			panic(fmt.Sprintf("card %q: %v", name, err))
		}
	}
	return c
}

// WithPower sets a creature's power.
func WithPower(p int) CardOption { return func(c *CardDefinition) { c.Power = p } }

// WithArmor sets a creature's armor.
func WithArmor(a int) CardOption { return func(c *CardDefinition) { c.Armor = a } }

// WithTraits appends traits to the card.
func WithTraits(traits ...Trait) CardOption {
	return func(c *CardDefinition) { c.Traits = append(c.Traits, traits...) }
}

// WithKeywords appends keywords to the card.
func WithKeywords(keywords ...Keyword) CardOption {
	return func(c *CardDefinition) { c.Keywords = append(c.Keywords, keywords...) }
}

// WithAssault gives a creature Assault N: it deals N damage to the creature it
// attacks, before fight damage.
func WithAssault(n int) CardOption { return func(c *CardDefinition) { c.Assault = n } }

// WithHazardous gives a creature Hazardous N: a creature that attacks it is dealt
// N damage before fight damage.
func WithHazardous(n int) CardOption { return func(c *CardDefinition) { c.Hazardous = n } }

// WithAemberBonus sets the number of Æmber pips on the card.
func WithAemberBonus(n int) CardOption { return func(c *CardDefinition) { c.AemberBonus = n } }

// WithStatic sets the continuous modifier an Upgrade applies to its host.
func WithStatic(m StaticModifier) CardOption { return func(c *CardDefinition) { c.Static = m } }

// WithAbility appends a triggered ability to the card.
func WithAbility(trigger Trigger, effect Effect) CardOption {
	return func(c *CardDefinition) {
		c.Abilities = append(c.Abilities, Ability{Trigger: trigger, Effect: effect})
	}
}
