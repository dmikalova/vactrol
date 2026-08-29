package engine

import "strings"

// Sequence resolves several effects in order, the way a card lists several
// sentences of rules text that happen one after another. Each child resolves
// fully before the next begins, and the rendered text joins them with ", and".
type Sequence struct {
	Effects []Effect
}

// combinable is a plain "verb the target" effect (e.g. "stun this creature")
// whose text can be folded together when several of them in a row act on the same
// target, so a Sequence reads "stun and exhaust this creature" instead of the
// clumsier "stun this creature, and exhaust this creature".
type combinable interface {
	verb() string
	targetText() string
}

// Text joins the child effect texts, folding a run of combinable effects that
// share a target into one "verb and verb ... target" phrase.
func (e Sequence) Text() string {
	parts := make([]string, 0, len(e.Effects))
	for i := 0; i < len(e.Effects); {
		c, ok := e.Effects[i].(combinable)
		if !ok {
			parts = append(parts, e.Effects[i].Text())
			i++
			continue
		}
		// Gather the run of following combinables that share this target.
		target := c.targetText()
		verbs := []string{c.verb()}
		i++
		for i < len(e.Effects) {
			next, ok := e.Effects[i].(combinable)
			if !ok || next.targetText() != target {
				break
			}
			verbs = append(verbs, next.verb())
			i++
		}
		parts = append(parts, strings.Join(verbs, " and ")+" "+target)
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
