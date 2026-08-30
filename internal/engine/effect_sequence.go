package engine

import "strings"

// Sequence resolves several effects in order, the way a card lists several
// sentences of rules text that happen one after another. Each child resolves
// fully before the next begins, and the rendered text joins them with ", and".
type Sequence struct {
	Effects []Effect
}

// combinable is a plain "verb the target" effect (e.g. "stun this creature")
// whose text can be folded together with its neighbours in a Sequence. A run of
// combinables folds along whichever axis they share: neighbours with the same
// target fold their verbs ("stun and exhaust this creature"), and neighbours with
// the same verb fold their targets ("destroy an enemy creature and a friendly
// creature"). Either way the Sequence reads as one phrase instead of the clumsier
// "stun this creature, and exhaust this creature".
type combinable interface {
	verb() string
	targetText() string
}

// Text joins the child effect texts, folding each run of combinables that shares
// a verb or a target into a single "verb and verb ... target" or "verb target and
// target ..." phrase.
func (e Sequence) Text() string {
	parts := make([]string, 0, len(e.Effects))
	for i := 0; i < len(e.Effects); {
		c, ok := e.Effects[i].(combinable)
		if !ok {
			parts = append(parts, e.Effects[i].Text())
			i++
			continue
		}
		phrase, next := foldCombinable(e.Effects, i, c)
		parts = append(parts, phrase)
		i = next
	}
	return joinSequenceParts(parts)
}

// joinSequenceParts preserves the usual ", and" flow for phrase effects, but lets
// a child effect that already rendered a full sentence continue with the next
// effect as a new sentence.
func joinSequenceParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	text := parts[0]
	for _, part := range parts[1:] {
		if strings.HasSuffix(text, ".") {
			text += " " + capitalizeFirst(part)
			continue
		}
		text += ", and " + part
	}
	return text
}

// foldCombinable folds the run of combinables starting at i into one phrase and
// reports the index just past the run. The fold axis is chosen from the first
// neighbour: a shared target folds the verbs, a shared verb folds the targets; a
// lone combinable renders as "verb target".
func foldCombinable(effects []Effect, i int, c combinable) (string, int) {
	verb, target := c.verb(), c.targetText()
	next, ok := peekCombinable(effects, i+1)
	switch {
	case ok && next.targetText() == target:
		verbs := []string{verb}
		i++
		for ; ; i++ {
			n, ok := peekCombinable(effects, i)
			if !ok || n.targetText() != target {
				break
			}
			verbs = append(verbs, n.verb())
		}
		return strings.Join(verbs, " and ") + " " + target, i
	case ok && next.verb() == verb:
		targets := []string{target}
		i++
		for ; ; i++ {
			n, ok := peekCombinable(effects, i)
			if !ok || n.verb() != verb {
				break
			}
			targets = append(targets, n.targetText())
		}
		return verb + " " + strings.Join(targets, " and "), i
	default:
		return verb + " " + target, i + 1
	}
}

// peekCombinable reports the effect at i as a combinable, if it is one and in range.
func peekCombinable(effects []Effect, i int) (combinable, bool) {
	if i >= len(effects) {
		return nil, false
	}
	c, ok := effects[i].(combinable)
	return c, ok
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
