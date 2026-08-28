package engine

// ReturnFromDiscard returns a card the controller chooses from their own discard
// pile — to their hand, or (ToDeck) to the top of their deck. CreaturesOnly
// restricts the choice to creatures. This is how cards recur from the discard
// pile, e.g. "Put a creature from your discard pile on top of your deck."
type ReturnFromDiscard struct {
	CreaturesOnly bool
	ToDeck        bool
}

// noun renders the kind of card the effect returns.
func (e ReturnFromDiscard) noun() string {
	if e.CreaturesOnly {
		return "creature"
	}
	return "card"
}

// Text renders the effect, e.g. "put a card from your discard pile into your
// hand" or "put a creature from your discard pile on top of your deck".
func (e ReturnFromDiscard) Text() string {
	dest := "into your hand"
	if e.ToDeck {
		dest = "on top of your deck"
	}
	return "put a " + e.noun() + " from your discard pile " + dest
}

// Resolve lets the controller choose an eligible card from their discard pile and
// moves it to the chosen destination. It does nothing if there is no candidate or
// the choice is declined.
func (e ReturnFromDiscard) Resolve(ctx *EffectContext) {
	discard := ctx.Resolver.Discard(ctx.Controller)
	candidates := discard
	if e.CreaturesOnly {
		candidates = nil
		for _, id := range discard {
			if ctx.Resolver.IsCreature(id) {
				candidates = append(candidates, id)
			}
		}
	}
	id, ok := ctx.Resolver.ChooseCreature(ctx.Controller, "Choose a "+e.noun()+" from your discard pile", candidates)
	if !ok {
		return
	}
	if e.ToDeck {
		ctx.Resolver.ReturnFromDiscardToTopOfDeck(id)
	} else {
		ctx.Resolver.ReturnFromDiscardToHand(id)
	}
}
