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
	// HasTrait reports whether a card has a trait.
	HasTrait(id LocalID, trait Trait) bool

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
	// OrderByChoice lets a player arrange ids into a resolution order.
	OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID
	// ChooseCreature asks a player to pick one creature from candidates.
	ChooseCreature(player int, prompt string, candidates []LocalID) (LocalID, bool)
	// ForceFight makes attacker fight defender (an ability-driven fight, with no
	// use/validation checks).
	ForceFight(attacker, defender LocalID)
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

// OrderByChoice is the Resolver entry point for orderByChoice.
func (g *Game) OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID {
	return g.orderByChoice(controller, prompt, ids)
}

// ChooseCreature asks a player to choose one creature from candidates.
func (g *Game) ChooseCreature(player int, prompt string, candidates []LocalID) (LocalID, bool) {
	return g.chooserFor(player).ChooseCreature(prompt, candidates)
}

// ForceFight makes attacker fight defender with no use/validation checks; it is
// how an ability drives a creature into combat.
func (g *Game) ForceFight(attacker, defender LocalID) { g.fight(attacker, defender) }

// Logf writes a line to the game log.
func (g *Game) Logf(format string, args ...any) { g.logf(format, args...) }
