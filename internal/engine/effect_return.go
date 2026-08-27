package engine

import "fmt"

// ReturnToDeck puts each card its Target selects on top of its owner's deck.
//
// A card returned to a deck leaves play and loses all the state it built up while
// in play (damage, armor spent, Æmber on it); it becomes the top card of its
// owner's deck, to be drawn again later. When several cards return at once the
// controller chooses the order they are stacked.
type ReturnToDeck struct {
	Target Target
}

// Text renders the effect, e.g. "put each artifact on top of its owner's deck".
func (e ReturnToDeck) Text() string {
	return fmt.Sprintf("put %s on top of its owner's deck", e.Target.Text())
}

// Resolve moves each selected card from play to the top of its owner's deck.
func (e ReturnToDeck) Resolve(ctx *EffectContext) {
	ids := ctx.Resolver.OrderByChoice(ctx.Controller, "Choose the next card to put on top of the deck", e.Target.Select(ctx))
	for _, id := range ids {
		ctx.Resolver.ReturnToTopOfDeck(id)
	}
}

// ReturnToHand puts each card its Target selects into its owner's hand.
//
// Like returning to a deck, a card put into hand leaves play and loses the state
// it built up there (damage, Æmber on it, and so on); it can be played again from
// hand later. This is how a "Destroyed:" ability can save its own creature — the
// creature is moved to hand as it is destroyed, so it never reaches the discard.
type ReturnToHand struct {
	Target Target
}

// Text renders the effect, e.g. "put this creature into its owner's hand".
func (e ReturnToHand) Text() string {
	return fmt.Sprintf("put %s into its owner's hand", e.Target.Text())
}

// Resolve moves each selected card from play to its owner's hand.
func (e ReturnToHand) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.ReturnToHand(id)
	}
}
