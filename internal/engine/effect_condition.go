package engine

import "fmt"

// A Conditional gates an effect behind a check on the current game state — the
// "If ..." clause many cards open with, e.g. "If your opponent has 7 or more
// Æmber, they lose 4 Æmber." The Condition renders its own English (CondText) and
// evaluates itself (Met), so the printed text always matches what is checked.

// Condition is a boolean predicate on the live game, used by Conditional.
type Condition interface {
	CondText() string
	Met(ctx *EffectContext) bool
}

// OpponentAemberAtLeast is met when the opponent's pool holds at least Amount
// Æmber.
type OpponentAemberAtLeast struct {
	Amount int
}

// CondText renders the condition, e.g. "if your opponent has 7 Æmber or more".
func (c OpponentAemberAtLeast) CondText() string {
	return fmt.Sprintf("if your opponent has %d Æmber or more", c.Amount)
}

// Met reports whether the opponent has at least Amount Æmber.
func (c OpponentAemberAtLeast) Met(ctx *EffectContext) bool {
	return ctx.Resolver.Aember(ctx.Opponent()) >= c.Amount
}

// OpponentAemberExactly is met when the opponent's pool holds exactly Amount
// Æmber.
type OpponentAemberExactly struct {
	Amount int
}

// CondText renders the condition, e.g. "if your opponent has exactly 1 Æmber".
func (c OpponentAemberExactly) CondText() string {
	return fmt.Sprintf("if your opponent has exactly %d Æmber", c.Amount)
}

// Met reports whether the opponent has exactly Amount Æmber.
func (c OpponentAemberExactly) Met(ctx *EffectContext) bool {
	return ctx.Resolver.Aember(ctx.Opponent()) == c.Amount
}

// OpponentAemberMoreThanYou is met while the opponent's pool holds strictly more
// Æmber than the controller's.
type OpponentAemberMoreThanYou struct{}

// CondText renders the condition.
func (OpponentAemberMoreThanYou) CondText() string {
	return "if your opponent has more Æmber than you"
}

// Met reports whether the opponent has more Æmber than the controller.
func (OpponentAemberMoreThanYou) Met(ctx *EffectContext) bool {
	return ctx.Resolver.Aember(ctx.Opponent()) > ctx.Resolver.Aember(ctx.Controller)
}

// Conditional resolves Then only when Cond is met. It renders as "<cond>, <then>",
// e.g. "if your opponent has 7 Æmber or more, your opponent loses 4 Æmber".
type Conditional struct {
	Cond Condition
	Then Effect
}

// Text joins the condition and the gated effect.
func (e Conditional) Text() string { return e.Cond.CondText() + ", " + e.Then.Text() }

// Resolve runs Then only if Cond is met.
func (e Conditional) Resolve(ctx *EffectContext) {
	if e.Cond.Met(ctx) {
		e.Then.Resolve(ctx)
	}
}

// validate checks the gated effect for configuration errors.
func (e Conditional) validate() error {
	return validateEffect(e.Then)
}

// RepeatWhile resolves Do again and again for as long as Cond holds, re-checking
// after each pass — a self-looping effect such as "if your opponent has more
// Æmber than you, steal 1 Æmber -> repeat this effect". Do must make progress
// toward ending the loop (each pass changes the state Cond checks).
type RepeatWhile struct {
	Cond Condition
	Do   Effect
}

// Text renders the loop, leading with the condition and closing with the
// self-repeat gate.
func (e RepeatWhile) Text() string {
	return e.Cond.CondText() + ", " + e.Do.Text() + " -> repeat this effect"
}

// Resolve runs Do while Cond is met.
func (e RepeatWhile) Resolve(ctx *EffectContext) {
	for e.Cond.Met(ctx) {
		e.Do.Resolve(ctx)
	}
}

// validate checks the looped effect for configuration errors.
func (e RepeatWhile) validate() error {
	return validateEffect(e.Do)
}
