package engine

// This file holds playing a card from hand: creatures, artifacts, actions, and
// upgrades, plus discarding from hand and the shared checks (house match, taking
// a card of the right type, granting its Æmber pips).

// PlayCreature plays a creature from hand onto the battleline. flankLeft places it
// on the left flank; otherwise it goes to the right flank.
func (g *Game) PlayCreature(player, handIndex int, flankLeft bool) (LocalID, error) {
	if g.cannotPlayCreatures(player) {
		return 0, ErrCannotPlayCreature
	}
	id, err := g.takeFromHand(player, handIndex, Creature)
	if err != nil {
		return 0, err
	}
	core := &g.State.Cards[id]
	core.Exhausted = true // enters play exhausted; readies during the end-of-turn ready step
	core.ArmorRemaining = int16(g.armor(id))
	if flankLeft {
		g.State.Battleline[player].addFront(id)
	} else {
		g.State.Battleline[player].add(id)
	}
	g.logf("%s plays %s to the battleline", g.names[player], g.Name(id))
	g.applyAemberBonus(id)
	g.triggerAbilities(id, TriggerAfterPlay, 0, false)
	g.fireCreatureEnters(id)
	g.fireLasting(EventCreaturePlayed, player, id)
	return id, nil
}

// PlayArtifact plays an artifact from hand into the artifact row. If the opponent
// controls a card that tolls artifact plays (Customs Office), the player must pay
// that toll first and cannot play the artifact if they cannot.
func (g *Game) PlayArtifact(player, handIndex int) (LocalID, error) {
	id, err := g.validateHandPlay(player, handIndex, Artifact)
	if err != nil {
		return 0, err
	}
	if err := g.chargeToll(player, TollPlayArtifact); err != nil {
		return 0, err
	}
	g.State.Hand[player].removeAt(handIndex)
	g.State.CardsPlayedThisTurn[player]++
	g.State.Cards[id].Exhausted = true // enters play exhausted; readies during the end-of-turn ready step
	g.State.Artifacts[player].add(id)
	g.logf("%s plays artifact %s", g.names[player], g.Name(id))
	g.applyAemberBonus(id)
	g.triggerAbilities(id, TriggerAfterPlay, 0, false)
	g.fireArtifactPlayed(player, id)
	return id, nil
}

// chargeToll makes player pay every toll an opponent's in-play card imposes for
// action (Customs Office, Tentacus), moving that Æmber to the opponent. It is the
// single cost gate the play and use sites share: it changes nothing and returns
// ErrCannotPayToll when player cannot pay the full amount owed.
func (g *Game) chargeToll(player int, action TollAction) error {
	payee := 1 - player
	owed := 0
	for _, id := range g.allInPlay(payee) {
		if t := g.cat.def(id).Restricts.Toll; t.Amount > 0 && t.Action == action {
			owed += t.Amount
		}
	}
	if owed == 0 {
		return nil
	}
	if g.State.Aember[player] < owed {
		return ErrCannotPayToll
	}
	g.State.Aember[player] -= owed
	g.State.Aember[payee] += owed
	g.logf("%s pays %d Æmber to %s to %s", g.names[player], owed, g.names[payee], action.phrase())
	return nil
}

// PlayAction plays an action card: its Æmber bonus and "Play:" abilities resolve,
// then it goes to the discard pile.
func (g *Game) PlayAction(player, handIndex int) error {
	id, err := g.takeFromHand(player, handIndex, Action)
	if err != nil {
		return err
	}
	g.logf("%s plays action %s", g.names[player], g.Name(id))
	g.applyAemberBonus(id)
	g.triggerAbilities(id, TriggerAfterPlay, 0, false)
	g.State.Discard[player].add(id)
	return nil
}

// DiscardFromHand discards a card of the active house from a player's hand,
// moving it to the discard pile. It performs no other effect.
func (g *Game) DiscardFromHand(player, handIndex int) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	hand := &g.State.Hand[player]
	if handIndex < 0 || handIndex >= int(hand.Count) {
		return ErrCardNotInHand
	}
	id := hand.IDs[handIndex]
	if !g.inActiveHouse(g.cat.def(id)) {
		return ErrWrongHouse
	}
	hand.removeAt(handIndex)
	g.State.Discard[player].add(id)
	g.logf("%s discards %s", g.names[player], g.Name(id))
	return nil
}

// PlayUpgrade plays an upgrade from hand, attaching it to a chosen creature
// (friendly or enemy). It returns the host creature it was attached to.
func (g *Game) PlayUpgrade(player, handIndex int) (LocalID, error) {
	if g.State.Winner >= 0 {
		return 0, ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return 0, ErrNotActivePlayer
	}
	if g.cannotPlayCard(player) {
		return 0, ErrCardPlayLimit
	}
	hand := &g.State.Hand[player]
	if handIndex < 0 || handIndex >= int(hand.Count) {
		return 0, ErrCardNotInHand
	}
	id := hand.IDs[handIndex]
	def := g.cat.def(id)
	if def.Type != Upgrade {
		return 0, ErrWrongType
	}
	if !g.inActiveHouse(def) {
		return 0, ErrWrongHouse
	}
	candidates := append(g.battlelineCopy(player), g.battlelineCopy(1-player)...)
	host, ok := g.pickCreature(player, g.Name(id), "Choose a creature to upgrade", candidates)
	if !ok {
		return 0, ErrNoTarget
	}
	hand.removeAt(handIndex)
	g.State.CardsPlayedThisTurn[player]++
	g.applyAemberBonus(id)
	hostCore := &g.State.Cards[host]
	hostCore.Upgrades[hostCore.UpgradeCount] = id
	hostCore.UpgradeCount++
	g.logf("%s attaches %s to %s", g.names[player], g.Name(id), g.Name(host))
	g.fireUpgradePlay(host, def)
	return host, nil
}

// DiscardCardFromHand moves a specific card from a player's hand to their discard
// zone, with no active-player or house checks (an effect may discard from either
// hand). It does nothing if the card is not in that hand.
func (g *Game) DiscardCardFromHand(owner int, id LocalID) {
	hand := &g.State.Hand[owner]
	i := hand.indexOf(id)
	if i < 0 {
		return
	}
	hand.removeAt(i)
	g.State.Discard[owner].add(id)
	g.logf("%s discards %s", g.names[owner], g.Name(id))
}

// inActiveHouse reports whether a card of the given definition matches the
// active house for the purpose of PLAYING or discarding it from hand: true when
// no house has been chosen or the card's own house is the active house. Manual
// (sandbox) mode lifts the restriction. Versatile does not apply here — it only
// relaxes using a card already in play (see usableInActiveHouse).
func (g *Game) inActiveHouse(def *CardDefinition) bool {
	return g.manual || g.State.ActiveHouse == HouseNone || def.House == g.State.ActiveHouse
}

// takeFromHand validates and removes a card of the wanted type from a hand.
func (g *Game) takeFromHand(player, handIndex int, want CardType) (LocalID, error) {
	id, err := g.validateHandPlay(player, handIndex, want)
	if err != nil {
		return 0, err
	}
	g.State.Hand[player].removeAt(handIndex)
	g.State.CardsPlayedThisTurn[player]++
	return id, nil
}

// validateHandPlay runs the read-only checks that a hand card may be played —
// game not over, active player, no card-play limit reached, index in hand, right
// type, active house — and returns the card without mutating the hand. It is the
// shared validation step for takeFromHand and for plays that must charge a cost
// (PlayArtifact) only once the play is otherwise legal.
func (g *Game) validateHandPlay(player, handIndex int, want CardType) (LocalID, error) {
	if g.State.Winner >= 0 {
		return 0, ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return 0, ErrNotActivePlayer
	}
	if g.cannotPlayCard(player) {
		return 0, ErrCardPlayLimit
	}
	hand := &g.State.Hand[player]
	if handIndex < 0 || handIndex >= int(hand.Count) {
		return 0, ErrCardNotInHand
	}
	id := hand.IDs[handIndex]
	def := g.cat.def(id)
	if def.Type != want {
		return 0, ErrWrongType
	}
	if !g.inActiveHouse(def) {
		return 0, ErrWrongHouse
	}
	return id, nil
}

// CanPlay reports whether the player can play the given hand card right now: nil
// if playable, otherwise the reason (wrong house, card-play limit, creatures
// barred, or an upgrade with no host). It mirrors the checks the PlayX methods
// enforce, so a UI can dim unplayable cards and explain why before the click.
func (g *Game) CanPlay(player int, id LocalID) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	def := g.cat.def(id)
	if def.Type == Creature && g.cannotPlayCreatures(player) {
		return ErrCannotPlayCreature
	}
	if g.cannotPlayCard(player) {
		return ErrCardPlayLimit
	}
	if !g.inActiveHouse(def) {
		return ErrWrongHouse
	}
	if def.Type == Upgrade &&
		len(g.State.Battleline[player].slice()) == 0 && len(g.State.Battleline[1-player].slice()) == 0 {
		return ErrNoTarget
	}
	return nil
}

// applyAemberBonus grants a card's Æmber pips to its controller.
func (g *Game) applyAemberBonus(id LocalID) {
	def := g.cat.def(id)
	if def.AemberBonus > 0 {
		o := g.owner(id)
		g.State.Aember[o] += def.AemberBonus
		g.logf("%s gains %d Æmber from %s", g.names[o], def.AemberBonus, def.Name)
	}
}
