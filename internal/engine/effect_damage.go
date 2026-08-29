package engine

import "fmt"

// Dealing damage places that many damage tokens on each creature the effect
// targets; armor reduces each separate instance of damage before it lands. A
// creature whose damage reaches or exceeds its power is destroyed. When one
// ability deals damage to several creatures they are damaged simultaneously and
// any that died are destroyed together, so no creature's destruction changes
// another's.
//
//rulebook:effect Deal Damage
type DealDamage struct {
	Amount int
	Per    Count
	Target Target
}

// Text renders the effect, e.g. "deal 2 damage to each enemy creature". A "for
// each" count leads the sentence (rule 9), e.g. "for each friendly creature in
// play, deal 1 damage to a creature".
func (e DealDamage) Text() string {
	return forEach(e.Per, fmt.Sprintf("deal %d damage to %s", e.Amount, e.Target.Text()))
}

// Resolve deals the damage to every selected creature simultaneously, resolving
// destruction as part of it. A Per count multiplies the amount dealt.
func (e DealDamage) Resolve(ctx *EffectContext) {
	amount := e.Amount
	if e.Per != nil {
		amount *= e.Per.Value(ctx)
	}
	ids := e.Target.Select(ctx)
	targets := make([]DamageTarget, len(ids))
	for i, id := range ids {
		targets[i] = DamageTarget{ID: id, Amount: amount}
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}

// SplashDamage deals damage to one chosen creature that is not on a flank and a
// smaller "splash" amount to each of that creature's neighbors, all at once. The
// not-on-a-flank restriction guarantees the chosen creature has two neighbors.
type SplashDamage struct {
	Amount int
	Splash int
}

// Text renders the effect.
func (e SplashDamage) Text() string {
	return fmt.Sprintf("deal %d damage to a creature that is not on a flank and %d damage to each of its neighbors", e.Amount, e.Splash)
}

// Resolve chooses a non-flank creature, then damages it and its neighbors as one
// simultaneous batch.
func (e SplashDamage) Resolve(ctx *EffectContext) {
	chosen := Target{Kind: TargetChosenCreature}.NotOnFlank().Select(ctx)
	if len(chosen) == 0 {
		return
	}
	targets := []DamageTarget{{ID: chosen[0], Amount: e.Amount}}
	for _, n := range neighbors(ctx, chosen[0]) {
		targets = append(targets, DamageTarget{ID: n, Amount: e.Splash})
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}
