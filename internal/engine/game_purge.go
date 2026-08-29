package engine

// This file holds the purge zone: cards set aside out of the game. Purging pulls
// a card from a discard zone into its owner's purge zone, from which nothing in
// the base game returns — so purge is the game's way of permanently answering
// recursion out of the discard.

// purgeFromDiscard moves a card from a player's discard zone to their purge zone.
// Callers pass a card already in that discard zone.
func (g *Game) purgeFromDiscard(owner int, id LocalID) {
	g.State.Discard[owner].remove(id)
	g.State.Purge[owner].add(id)
	g.logf("%s purges %s from a discard zone", g.names[g.State.ActivePlayer], g.Name(id))
}

// purgeFromHand moves a card from a player's hand to their purge zone. Callers
// pass a card already in that hand.
func (g *Game) purgeFromHand(owner int, id LocalID) {
	g.State.Hand[owner].remove(id)
	g.State.Purge[owner].add(id)
	g.logf("%s purges %s from a hand", g.names[g.State.ActivePlayer], g.Name(id))
}
