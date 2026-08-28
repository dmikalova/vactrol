package engine

import "errors"

// Engine action errors.
var (
	ErrNotActivePlayer = errors.New("not the active player")
	ErrCardNotInHand   = errors.New("card index is not in hand")
	ErrWrongType       = errors.New("card is the wrong type for this action")
	ErrWrongHouse      = errors.New("card's house is not the active house")
	ErrCardExhausted   = errors.New("card is exhausted")
	ErrNoTarget        = errors.New("no legal target")
	ErrGameOver        = errors.New("the game is already over")
)

// BeginTurn starts a player's turn: it becomes the active player, the active
// house is cleared, and the forge step runs (keys are forged if affordable).
func (g *Game) BeginTurn(player int) {
	if g.State.Winner >= 0 {
		return
	}
	g.State.ActivePlayer = player
	g.State.ActiveHouse = HouseNone
	g.State.Turn++
	g.logf("--- %s begins turn %d ---", g.names[player], g.State.Turn)
	g.forgeKeys(player)
}

// ChooseHouse sets the active house for the current turn.
func (g *Game) ChooseHouse(player int, house House) error {
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	g.State.ActiveHouse = house
	g.logf("%s chooses house %s", g.names[player], house)
	g.offerArchives(player)
	return nil
}

// EndTurn readies the active player's cards (creatures and artifacts enter play
// exhausted and become ready in this step), refreshes armor, and draws back up to
// the hand size.
func (g *Game) EndTurn(player int) {
	for _, id := range g.allInPlay(player) {
		core := &g.State.Cards[id]
		core.Exhausted = false
		if g.cat.def(id).Type == Creature {
			core.ArmorRemaining = int16(g.armor(id))
		}
	}
	g.drawTo(player, HandSize)
	g.logf("%s ends their turn", g.names[player])
}

// drawOne draws a single card to the player's hand. If the deck is empty, the
// discard pile is shuffled to form a new deck first (KeyForge recycles the
// discard when the draw pile runs out). It returns false only when both the deck
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

// forgeKeys forges as many keys as the player can afford at the start of their
// turn, one at a time.
func (g *Game) forgeKeys(player int) {
	for g.State.Winner < 0 && g.State.Aember[player] >= KeyCost {
		g.forgeOneKey(player)
	}
}

// forgeOneKey forges a single key for a player if they can afford the current key
// cost, paying it and firing "after you forge a key" abilities. Forging the final
// key wins the game.
func (g *Game) forgeOneKey(player int) {
	if g.State.Aember[player] < KeyCost {
		return
	}
	g.State.Aember[player] -= KeyCost
	g.State.Keys[player]++
	g.logf("%s forges a key (%d/%d)", g.names[player], g.State.Keys[player], KeysToWin)
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerAfterForgeKey, 0, false)
	}
	if g.State.Keys[player] >= KeysToWin {
		g.State.Winner = player
		g.logf("%s wins the game!", g.names[player])
	}
}

// offerArchives asks a player — after they have chosen their house — whether to
// take their archived cards into their hand, moving the archives to hand if they
// accept. A player with no archived cards is not prompted.
func (g *Game) offerArchives(player int) {
	arc := &g.State.Archives[player]
	if arc.Count == 0 {
		return
	}
	if g.ChooseOption(player, "Take your archived cards into your hand?", []string{"Take them", "Leave them archived"}) != 0 {
		return
	}
	n := arc.Count
	for _, id := range arc.slice() {
		g.State.Hand[player].add(id)
	}
	*arc = Zone{}
	g.logf("%s takes %d card(s) from their archives into hand", g.names[player], n)
}

// archiveFromHand moves a card from a player's hand to their archives.
func (g *Game) archiveFromHand(player int, id LocalID) {
	if g.State.Hand[player].remove(id) {
		g.State.Archives[player].add(id)
		g.logf("%s archives a card", g.names[player])
	}
}

// archiveTopOfDeck moves the top card of a player's deck to their archives,
// reporting whether a card was available to archive.
func (g *Game) archiveTopOfDeck(player int) bool {
	deck := &g.State.Deck[player]
	if deck.Count == 0 {
		return false
	}
	id := deck.removeAt(0)
	g.State.Archives[player].add(id)
	g.logf("%s archives the top card of their deck", g.names[player])
	return true
}

// discardArchives moves all of a player's archived cards to their discard pile.
// The active player performs the discard, so they choose the order when it is
// their own archives but cannot when it is an opponent's — those enter the
// discard in a random order, since the active player cannot see them.
func (g *Game) discardArchives(owner int) {
	arc := &g.State.Archives[owner]
	if arc.Count == 0 {
		return
	}
	ids := cloneIDs(arc.slice())
	if owner == g.State.ActivePlayer {
		ids = g.orderByChoice(owner, "Choose the order to discard your archives", ids)
	} else {
		g.rng.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	}
	*arc = Zone{}
	for _, id := range ids {
		g.State.Discard[owner].add(id)
	}
	g.logf("%s discards %d archived card(s)", g.names[owner], len(ids))
}

// PlayCreature plays a creature from hand onto the battleline. flankLeft places it
// on the left flank; otherwise it goes to the right flank.
func (g *Game) PlayCreature(player, handIndex int, flankLeft bool) (LocalID, error) {
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
	return id, nil
}

// PlayArtifact plays an artifact from hand into the artifact row.
func (g *Game) PlayArtifact(player, handIndex int) (LocalID, error) {
	id, err := g.takeFromHand(player, handIndex, Artifact)
	if err != nil {
		return 0, err
	}
	g.State.Cards[id].Exhausted = true // enters play exhausted; readies during the end-of-turn ready step
	g.State.Artifacts[player].add(id)
	g.logf("%s plays artifact %s", g.names[player], g.Name(id))
	g.applyAemberBonus(id)
	g.triggerAbilities(id, TriggerAfterPlay, 0, false)
	g.fireArtifactPlayed(player, id)
	return id, nil
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
	host, ok := g.chooserFor(player).ChooseCreature("Choose a creature to upgrade", candidates)
	if !ok {
		return 0, ErrNoTarget
	}
	hand.removeAt(handIndex)
	g.applyAemberBonus(id)
	hostCore := &g.State.Cards[host]
	hostCore.Upgrades[hostCore.UpgradeCount] = id
	hostCore.UpgradeCount++
	g.logf("%s attaches %s to %s", g.names[player], g.Name(id), g.Name(host))
	return host, nil
}

// Reap uses a creature to reap, gaining 1 Æmber and firing "Reap:" abilities.
func (g *Game) Reap(player int, id LocalID) error {
	if err := g.canUse(player, id); err != nil {
		return err
	}
	g.reapWith(id)
	return nil
}

// reapWith performs a reap driven by a rule or ability, with no active-player or
// active-house checks: a stunned creature recovers instead; otherwise it
// exhausts, its controller gains 1 Æmber, and its "Reap:" abilities fire.
func (g *Game) reapWith(id LocalID) {
	if g.recoverFromStun(id) {
		return
	}
	p := g.owner(id)
	g.State.Cards[id].Exhausted = true
	g.State.Aember[p]++
	g.logf("%s reaps with %s (+1 Æmber)", g.names[p], g.Name(id))
	g.triggerAbilities(id, TriggerAfterReap, 0, false)
}

// UseAction uses a creature's or artifact's "Action:" ability.
func (g *Game) UseAction(player int, id LocalID) error {
	if err := g.usable(player, id); err != nil {
		return err
	}
	if !g.cat.def(id).hasTrigger(TriggerAction) {
		return ErrWrongType
	}
	g.useActionOf(id)
	return nil
}

// useActionOf fires a card's "Action:" ability, driven by a rule or ability: a
// stunned card recovers instead; otherwise it exhausts and its "Action:" fires.
func (g *Game) useActionOf(id LocalID) {
	if g.recoverFromStun(id) {
		return
	}
	g.State.Cards[id].Exhausted = true
	g.logf("%s uses %s's action ability", g.names[g.owner(id)], g.Name(id))
	g.triggerAbilities(id, TriggerAction, 0, false)
}

// Fight uses attacker to fight the enemy creature defender.
func (g *Game) Fight(player int, attacker, defender LocalID) error {
	if err := g.canUse(player, attacker); err != nil {
		return err
	}
	if g.cat.def(defender).Type != Creature || g.owner(defender) == player || !g.inPlay(defender) {
		return ErrNoTarget
	}
	if g.recoverFromStun(attacker) {
		return nil
	}
	g.fight(attacker, defender)
	return nil
}

// inActiveHouse reports whether a card of the given definition matches the
// active house for the purpose of PLAYING or discarding it from hand: true when
// no house has been chosen or the card's own house is the active house. Versatile
// does not apply here — it only relaxes using a card already in play (see
// usableInActiveHouse).
func (g *Game) inActiveHouse(def *CardDefinition) bool {
	return g.State.ActiveHouse == HouseNone || def.House == g.State.ActiveHouse
}

// usableInActiveHouse reports whether the in-play card id may be USED (reaped,
// fought, or its action ability activated) under the active house. It holds when
// the card is in the active house or is Versatile — printed on it or granted by
// an attached upgrade — letting it be used as if it belonged to the active house.
// Versatile only relaxes using: a Versatile card is still played from hand only
// when its own house is active.
func (g *Game) usableInActiveHouse(id LocalID) bool {
	return g.inActiveHouse(g.cat.def(id)) || g.hasKeyword(id, Versatile)
}

// usable runs the checks shared by reaping, fighting, and using an action
// ability: the card must be owned by the active player, in play, unexhausted, and
// usable under the active house. It does not restrict by card type — callers add
// that.
func (g *Game) usable(player int, id LocalID) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	if g.owner(id) != player || !g.inPlay(id) {
		return ErrWrongType
	}
	if g.State.Cards[id].Exhausted {
		return ErrCardExhausted
	}
	if !g.usableInActiveHouse(id) {
		return ErrWrongHouse
	}
	return nil
}

// canUse validates that a creature may be used to reap or fight right now.
func (g *Game) canUse(player int, id LocalID) error {
	if err := g.usable(player, id); err != nil {
		return err
	}
	if g.cat.def(id).Type != Creature {
		return ErrWrongType
	}
	return nil
}

// CanUse reports whether a creature may currently be used (reap/fight/action) by
// the player: nil if usable, otherwise the reason. UIs can call it to reject an
// action before prompting for a target.
func (g *Game) CanUse(player int, id LocalID) error { return g.canUse(player, id) }

// recoverFromStun handles using a stunned creature: it exhausts and clears the
// stun instead of performing the reap/fight/action. It reports whether the
// creature was stunned, in which case the caller should stop (the use is spent
// removing the stun and nothing else happens).
func (g *Game) recoverFromStun(id LocalID) bool {
	core := &g.State.Cards[id]
	if !core.Stunned {
		return false
	}
	core.Stunned = false
	core.Exhausted = true
	g.logf("%s recovers from stun instead of acting", g.Name(id))
	return true
}

// readyToUse reports whether a creature may be used by an ability right now. A
// creature can only be used while ready: an ability may still target an exhausted
// creature, but using it then does nothing. Unlike canUse this ignores the active
// player and active house — ability-driven use cares only about readiness.
func (g *Game) readyToUse(id LocalID) bool {
	if g.State.Cards[id].Exhausted {
		g.logf("%s is exhausted and cannot be used", g.Name(id))
		return false
	}
	return true
}

// takeFromHand validates and removes a card of the wanted type from a hand.
func (g *Game) takeFromHand(player, handIndex int, want CardType) (LocalID, error) {
	if g.State.Winner >= 0 {
		return 0, ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return 0, ErrNotActivePlayer
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
	hand.removeAt(handIndex)
	return id, nil
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

// fight resolves combat between an attacker and a defender. Both deal damage equal
// to their power simultaneously; Skirmish prevents the attacker taking damage back.
func (g *Game) fight(attacker, defender LocalID) {
	// Using a creature to fight exhausts it before anything else resolves, so a
	// "Before Fight" ability already sees the attacker exhausted.
	g.State.Cards[attacker].Exhausted = true
	g.triggerAbilities(attacker, TriggerBeforeFight, 0, false)

	// Assault and Hazardous deal their damage before fight damage: the attacker's
	// Assault hits the defender, the defender's Hazardous hits the attacker.
	// Either can destroy a fighter before combat. Skipped if a "Before Fight"
	// effect already removed the defender.
	if g.inPlay(defender) {
		var pre []DamageTarget
		if a := g.assault(attacker); a > 0 {
			pre = append(pre, DamageTarget{ID: defender, Amount: a})
		}
		if h := g.hazardous(defender); h > 0 {
			pre = append(pre, DamageTarget{ID: attacker, Amount: h})
		}
		if len(pre) > 0 {
			g.dealDamage(g.owner(attacker), pre...)
		}
	}

	// Combat damage is exchanged only while both fighters are still in play (a
	// "Before Fight" effect, Assault, or Hazardous can remove one first).
	if g.inPlay(attacker) && g.inPlay(defender) {
		ap, dp := g.Power(attacker), g.Power(defender)
		g.logf("%s (%d power) fights %s (%d power)", g.Name(attacker), ap, g.Name(defender), dp)
		// Both fighters take their damage simultaneously, then destruction is
		// resolved together as part of dealing it — so each dying creature is
		// already in the discard before the other's "Destroyed:" ability (or the
		// attacker's "After Fight") fires, and neither death changes the other's.
		targets := []DamageTarget{{ID: defender, Amount: ap}}
		if !g.hasKeyword(attacker, Skirmish) {
			targets = append(targets, DamageTarget{ID: attacker, Amount: dp})
		}
		g.dealDamage(g.owner(attacker), targets...)
	}
	g.triggerAbilities(attacker, TriggerAfterFight, 0, false)

	// "After a creature is destroyed fighting X": when exactly one combatant is
	// removed by the fight, the survivor's ability fires with the destroyed
	// creature as `it`.
	attackerDead, defenderDead := !g.inPlay(attacker), !g.inPlay(defender)
	if defenderDead && !attackerDead {
		g.triggerAbilities(attacker, TriggerAfterDestroyedFighting, defender, true)
	}
	if attackerDead && !defenderDead {
		g.triggerAbilities(defender, TriggerAfterDestroyedFighting, attacker, true)
	}
}

// applyRawDamage records damage on a creature, letting armor absorb it first. It
// only updates the damage counters; destruction is resolved by dealDamage, of
// which this is the per-creature step.
func (g *Game) applyRawDamage(id LocalID, amount int) {
	if amount <= 0 {
		return
	}
	core := &g.State.Cards[id]
	absorbed := min(int(core.ArmorRemaining), amount)
	if absorbed > 0 {
		core.ArmorRemaining -= int16(absorbed)
		amount -= absorbed
		g.logf("%s's armor absorbs %d damage", g.Name(id), absorbed)
	}
	if amount > 0 {
		core.Damage += int16(amount)
		g.logf("%s takes %d damage (%d total)", g.Name(id), amount, core.Damage)
	}
}

// DamageTarget pairs a creature with the amount of damage to deal it within a
// single simultaneous batch (see dealDamage). It is exported because effects
// build damage batches and pass them through the Resolver.
type DamageTarget struct {
	ID     LocalID
	Amount int
}

// dealDamage deals damage to every target simultaneously: each creature takes its
// damage (armor absorbs first) before any destruction is resolved, then all that
// were dealt lethal damage are destroyed together, in an order the controller
// chooses. Resolving destruction is part of dealing damage, so once dealDamage
// returns the dead creatures are already in the discard.
func (g *Game) dealDamage(controller int, targets ...DamageTarget) {
	for _, t := range targets {
		g.applyRawDamage(t.ID, t.Amount)
	}
	var dying []LocalID
	for _, t := range targets {
		if g.shouldDestroy(t.ID) {
			dying = append(dying, t.ID)
		}
	}
	g.destroyEach(controller, dying)
}

// shouldDestroy reports whether a creature is currently in a destroyable state:
// its damage meets or exceeds its power, its power is zero or less, or it has
// Poison and any damage on it.
func (g *Game) shouldDestroy(id LocalID) bool {
	def := g.cat.def(id)
	if def.Type != Creature || !g.inPlay(id) {
		return false
	}
	core := &g.State.Cards[id]
	poisoned := g.hasKeyword(id, Poison) && core.Damage > 0
	return int(core.Damage) >= g.Power(id) || g.Power(id) <= 0 || poisoned
}

// resetCore returns a card to its fresh, out-of-play state by zeroing its entire
// CardCore. Every field there is per-match, in-play state (damage, armor, Æmber,
// exhaustion, attached upgrades), so one zeroing covers them all and any field
// added later is reset automatically — no leaves-play path has to remember to
// clear it. Callers apply field-specific side effects (moving upgrades to the
// discard, handing Æmber to the opponent) BEFORE resetting.
func (g *Game) resetCore(id LocalID) { g.State.Cards[id] = CardCore{} }

// discardDestroyed moves a destroyed card from play to its owner's discard: it
// discards attached upgrades, hands any Æmber on the card to the owner's
// opponent, resets the card's per-match state, and adds it to the discard. It is
// the final step of destruction, run only for creatures still in play after their
// "Destroyed:" abilities resolve (see destroyTogether).
func (g *Game) discardDestroyed(id LocalID) {
	o := g.owner(id)
	g.State.Battleline[o].remove(id)
	g.State.Artifacts[o].remove(id)
	g.discardUpgrades(id)
	core := &g.State.Cards[id]
	if core.Amber > 0 {
		g.State.Aember[1-o] += int(core.Amber)
		g.logf("%d Æmber on %s goes to %s's pool", core.Amber, g.Name(id), g.names[1-o])
	}
	g.resetCore(id)
	g.State.Discard[o].add(id)
}

// discardUpgrades moves a card's attached upgrades to their owner's discard pile.
// A card that leaves play — destroyed or relocated — sheds its upgrades this way;
// they do not follow it to hand, deck, or archives.
func (g *Game) discardUpgrades(id LocalID) {
	o := g.owner(id)
	core := &g.State.Cards[id]
	for i := 0; i < int(core.UpgradeCount); i++ {
		g.State.Discard[o].add(core.Upgrades[i])
	}
}

// destroyTogether destroys several creatures as one simultaneous event, matching
// KeyForge timing. Every creature is tagged for destruction but stays in play
// while the "Destroyed:" abilities resolve, in the given order, so each ability
// sees the others still present. Then every creature still in play is moved to
// its owner's discard — a "Destroyed:" ability that relocated its creature (e.g.
// to hand or deck) has already removed it from play, so it is not also discarded.
func (g *Game) destroyTogether(ids []LocalID) {
	for _, id := range ids {
		g.logf("%s is destroyed", g.Name(id))
		g.triggerAbilities(id, TriggerDestroyed, 0, false)
	}
	for _, id := range ids {
		if g.inPlay(id) {
			g.discardDestroyed(id)
		}
	}
}

// destroyEach destroys each creature in ids simultaneously, letting the
// controller choose the order their "Destroyed:" abilities resolve. Callers pass
// a snapshot of distinct ids, so each is destroyed once.
func (g *Game) destroyEach(controller int, ids []LocalID) {
	ids = g.orderByChoice(controller, "Choose the next creature to destroy", ids)
	g.destroyTogether(ids)
}

// returnToTopOfDeck removes a card from play and places it on top of its owner's
// deck, clearing the per-match state it accrued while in play.
func (g *Game) returnToTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.State.Battleline[o].remove(id)
	g.State.Artifacts[o].remove(id)
	g.discardUpgrades(id)
	g.resetCore(id)
	g.State.Deck[o].addFront(id)
	g.logf("%s is put on top of %s's deck", g.Name(id), g.names[o])
}

// returnToHand removes a card from play and places it into its owner's hand,
// clearing the per-match state it accrued while in play.
func (g *Game) returnToHand(id LocalID) {
	o := g.owner(id)
	g.State.Battleline[o].remove(id)
	g.State.Artifacts[o].remove(id)
	g.discardUpgrades(id)
	g.resetCore(id)
	g.State.Hand[o].add(id)
	g.logf("%s is returned to %s's hand", g.Name(id), g.names[o])
}

// returnToArchives removes a card from play and places it into its owner's
// archives, clearing the per-match state it accrued while in play.
func (g *Game) returnToArchives(id LocalID) {
	o := g.owner(id)
	g.State.Battleline[o].remove(id)
	g.State.Artifacts[o].remove(id)
	g.discardUpgrades(id)
	g.resetCore(id)
	g.State.Archives[o].add(id)
	g.logf("%s is put into %s's archives", g.Name(id), g.names[o])
}

// fireCreatureEnters fires "after a creature enters play" abilities on every other
// in-play card, with the entering creature as the trigger target ("it").
func (g *Game) fireCreatureEnters(entered LocalID) {
	for player := 0; player < 2; player++ {
		for _, id := range g.allInPlay(player) {
			if id == entered {
				continue
			}
			g.triggerAbilities(id, TriggerAfterCreatureEnters, entered, true)
		}
	}
}

// fireArtifactPlayed fires "after you play an artifact" abilities on the playing
// player's other in-play cards, with the played artifact as "it".
func (g *Game) fireArtifactPlayed(player int, artifact LocalID) {
	for _, id := range g.allInPlay(player) {
		if id == artifact {
			continue
		}
		g.triggerAbilities(id, TriggerAfterArtifactPlayed, artifact, true)
	}
}

// triggerAbilities resolves every ability matching the trigger that the card
// carries itself or is granted by an attached upgrade.
func (g *Game) triggerAbilities(src LocalID, trigger Trigger, it LocalID, hasIt bool) {
	def := g.cat.def(src)
	fire := func(ab Ability) {
		if ab.Trigger != trigger {
			return
		}
		g.logf("%s: %s", def.Name, renderAbilityLine(def, ab))
		ab.Effect.Resolve(&EffectContext{
			Resolver:   g,
			Source:     src,
			Controller: g.owner(src),
			It:         it,
			HasIt:      hasIt,
		})
	}
	for _, ab := range def.Abilities {
		fire(ab)
	}
	core := &g.State.Cards[src]
	for i := 0; i < int(core.UpgradeCount); i++ {
		for _, ab := range g.cat.def(core.Upgrades[i]).Static.Granted {
			fire(ab)
		}
	}
}
