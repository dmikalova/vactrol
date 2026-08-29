package engine

import "fmt"

// This file holds the authoring effects for lasting "for the remainder of the
// turn" behavior: ForRemainderOfTurn installs a reaction (do something after an
// event), Instead installs a replacement (change an event's own outcome). Both are
// thin: they translate a composed effect into a flat lasting record via AddLasting
// (see game_lasting.go for how the records are fired and queried).

// ForRemainderOfTurn installs a reaction that runs for the rest of the controller's
// turn each time On occurs — Full Moon (gain Æmber when you play a creature),
// Charge! (deal damage when you play a creature), Crystal Hive (gain Æmber after a
// creature reaps). Do is the effect that runs each time; the supported effects are
// GainAember (crediting the controller) and DealDamage to an enemy creature.
type ForRemainderOfTurn struct {
	On Event
	Do Effect
}

// validate rejects a non-reaction event or a Do the reaction cannot carry.
func (e ForRemainderOfTurn) validate() error {
	if !e.On.isReaction() {
		return fmt.Errorf("ForRemainderOfTurn: On must be a reaction event")
	}
	if _, _, ok := lastingActionOf(e.Do); !ok {
		return fmt.Errorf("ForRemainderOfTurn: unsupported Do %T", e.Do)
	}
	if d, ok := e.Do.(DealDamage); ok && d.Target != (Target{Kind: TargetChosenEnemyCreature}) {
		return fmt.Errorf("ForRemainderOfTurn: DealDamage must target an enemy creature")
	}
	return validateEffect(e.Do)
}

// Text renders the effect, e.g. "for the remainder of the turn, each time you play
// a creature, gain 1 Æmber".
func (e ForRemainderOfTurn) Text() string {
	return "for the remainder of the turn, " + e.On.clause() + ", " + e.Do.Text()
}

// Resolve registers the reaction on the controller for the rest of the turn.
func (e ForRemainderOfTurn) Resolve(ctx *EffectContext) {
	action, amount, _ := lastingActionOf(e.Do)
	ctx.Resolver.AddLasting(e.On, action, ctx.Controller, amount)
}

// lastingActionOf maps a reaction's Do effect to the flat action and amount stored
// in the registry, reporting whether the effect is one the registry can carry.
func lastingActionOf(e Effect) (lastingAction, int, bool) {
	switch d := e.(type) {
	case GainAember:
		return actGainAember, d.Amount, true
	case DealDamage:
		return actDealDamage, d.Amount, true
	}
	return 0, 0, false
}

// Replacement is a lasting change to an event's own outcome, used by Instead.
type Replacement uint8

const (
	// Steal replaces gaining Æmber with stealing that much from the opponent.
	Steal Replacement = iota
)

// action maps the replacement to the flat action stored in the registry.
func (Replacement) action() lastingAction { return actSteal }

// text renders the replacement clause, e.g. "steal the same amount".
func (Replacement) text() string { return "steal the same amount" }

// Instead installs a replacement that, for the rest of the controller's turn,
// changes the outcome of the event Of before it happens — Dimension Door replaces
// gaining Æmber from reaping with stealing it.
type Instead struct {
	Of   Event
	With Replacement
}

// validate rejects an Of that is not a replacement event.
func (e Instead) validate() error {
	if e.Of.isReaction() {
		return fmt.Errorf("Instead: Of must be a replacement event")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, instead of gaining
// Æmber from reaping, steal the same amount".
func (e Instead) Text() string {
	return "for the remainder of the turn, instead of " + e.Of.gerund() + ", " + e.With.text()
}

// Resolve registers the replacement on the controller for the rest of the turn.
func (e Instead) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(e.Of, e.With.action(), ctx.Controller, 0)
}
