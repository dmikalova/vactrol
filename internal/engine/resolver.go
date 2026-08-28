package engine

// Resolver is the complete interface an effect uses to inspect and change the
// game. Effects hold only a Resolver (through EffectContext) — never the *Game or
// its GameState — so every change a card is able to make goes through one of
// these methods. The list below is therefore the full, explicit catalogue of
// what an effect is allowed to do; *Game implements it. The surface is
// deliberately wide: each new mechanic adds a method here, which surfaces the
// engine's effect-facing capability (and any gaps worth refactoring) in one place.
type Resolver interface {
	// ---- reads ----

	// Aember returns a player's Æmber pool.
	Aember(player int) int
	// AmberOn returns the Æmber sitting on a card (from capture, exalt, ...).
	AmberOn(id LocalID) int
	// Damage returns the damage currently on a creature.
	Damage(id LocalID) int
	// Power returns a creature's current power (including upgrades).
	Power(id LocalID) int
	// Name returns a card's printed name.
	Name(id LocalID) string
	// PlayerName returns a player's display name.
	PlayerName(player int) string
	// Owner returns the player who owns a card.
	Owner(id LocalID) int
	// Battleline returns a copy of a player's creatures, safe to hold across
	// mutations.
	Battleline(player int) []LocalID
	// Artifacts returns a copy of a player's artifacts.
	Artifacts(player int) []LocalID
	// Hand returns a copy of a player's hand.
	Hand(player int) []LocalID
	// Discard returns a copy of a player's discard pile.
	Discard(player int) []LocalID
	// Archives returns a copy of a player's archived cards.
	Archives(player int) []LocalID
	// IsCreature reports whether a card is a creature.
	IsCreature(id LocalID) bool
	// HasTrait reports whether a card has a trait.
	HasTrait(id LocalID, trait Trait) bool
	// House returns a card's house.
	House(id LocalID) House
	// Keys returns the number of keys a player has forged.
	Keys(player int) int
	// ForgeKey has a player forge one key at its current cost, if affordable.
	ForgeKey(player int)

	// ---- single-card / pool changes ----

	// SetAember sets a player's Æmber pool (never below zero).
	SetAember(player, amount int)
	// SetDamage sets the damage on a creature (never below zero).
	SetDamage(id LocalID, amount int)
	// SetStunned sets a creature's stun status.
	SetStunned(id LocalID, stunned bool)
	// SetExhausted sets a creature's exhausted status.
	SetExhausted(id LocalID, exhausted bool)
	// AddAmberOn changes the Æmber sitting on a card.
	AddAmberOn(id LocalID, delta int)

	// ---- compound actions the engine coordinates ----

	// DealDamage deals damage to each target simultaneously, then resolves
	// destruction (see the internal dealDamage).
	DealDamage(controller int, targets []DamageTarget)
	// DestroyEach destroys the given creatures as one simultaneous event.
	DestroyEach(controller int, ids []LocalID)
	// Draw makes a player draw count cards.
	Draw(controller, count int)
	// ReturnToTopOfDeck moves a card from play to the top of its owner's deck.
	ReturnToTopOfDeck(id LocalID)
	// ReturnToHand moves a card from play to its owner's hand.
	ReturnToHand(id LocalID)
	// ReturnToArchives moves a card from play to its owner's archives.
	ReturnToArchives(id LocalID)
	// ArchiveFromHand moves a card from its owner's hand to their archives.
	ArchiveFromHand(id LocalID)
	// ArchiveTopOfDeck moves the top card of a player's deck to their archives,
	// reporting whether a card was available.
	ArchiveTopOfDeck(player int) bool
	// DiscardArchives moves all of a player's archived cards to their discard pile.
	// The active player performs the discard, so they choose the order for their own
	// archives but get a random order for an opponent's (which they cannot see).
	DiscardArchives(owner int)
	// ReturnFromDiscardToHand moves a card from its owner's discard to their hand.
	ReturnFromDiscardToHand(id LocalID)
	// ReturnFromDiscardToTopOfDeck moves a card from its owner's discard to the top
	// of their deck.
	ReturnFromDiscardToTopOfDeck(id LocalID)
	// OrderByChoice lets a player arrange ids into a resolution order.
	OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID
	// ChooseCreature asks a player to pick one creature from candidates.
	ChooseCreature(player int, prompt string, candidates []LocalID) (LocalID, bool)
	// ChooseOption asks a player to pick one of several labeled options, returning
	// its index (0 when the player's chooser expresses no preference).
	ChooseOption(player int, prompt string, options []string) int
	// FightWith makes attacker fight defender (ability-driven, ignoring active
	// player and house). A creature can only be used while ready, so an exhausted
	// attacker may be chosen but does nothing.
	FightWith(attacker, defender LocalID)
	// ReapWith reaps with a creature (ability-driven, ignoring active player and
	// house). A creature can only be used while ready, so an exhausted creature may
	// be chosen but does nothing.
	ReapWith(id LocalID)
	// UseActionOf fires a card's "Action:" ability (ability-driven, ignoring active
	// player and house). A card can only be used while ready, so an exhausted card
	// may be chosen but does nothing.
	UseActionOf(id LocalID)
	// HasAction reports whether a card has an "Action:" ability.
	HasAction(id LocalID) bool
	// Logf writes a line to the game log.
	Logf(format string, args ...any)
}

// The read accessors Aember, AmberOn, Damage, Power, Name, PlayerName,
// Battleline, and Artifacts are defined in game.go; the remaining Resolver
// methods follow.

// Owner returns the player who owns a card.
func (g *Game) Owner(id LocalID) int { return g.owner(id) }

// HasTrait reports whether a card has a trait.
func (g *Game) HasTrait(id LocalID, trait Trait) bool { return g.cat.def(id).hasTrait(trait) }

// House returns a card's house.
func (g *Game) House(id LocalID) House { return g.cat.def(id).House }

// ForgeKey has a player forge one key at its current cost, if affordable.
func (g *Game) ForgeKey(player int) { g.forgeOneKey(player) }

// IsCreature reports whether a card is a creature.
func (g *Game) IsCreature(id LocalID) bool { return g.cat.def(id).Type == Creature }

// SetAember sets a player's Æmber pool, clamped at zero.
func (g *Game) SetAember(player, amount int) {
	if amount < 0 {
		amount = 0
	}
	g.State.Aember[player] = amount
}

// SetDamage sets the damage on a creature, clamped at zero.
func (g *Game) SetDamage(id LocalID, amount int) {
	if amount < 0 {
		amount = 0
	}
	g.State.Cards[id].Damage = int16(amount)
}

// SetStunned sets a creature's stun status.
func (g *Game) SetStunned(id LocalID, stunned bool) { g.State.Cards[id].Stunned = stunned }

// SetExhausted sets a creature's exhausted status.
func (g *Game) SetExhausted(id LocalID, exhausted bool) { g.State.Cards[id].Exhausted = exhausted }

// AddAmberOn changes the Æmber sitting on a card.
func (g *Game) AddAmberOn(id LocalID, delta int) { g.State.Cards[id].Amber += int16(delta) }

// DealDamage is the Resolver entry point for the internal dealDamage.
func (g *Game) DealDamage(controller int, targets []DamageTarget) {
	g.dealDamage(controller, targets...)
}

// DestroyEach is the Resolver entry point for destroyEach.
func (g *Game) DestroyEach(controller int, ids []LocalID) { g.destroyEach(controller, ids) }

// Draw is the Resolver entry point for the internal draw.
func (g *Game) Draw(controller, count int) { g.draw(controller, count) }

// ReturnToTopOfDeck is the Resolver entry point for returnToTopOfDeck.
func (g *Game) ReturnToTopOfDeck(id LocalID) { g.returnToTopOfDeck(id) }

// ReturnToHand is the Resolver entry point for returnToHand.
func (g *Game) ReturnToHand(id LocalID) { g.returnToHand(id) }

// ReturnToArchives is the Resolver entry point for returnToArchives.
func (g *Game) ReturnToArchives(id LocalID) { g.returnToArchives(id) }

// ArchiveFromHand moves a card from its owner's hand to their archives.
func (g *Game) ArchiveFromHand(id LocalID) { g.archiveFromHand(g.owner(id), id) }

// ArchiveTopOfDeck moves the top card of a player's deck to their archives.
func (g *Game) ArchiveTopOfDeck(player int) bool { return g.archiveTopOfDeck(player) }

// DiscardArchives moves all of a player's archived cards to their discard pile.
func (g *Game) DiscardArchives(owner int) { g.discardArchives(owner) }

// ReturnFromDiscardToHand moves a card from its owner's discard pile to their hand.
func (g *Game) ReturnFromDiscardToHand(id LocalID) {
	o := g.owner(id)
	g.State.Discard[o].remove(id)
	g.State.Hand[o].add(id)
	g.logf("%s returns %s from their discard to hand", g.names[o], g.Name(id))
}

// ReturnFromDiscardToTopOfDeck moves a card from its owner's discard pile to the
// top of their deck.
func (g *Game) ReturnFromDiscardToTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.State.Discard[o].remove(id)
	g.State.Deck[o].addFront(id)
	g.logf("%s puts %s from their discard on top of their deck", g.names[o], g.Name(id))
}

// OrderByChoice is the Resolver entry point for orderByChoice.
func (g *Game) OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID {
	return g.orderByChoice(controller, prompt, ids)
}

// ChooseCreature asks a player to choose one creature from candidates.
func (g *Game) ChooseCreature(player int, prompt string, candidates []LocalID) (LocalID, bool) {
	return g.chooserFor(player).ChooseCreature(prompt, candidates)
}

// ChooseOption asks a player to choose one of several labeled options. If the
// player's chooser expresses no preference (does not implement OptionChooser),
// the first option is taken.
func (g *Game) ChooseOption(player int, prompt string, options []string) int {
	if oc, ok := g.chooserFor(player).(OptionChooser); ok {
		return oc.ChooseOption(prompt, options)
	}
	return 0
}

// FightWith makes attacker fight defender, ability-driven (ignoring active player
// and house). A creature can only be used while ready, so an exhausted attacker
// does nothing.
func (g *Game) FightWith(attacker, defender LocalID) {
	if g.readyToUse(attacker) {
		g.fight(attacker, defender)
	}
}

// ReapWith reaps with a creature, ability-driven (ignoring active player and
// house). A creature can only be used while ready, so an exhausted creature does
// nothing.
func (g *Game) ReapWith(id LocalID) {
	if g.readyToUse(id) {
		g.reapWith(id)
	}
}

// UseActionOf fires a card's "Action:" ability, ability-driven (ignoring active
// player and house). A card can only be used while ready, so an exhausted card
// does nothing.
func (g *Game) UseActionOf(id LocalID) {
	if g.readyToUse(id) {
		g.useActionOf(id)
	}
}

// HasAction reports whether a card has an "Action:" ability.
func (g *Game) HasAction(id LocalID) bool { return g.cat.def(id).hasTrigger(TriggerAction) }

// Logf writes a line to the game log.
func (g *Game) Logf(format string, args ...any) { g.logf(format, args...) }
