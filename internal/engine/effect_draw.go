package engine

import "fmt"

// Draw makes the controller draw cards from the top of their deck into their
// hand.
//
// Drawing puts the top card of the deck into hand. If the deck is empty when a
// card must be drawn, the discard pile is shuffled to form a new deck first, so a
// player only fails to draw when both deck and discard are empty (see
// Game.drawOne).
type Draw struct {
	Amount int
}

// Text renders the effect, e.g. "draw a card" or "draw 2 cards".
func (e Draw) Text() string {
	if e.Amount == 1 {
		return "draw a card"
	}
	return fmt.Sprintf("draw %d cards", e.Amount)
}

// Resolve draws the cards.
func (e Draw) Resolve(ctx *EffectContext) {
	ctx.Resolver.Draw(ctx.Controller, e.Amount)
}
