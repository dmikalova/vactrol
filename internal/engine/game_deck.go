package engine

// This file holds the deck: drawing cards (recycling the discard when the draw
// zone runs out) and shuffling.

// drawOne draws a single card to the player's hand. If the deck is empty, the
// discard pile is shuffled to form a new deck first (KeyForge recycles the
// discard when the draw zone runs out). It returns false only when both the deck
// and discard are empty, so nothing can be drawn.
func (g *Game) drawOne(player int) bool {
	deck := &g.State.Deck[player]
	if deck.Count == 0 {
		discard := &g.State.Discard[player]
		if discard.Count == 0 {
			return false
		}
		for discard.Count > 0 {
			deck.add(discard.removeAt(0))
		}
		g.Shuffle(player)
	}
	g.State.Hand[player].add(deck.removeAt(0))
	return true
}

// draw draws count cards into the player's hand, stopping early only when the
// deck and discard are both exhausted.
func (g *Game) draw(player, count int) {
	for i := 0; i < count; i++ {
		if !g.drawOne(player) {
			break
		}
	}
}

// drawTo draws until the hand holds n cards (or nothing is left to draw).
func (g *Game) drawTo(player, n int) {
	for int(g.State.Hand[player].Count) < n {
		if !g.drawOne(player) {
			break
		}
	}
}

// canDraw reports whether a player still has a card to draw — one in their deck, or
// one in their discard pile that would be reshuffled in.
func (g *Game) canDraw(player int) bool {
	return g.State.Deck[player].Count > 0 || g.State.Discard[player].Count > 0
}

// Shuffle randomizes a player's deck using the game's seeded RNG.
func (g *Game) Shuffle(player int) {
	d := &g.State.Deck[player]
	for i := int(d.Count) - 1; i > 0; i-- {
		j := g.rng.Intn(i + 1)
		d.IDs[i], d.IDs[j] = d.IDs[j], d.IDs[i]
	}
}
