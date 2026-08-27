package engine

import "fmt"

// DealDamage deals a fixed amount of damage to each creature its Target selects.
//
// In KeyForge, dealing damage places that many damage tokens on a creature;
// armor reduces each separate instance of damage before it lands. A creature
// whose damage reaches or exceeds its power is destroyed. When one ability deals
// damage to several creatures they are dealt damage simultaneously and any that
// died are destroyed together, so no creature's destruction changes another's
// (see Game.dealDamage).
type DealDamage struct {
	Amount int
	Target Target
}

// Text renders the effect, e.g. "deal 2 damage to each enemy creature".
func (e DealDamage) Text() string {
	return fmt.Sprintf("deal %d damage to %s", e.Amount, e.Target.Text())
}

// Resolve deals the damage to every selected creature simultaneously, resolving
// destruction as part of it.
func (e DealDamage) Resolve(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	targets := make([]DamageTarget, len(ids))
	for i, id := range ids {
		targets[i] = DamageTarget{ID: id, Amount: e.Amount}
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}
