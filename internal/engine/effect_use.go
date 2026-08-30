package engine

import "fmt"

// UseFriendlyCardsOfHouse uses the controller's ready cards of a named house, one
// at a time. It reaches both creatures and artifacts: a creature is used by
// choosing whether it reaps, fights, or uses its Action: ability, while an artifact
// uses its Action: ability. The source can be excluded for "other" card text.
type UseFriendlyCardsOfHouse struct {
	House House
	Count int
	Other bool
}

// validate requires a real house and a positive count.
func (e UseFriendlyCardsOfHouse) validate() error {
	if e.House == HouseNone {
		return fmt.Errorf("UseFriendlyCardsOfHouse: house must be set")
	}
	if e.Count <= 0 {
		return fmt.Errorf("UseFriendlyCardsOfHouse: count must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "use 2 other Mars cards, one at a time".
func (e UseFriendlyCardsOfHouse) Text() string {
	count := fmt.Sprintf("%d", e.Count)
	noun := "cards"
	if e.Count == 1 {
		count = "a"
		noun = "card"
	}
	other := ""
	if e.Other {
		other = "other "
	}
	return "use " + count + " " + other + e.House.String() + " " + noun + ", one at a time"
}

// Resolve chooses and uses up to Count matching ready cards. Each use fully
// resolves before the next choice is offered, so an exhausted card drops out of
// later choices.
func (e UseFriendlyCardsOfHouse) Resolve(ctx *EffectContext) {
	for i := 0; i < e.Count; i++ {
		candidates := e.candidates(ctx)
		if len(candidates) == 0 {
			return
		}
		id, ok := ctx.ChooseCard("Choose a "+e.House.String()+" card to use", candidates)
		if !ok {
			return
		}
		e.use(ctx, id)
	}
}

// candidates returns the friendly, ready cards of the effect's house that can be
// used by an ability right now.
func (e UseFriendlyCardsOfHouse) candidates(ctx *EffectContext) []LocalID {
	ids := append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Artifacts(ctx.Controller)...)
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if e.Other && id == ctx.Source {
			continue
		}
		if ctx.Resolver.House(id) != e.House || ctx.Resolver.Exhausted(id) {
			continue
		}
		switch ctx.Resolver.TypeOf(id) {
		case Creature:
			out = append(out, id)
		case Artifact:
			if ctx.Resolver.HasAction(id) {
				out = append(out, id)
			}
		}
	}
	return out
}

// use resolves the chosen card's use.
func (e UseFriendlyCardsOfHouse) use(ctx *EffectContext, id LocalID) {
	if ctx.Resolver.TypeOf(id) == Artifact {
		ctx.Resolver.UseActionOf(id)
		return
	}
	UseVerb{}.Apply(ctx, id)
}

// Sentence resolves its child normally but renders it as a complete sentence, so a
// following Sequence child starts a new sentence instead of joining with ", and".
type Sentence struct {
	Effect Effect
}

// Text renders the child as a punctuated sentence.
func (e Sentence) Text() string { return punctuate(e.Effect.Text()) }

// Resolve resolves the child effect.
func (e Sentence) Resolve(ctx *EffectContext) { e.Effect.Resolve(ctx) }

// validate surfaces a configuration error from the child effect.
func (e Sentence) validate() error { return validateEffect(e.Effect) }
