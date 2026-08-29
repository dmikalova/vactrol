package engine

import "fmt"

// A card returned to a deck leaves play and loses all the state it built up while
// in play — damage, spent armor, Æmber on it — and becomes the top card of its
// owner's deck, to be drawn again later. When several cards return at once the
// controller chooses the order they are stacked.
//
//rulebook:effect Return to Deck
type ReturnToDeck struct {
	Target Target
}

// Text renders the effect, e.g. "put each artifact on top of its owner's deck".
func (e ReturnToDeck) Text() string {
	return fmt.Sprintf("put %s on top of its owner's deck", e.Target.Text())
}

// Resolve moves each selected card from play to the top of its owner's deck.
func (e ReturnToDeck) Resolve(ctx *EffectContext) {
	ids := ctx.OrderByChoice("Choose the next card to put on top of the deck", e.Target.Select(ctx))
	for _, id := range ids {
		ctx.Resolver.ReturnToTopOfDeck(id)
	}
}

// Like returning to a deck, a card put into its owner's hand leaves play and
// loses the state it built up there — damage, Æmber on it, and so on — and can be
// played again later. This is how a "Destroyed:" ability can save its own
// creature: the creature is moved to hand as it is destroyed, so it never reaches
// the discard zone.
//
//rulebook:effect Return to Hand
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

// ReturnToArchives puts each card its Target selects into its owner's archives.
//
// Like returning to hand or deck, a card put into archives leaves play and sheds
// the state it built up there (damage, Æmber on it, upgrades). Archived cards
// return to their owner's hand when that player later chooses to take them. This
// is how a "Destroyed:" ability can bank its own creature instead of discarding it.
type ReturnToArchives struct {
	Target Target
}

// Text renders the effect, e.g. "put this creature into its owner's archives".
func (e ReturnToArchives) Text() string {
	return fmt.Sprintf("put %s into its owner's archives", e.Target.Text())
}

// Resolve moves each selected card from play into its owner's archives.
func (e ReturnToArchives) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.ReturnToArchives(id)
	}
}

// ReturnArtifactsToHand puts up to Max artifacts (either player's) into their
// owners' hands. The controller chooses them one at a time and may stop early,
// so it is "up to" rather than exactly Max.
type ReturnArtifactsToHand struct {
	Max int
}

// Text renders the effect, e.g. "put up to 3 artifacts into their owners' hands".
func (e ReturnArtifactsToHand) Text() string {
	return fmt.Sprintf("put up to %d artifacts into their owners' hands", e.Max)
}

// Resolve returns artifacts to hand one at a time, up to Max. Each step offers the
// artifacts in play plus a "Done" option to stop early; when no artifacts remain,
// "Done" is the only option and is chosen automatically.
func (e ReturnArtifactsToHand) Resolve(ctx *EffectContext) {
	const done = "Done"
	for i := 0; i < e.Max; i++ {
		cands := append(ctx.Resolver.Artifacts(ctx.Controller), ctx.Resolver.Artifacts(ctx.Opponent())...)
		options := make([]string, 0, len(cands)+1)
		for _, id := range cands {
			options = append(options, ctx.Resolver.Name(id))
		}
		options = append(options, done)
		choice := ctx.ChooseOption("Choose an artifact to return to hand", options)
		if choice >= len(cands) {
			return // "Done" (the last option), or an out-of-range choice
		}
		ctx.Resolver.ReturnToHand(cands[choice])
	}
}
