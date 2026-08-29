package engine

// Exhaust turns a creature sideways so it cannot be used again until it readies.
// Like Stun and Unstun it is a simple "verb the target" effect, so it folds
// together with a neighbouring status change on the same target when they share a
// Sequence (see combinable) — e.g. "stun and exhaust this creature".

// Exhausting a creature turns it sideways so it cannot be used again until it
// readies at the end of its controller's turn. It exhausts each creature the
// effect targets.
//
//rulebook:effect Exhaust
type Exhaust struct {
	Target Target
}

func (e Exhaust) verb() string       { return "exhaust" }
func (e Exhaust) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "exhaust this creature".
func (e Exhaust) Text() string { return e.verb() + " " + e.targetText() }

// Resolve exhausts each selected creature.
func (e Exhaust) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetExhausted(id, true)
	}
}
