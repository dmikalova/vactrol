package engine

// Stun and Unstun apply and remove the stun status.

// A stun is a status placed on a creature. A stunned creature must shake off the
// stun before it can do anything else: the next time it is used to reap, fight,
// or use an "Action:" ability, it is exhausted and the stun is removed instead of
// that action happening. Stunning applies this status to each creature the effect
// targets.
//
//rulebook:effect Stun
type Stun struct {
	Target Target
}

// Text renders the effect, e.g. "stun each friendly creature".
func (e Stun) Text() string { return "stun " + e.Target.Text() }

// Resolve stuns each selected creature.
func (e Stun) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetStunned(id, true)
	}
}

// Unstunning a creature removes the stun status from each creature the effect
// targets, freeing it to act normally instead of having to shake the stun off.
//
//rulebook:effect Unstun
type Unstun struct {
	Target Target
}

// Text renders the effect, e.g. "unstun each friendly creature".
func (e Unstun) Text() string { return "unstun " + e.Target.Text() }

// Resolve clears the stun on each selected creature.
func (e Unstun) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetStunned(id, false)
	}
}
