package engine

// ShuffleDiscard shuffles the controller's discard pile into their deck (Help from
// Future Self, after tutoring).
//
//rulebook:effect Shuffle Discard
type ShuffleDiscard struct{}

// Text renders the effect, "shuffle your discard pile into your deck".
func (ShuffleDiscard) Text() string {
	return "shuffle your discard pile into your deck"
}

// Resolve moves the controller's discard pile into their deck and shuffles it.
func (ShuffleDiscard) Resolve(ctx *EffectContext) {
	ctx.Resolver.ShuffleDiscardIntoDeck(ctx.Controller)
}
