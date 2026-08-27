package engine

// Stun and Unstun apply and remove the stun status.
//
// A stun is a status placed on a creature. A stunned creature must shake off the
// stun before it can do anything else: the next time it is used to reap, fight,
// or use an "Action:" ability, it is exhausted and the stun is removed instead of
// that action happening (see Game.recoverFromStun). Removing a stun by other
// means (unstunning) frees the creature to act normally.

// Stun stuns each creature the Target selects.
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

// Unstun removes the stun from each creature the Target selects.
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
