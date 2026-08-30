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

// validate requires an explicit target.
func (e Exhaust) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Exhaust")
	}
	return nil
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

// Readying a creature turns it upright again so it can be used, the opposite of
// exhausting. It readies each creature the effect targets.
//
//rulebook:effect Ready
type Ready struct {
	Target Target
}

// validate requires an explicit target.
func (e Ready) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Ready")
	}
	return nil
}

func (e Ready) verb() string       { return "ready" }
func (e Ready) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "ready this creature".
func (e Ready) Text() string { return e.verb() + " " + e.targetText() }

// Resolve readies each selected creature.
func (e Ready) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetExhausted(id, false)
	}
}
