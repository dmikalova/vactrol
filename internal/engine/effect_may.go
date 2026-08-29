package engine

// A "you may" effect is optional: it offers the controller the choice to resolve
// its inner effect or to decline it entirely. It models KeyForge's "You may <do
// X>", where passing is always allowed even when a legal target exists — the
// distinction that keeps Chuff Ape's "you may destroy another friendly creature"
// from ever being forced.
//
//rulebook:effect May
type May struct {
	Do Effect
}

// Text renders the effect, e.g. "you may destroy another friendly creature -> fully
// heal Chuff Ape".
func (e May) Text() string {
	return "you may " + e.Do.Text()
}

// Resolve asks the controller whether to resolve the inner effect, and does so only
// when they accept.
func (e May) Resolve(ctx *EffectContext) {
	if ctx.ChooseOption(capitalizeFirst(e.Text()), []string{"Yes", "No"}) == 0 {
		e.Do.Resolve(ctx)
	}
}

// validate descends into the wrapped effect.
func (e May) validate() error { return validateEffect(e.Do) }
