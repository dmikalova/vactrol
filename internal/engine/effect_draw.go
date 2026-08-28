package engine

import "fmt"

// Drawing puts the top card of your deck into your hand. If your deck is empty
// when you must draw, your discard pile is shuffled to form a new deck first, so
// you only fail to draw when both deck and discard are empty.
//
//rulebook:effect Draw
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
