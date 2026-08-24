package game

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

// drawTo draws from the top of the deck until the hand reaches n (or deck empty).
func (g *Game) drawTo(player, n int) {
	hand := &g.State.Hand[player]
	deck := &g.State.Deck[player]
	for int(hand.Count) < n && deck.Count > 0 {
		hand.add(deck.removeAt(0))
	}
}

// forgeKeys forges as many keys as the player can afford, firing "after you forge
// a key" abilities for each key forged.
func (g *Game) forgeKeys(player int) {
	for g.State.Aember[player] >= KeyCost {
		g.State.Aember[player] -= KeyCost
		g.State.Keys[player]++
		g.logf("%s forges a key (%d/%d)", g.names[player], g.State.Keys[player], KeysToWin)
		for _, id := range g.allInPlay(player) {
			g.triggerAbilities(id, TriggerAfterForgeKey, 0, false)
		}
		if g.State.Keys[player] >= KeysToWin {
			g.State.Winner = player
			g.logf("%s wins the game!", g.names[player])
			return
		}
	}
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
	g.triggerAbilities(id, TriggerPlay, 0, false)
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
	g.triggerAbilities(id, TriggerPlay, 0, false)
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
	g.triggerAbilities(id, TriggerPlay, 0, false)
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
	def := g.cat.def(id)
	if g.State.ActiveHouse != HouseNone && def.House != g.State.ActiveHouse {
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
	if g.State.ActiveHouse != HouseNone && def.House != g.State.ActiveHouse {
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
	g.State.Cards[id].Exhausted = true
	g.State.Aember[player]++
	g.logf("%s reaps with %s (+1 Æmber)", g.names[player], g.Name(id))
	g.triggerAbilities(id, TriggerReap, 0, false)
	return nil
}

// UseAction uses a creature's "Action:" ability.
func (g *Game) UseAction(player int, id LocalID) error {
	if err := g.canUse(player, id); err != nil {
		return err
	}
	if !g.cat.def(id).hasTrigger(TriggerAction) {
		return ErrWrongType
	}
	g.State.Cards[id].Exhausted = true
	g.logf("%s uses %s's action ability", g.names[player], g.Name(id))
	g.triggerAbilities(id, TriggerAction, 0, false)
	return nil
}

// Fight uses attacker to fight the enemy creature defender.
func (g *Game) Fight(player int, attacker, defender LocalID) error {
	if err := g.canUse(player, attacker); err != nil {
		return err
	}
	if g.cat.def(defender).Type != Creature || g.owner(defender) == player || !g.inPlay(defender) {
		return ErrNoTarget
	}
	g.fight(attacker, defender)
	return nil
}

// canUse validates that a creature may be used (reap/fight/action) right now.
func (g *Game) canUse(player int, id LocalID) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	def := g.cat.def(id)
	if g.owner(id) != player || def.Type != Creature || !g.inPlay(id) {
		return ErrWrongType
	}
	core := &g.State.Cards[id]
	if core.Exhausted {
		return ErrCardExhausted
	}
	if g.State.ActiveHouse != HouseNone && def.House != g.State.ActiveHouse {
		return ErrWrongHouse
	}
	return nil
}

// CanUse reports whether a creature may currently be used (reap/fight/action) by
// the player: nil if usable, otherwise the reason. UIs can call it to reject an
// action before prompting for a target.
func (g *Game) CanUse(player int, id LocalID) error { return g.canUse(player, id) }

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
	if g.State.ActiveHouse != HouseNone && def.House != g.State.ActiveHouse {
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
	g.State.Cards[attacker].Exhausted = true
	ap, dp := g.Power(attacker), g.Power(defender)
	g.logf("%s (%d power) fights %s (%d power)", g.Name(attacker), ap, g.Name(defender), dp)
	g.applyDamage(defender, ap)
	if !g.cat.def(attacker).hasKeyword(Skirmish) {
		g.applyDamage(attacker, dp)
	}
	g.checkDestroyed(defender)
	g.checkDestroyed(attacker)
	g.triggerAbilities(attacker, TriggerFight, 0, false)
}

// applyDamage applies damage to a creature, absorbing it with armor first.
func (g *Game) applyDamage(id LocalID, amount int) {
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

// checkDestroyed destroys a creature whose damage meets its power, whose power is
// zero or less, or that has Poison and any damage on it.
func (g *Game) checkDestroyed(id LocalID) {
	def := g.cat.def(id)
	if def.Type != Creature || !g.inPlay(id) {
		return
	}
	core := &g.State.Cards[id]
	poisoned := def.hasKeyword(Poison) && core.Damage > 0
	if int(core.Damage) >= g.Power(id) || g.Power(id) <= 0 || poisoned {
		g.destroy(id)
	}
}

// resetCore returns a card to its fresh, out-of-play state by zeroing its entire
// CardCore. Every field there is per-match, in-play state (damage, armor, Æmber,
// exhaustion, attached upgrades), so one zeroing covers them all and any field
// added later is reset automatically — no leaves-play path has to remember to
// clear it. Callers apply field-specific side effects (moving upgrades to the
// discard, handing Æmber to the opponent) BEFORE resetting.
func (g *Game) resetCore(id LocalID) { g.State.Cards[id] = CardCore{} }

// destroy removes a creature from play, discards it and its upgrades, and fires
// "Destroyed:" abilities.
func (g *Game) destroy(id LocalID) {
	o := g.owner(id)
	g.State.Battleline[o].remove(id)
	core := &g.State.Cards[id]
	for i := 0; i < int(core.UpgradeCount); i++ {
		g.State.Discard[o].add(core.Upgrades[i])
	}
	if core.Amber > 0 {
		g.State.Aember[1-o] += int(core.Amber)
		g.logf("%d Æmber on %s goes to %s's pool", core.Amber, g.Name(id), g.names[1-o])
	}
	g.resetCore(id)
	g.State.Discard[o].add(id)
	g.logf("%s is destroyed", g.Name(id))
	g.triggerAbilities(id, TriggerDestroyed, 0, false)
}

// returnToTopOfDeck removes a card from play and places it on top of its owner's
// deck, clearing the per-match state it accrued while in play.
func (g *Game) returnToTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.State.Battleline[o].remove(id)
	g.State.Artifacts[o].remove(id)
	g.resetCore(id)
	g.State.Deck[o].addFront(id)
	g.logf("%s is put on top of %s's deck", g.Name(id), g.names[o])
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

// triggerAbilities resolves every ability on a card matching the trigger.
func (g *Game) triggerAbilities(src LocalID, trigger Trigger, it LocalID, hasIt bool) {
	def := g.cat.def(src)
	for _, ab := range def.Abilities {
		if ab.Trigger != trigger {
			continue
		}
		g.logf("%s: %s", def.Name, RenderAbility(ab))
		ab.Effect.Resolve(&EffectContext{
			Game:       g,
			Source:     src,
			Controller: g.owner(src),
			It:         it,
			HasIt:      hasIt,
		})
	}
}
