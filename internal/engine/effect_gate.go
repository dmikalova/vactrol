package engine

// A result gate resolves one action and then a follow-up, but only when the first
// action actually happened — KeyForge's "A -> B" (destroy a creature -> steal 1
// Æmber; purge a creature -> give a +1 power counter). The follow-up never runs
// when the gate does nothing (no valid target, an empty zone, a declined choice).

// GatingEffect is an effect that can be the first half of a Then: besides
// resolving, it reports whether it did anything. The report is an unexported
// method, so only engine effects (Destroy, Purge) can be a gate.
type GatingEffect interface {
	Effect
	// resolveGate resolves the effect and reports whether it did anything.
	resolveGate(ctx *EffectContext) bool
}

// Then is the "A -> B" result gate: it resolves First and, only when First did
// something, resolves Result.
type Then struct {
	First  GatingEffect
	Result Effect
}

// Text renders the gate, e.g. "destroy a damaged creature -> steal 1 Æmber".
func (e Then) Text() string {
	return e.First.Text() + " -> " + e.Result.Text()
}

// Resolve runs First and then Result only if First did something.
func (e Then) Resolve(ctx *EffectContext) {
	if e.First.resolveGate(ctx) {
		e.Result.Resolve(ctx)
	}
}

// validate surfaces a configuration error in either half of the gate.
func (e Then) validate() error {
	if err := validateEffect(e.First); err != nil {
		return err
	}
	return validateEffect(e.Result)
}
