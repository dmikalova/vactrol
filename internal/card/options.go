package card

import "github.com/dmikalova/vactrol/internal/engine"

// Option helpers — thin wrappers over the engine's card options, each returning
// an authoring Option that appends its engine option to the builder. New (in
// register.go) applies them and enrolls the card.
var (
	WithPower                 = func(p int) Option { return gameplay(engine.WithPower(p)) }
	WithArmor                 = func(a int) Option { return gameplay(engine.WithArmor(a)) }
	WithTraits                = func(t ...Trait) Option { return gameplay(engine.WithTraits(t...)) }
	WithKeywords              = func(k ...engine.Keyword) Option { return gameplay(engine.WithKeywords(k...)) }
	WithAssault               = func(n int) Option { return gameplay(engine.WithAssault(n)) }
	WithHazardous             = func(n int) Option { return gameplay(engine.WithHazardous(n)) }
	WithAttackDamage          = func(ad engine.AttackDamage) Option { return gameplay(engine.WithAttackDamage(ad)) }
	WithFightRestriction      = func(t engine.Target) Option { return gameplay(engine.WithFightRestriction(t)) }
	WithEntersPlay            = func(e Effect) Option { return gameplay(engine.WithEntersPlay(e)) }
	WithAemberBonus           = func(n int) Option { return gameplay(engine.WithAemberBonus(n)) }
	WithStatic                = func(m StaticModifier) Option { return gameplay(engine.WithStatic(m)) }
	WithConstantAbility       = func(c ConstantAbility) Option { return gameplay(engine.WithConstantAbility(c)) }
	WithRestrictions          = func(r Restrictions) Option { return gameplay(engine.WithRestrictions(r)) }
	WithKeyCost               = func(kc engine.KeyCostChange) Option { return gameplay(engine.WithKeyCost(kc)) }
	WithOffHousePlayGrant     = func(h engine.House) Option { return gameplay(engine.WithOffHousePlayGrant(h)) }
	WithCaptureOpponentAember = func() Option {
		return gameplay(engine.WithCaptureOpponentAember())
	}
	WithAbility = func(t engine.Trigger, e Effect) Option { return gameplay(engine.WithAbility(t, e)) }
	// WithFightOrReap adds effect as both a Fight and a Reap ability, so it resolves
	// whenever the creature is used to fight or to reap; the two print as one
	// "Fight/Reap:" line.
	WithFightOrReap = func(e Effect) Option {
		return gameplay(func(d *engine.CardDefinition) {
			engine.WithAbility(engine.TriggerAfterReap, e)(d)
			engine.WithAbility(engine.TriggerAfterFight, e)(d)
		})
	}
)
