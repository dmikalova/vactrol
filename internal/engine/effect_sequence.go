package engine

import "strings"

// Sequence resolves several effects in order, the way a card lists several
// sentences of rules text that happen one after another. Each child resolves
// fully before the next begins, and the rendered text joins them with ", and".
type Sequence struct {
	Effects []Effect
}

// Text joins the child effect texts.
func (e Sequence) Text() string {
	parts := make([]string, 0, len(e.Effects))
	for _, child := range e.Effects {
		parts = append(parts, child.Text())
	}
	return strings.Join(parts, ", and ")
}

// Resolve resolves each child effect in order.
func (e Sequence) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

// validate surfaces the first configuration error among the child effects.
func (e Sequence) validate() error {
	for _, child := range e.Effects {
		if err := validateEffect(child); err != nil {
			return err
		}
	}
	return nil
}
