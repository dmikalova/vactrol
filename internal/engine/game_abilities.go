package engine

import "strings"

// This file holds how a card is USED — the player actions that reap, fight, or
// activate an "Action:" ability — the checks that gate them, and the machinery
// that fires a card's triggered abilities. Combat resolution lives in
// game_combat.go and destruction in game_destroy.go.

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
	g.gainReapAember(p, id)
	g.triggerAbilities(id, TriggerAfterReap, 0, false)
	g.fireLasting(EventReap, p, id)
}

// gainReapAember pays out the Æmber a reap grants to player p, applying any lasting
// replacement of that payout (Dimension Door makes it steal instead of gain).
func (g *Game) gainReapAember(p int, source LocalID) {
	if act, ok := g.lastingReplacement(p, EventReapAember); ok && act == actSteal {
		if stolen := min(1, g.State.Aember[1-p]); stolen > 0 {
			g.State.Aember[1-p] -= stolen
			g.State.Aember[p] += stolen
			g.logf("%s reaps with %s, stealing %d Æmber", g.names[p], g.Name(source), stolen)
		} else {
			g.logf("%s reaps with %s (no Æmber to steal)", g.names[p], g.Name(source))
		}
		return
	}
	g.State.Aember[p]++
	g.logf("%s reaps with %s (+1 Æmber)", g.names[p], g.Name(source))
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
	if g.cannotFight(player) {
		return ErrCannotFight
	}
	if err := g.canUse(player, attacker); err != nil {
		// A grant such as Brothers in Battle lets a creature of a chosen house fight
		// out of the active house, so a wrong-house error is forgiven when the grant
		// covers this attacker; every other check still applies.
		if !(err == ErrWrongHouse && g.mayFightOutOfHouse(attacker)) {
			return err
		}
	}
	if g.cat.def(defender).Type != Creature || g.owner(defender) == player || !g.inPlay(defender) {
		return ErrNoTarget
	}
	if fr := g.cat.def(attacker).FightRestriction; fr != (Target{}) &&
		!fr.allows(&EffectContext{Resolver: g, Source: attacker, Controller: player}, defender) {
		return ErrNoTarget
	}
	if g.recoverFromStun(attacker) {
		return nil
	}
	g.fight(attacker, defender)
	return nil
}

// mayFightOutOfHouse reports whether a fight grant (Brothers in Battle) lets the
// attacker fight this turn despite not being in the active house.
func (g *Game) mayFightOutOfHouse(attacker LocalID) bool {
	h := g.State.MayFightHouse[g.owner(attacker)]
	return h != HouseNone && g.cat.def(attacker).House == h
}

// FightTargets returns the enemy creatures the attacker may legally fight right
// now, mirroring the checks in Fight. It is empty when the player is barred from
// fighting, the attacker cannot be used (readiness, or wrong house without a
// fight grant), or no enemy creature satisfies the attacker's fight restriction.
// A UI offers Fight only when this is non-empty.
func (g *Game) FightTargets(player int, attacker LocalID) []LocalID {
	if g.cannotFight(player) {
		return nil
	}
	if err := g.canUse(player, attacker); err != nil &&
		!(err == ErrWrongHouse && g.mayFightOutOfHouse(attacker)) {
		return nil
	}
	fr := g.cat.def(attacker).FightRestriction
	var targets []LocalID
	for _, def := range g.State.Battleline[1-player].slice() {
		if fr != (Target{}) &&
			!fr.allows(&EffectContext{Resolver: g, Source: attacker, Controller: player}, def) {
			continue
		}
		targets = append(targets, def)
	}
	return targets
}

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

// fireUpgradePlay resolves an upgrade's own "Play:" abilities when it is attached,
// acting on its host creature. The upgrade is the card that was played, but its
// ability speaks of "this creature" — the host — so the host is the effect source
// that self-references resolve to.
func (g *Game) fireUpgradePlay(host LocalID, up *CardDefinition) {
	for _, ab := range up.Abilities {
		if ab.Trigger != TriggerAfterPlay {
			continue
		}
		g.logf("%s: %s", up.Name, strings.ReplaceAll(RenderAbility(ab), SelfName, g.Name(host)))
		ab.Effect.Resolve(&EffectContext{
			Resolver:   g,
			Source:     host,
			Controller: g.owner(host),
		})
	}
}

// fireCreatureEnters is the enter-play event for a creature. It first resolves the
// entering creature's own "enters play" abilities (Chuff Ape entering stunned),
// then fires "after a creature enters play" abilities on every other in-play card,
// with the entering creature as the trigger target ("it").
func (g *Game) fireCreatureEnters(entered LocalID) {
	g.triggerAbilities(entered, TriggerEntersPlay, 0, false)
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
// carries itself, is granted by an attached upgrade, or is granted by an in-play
// card's constant ability (Annihilation Ritual's "Destroyed: purge this creature"
// on each creature). A "Destroyed:" ability may take the card out of play (purge,
// return to hand); once it does, the remaining abilities are skipped, since the
// card is no longer there to fire them.
func (g *Game) triggerAbilities(src LocalID, trigger Trigger, it LocalID, hasIt bool) {
	def := g.cat.def(src)
	fire := func(ab Ability) {
		// A Play ability on an action resolves while its source is between hand and
		// discard, so only Destroyed abilities require their source to remain in play.
		if ab.Trigger != trigger || (trigger == TriggerDestroyed && !g.inPlay(src)) {
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
	for player := 0; player < 2; player++ {
		for _, srcCard := range g.allInPlay(player) {
			c := g.cat.def(srcCard).Constant
			if len(c.Granted) == 0 || !g.constantAffects(srcCard, c, src) {
				continue
			}
			for _, ab := range c.Granted {
				fire(ab)
			}
		}
	}
}

// triggeredAbility is an ability waiting to resolve from a creature being
// destroyed. It retains the creature as Source so a granted "purge this creature"
// resolves against the creature that gained it, rather than the card that granted it.
type triggeredAbility struct {
	source  LocalID
	ability Ability
}

// destroyedAbilities collects every Destroyed ability the creatures about to be
// destroyed carry: printed, upgrade-granted, and constant-granted. The collection
// happens before any resolves, then destroyTogether lets the active player order
// the whole set as KeyForge requires.
func (g *Game) destroyedAbilities(ids []LocalID) []triggeredAbility {
	var pending []triggeredAbility
	for _, id := range ids {
		def := g.cat.def(id)
		for _, ab := range def.Abilities {
			if ab.Trigger == TriggerDestroyed {
				pending = append(pending, triggeredAbility{source: id, ability: ab})
			}
		}
		for _, upgrade := range g.Upgrades(id) {
			for _, ab := range g.cat.def(upgrade).Static.Granted {
				if ab.Trigger == TriggerDestroyed {
					pending = append(pending, triggeredAbility{source: id, ability: ab})
				}
			}
		}
		for player := 0; player < 2; player++ {
			for _, grantor := range g.allInPlay(player) {
				c := g.cat.def(grantor).Constant
				if !g.constantAffects(grantor, c, id) {
					continue
				}
				for _, ab := range c.Granted {
					if ab.Trigger == TriggerDestroyed {
						pending = append(pending, triggeredAbility{source: id, ability: ab})
					}
				}
			}
		}
	}
	return pending
}
