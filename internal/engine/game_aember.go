package engine

// gainAember adds Æmber from the common supply to a player's pool. It is the
// single seam for pool gains: before the Æmber reaches the pool, continuous
// replacement effects such as Ether Spider may capture it instead. Movements that
// already come from another pool or card (steal, capture, returning captured
// Æmber) are not gains and intentionally bypass this helper.
func (g *Game) gainAember(player, amount int) (LocalID, bool) {
	if capturer, ok := g.opponentAemberCaptor(player); ok {
		g.State.Cards[capturer].Amber += int16(amount)
		return capturer, true
	}
	g.State.Aember[player] += amount
	return 0, false
}

// GainAember is the Resolver entry point for gainAember.
func (g *Game) GainAember(player, amount int) (LocalID, bool) { return g.gainAember(player, amount) }

// opponentAemberCaptor returns the first opposing in-play creature whose static
// replacement captures Æmber that would be added to player's pool.
func (g *Game) opponentAemberCaptor(player int) (LocalID, bool) {
	opponent := 1 - player
	for _, id := range g.allInPlay(opponent) {
		def := g.cat.def(id)
		if def.Type == Creature && def.CapturesOpponentAember {
			return id, true
		}
	}
	return 0, false
}
