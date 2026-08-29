package engine

import "fmt"

// This file gathers the Æmber-economy effects. Æmber is the resource players
// spend to forge keys (the way to win). It lives in a player's pool, except when
// it is captured or placed on a creature, where it belongs to no one until that
// creature leaves play.

// To gain Æmber, a player moves that many Æmber from the common supply into
// their pool — the ability's controller by default, or their opponent when the
// card says so. A "for each" clause multiplies the amount by a running count.
//
//rulebook:effect Gain Æmber
type GainAember struct {
	Player Player
	Amount int
	Per    Count
}

// validate rejects a GainAember whose player was left unset.
func (e GainAember) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("GainAember")
	}
	return nil
}

// Text renders the effect, e.g. "gain 1 Æmber" or "your opponent gains 2 Æmber".
// A "for each" count leads the sentence (rule 9), e.g. "for each key your opponent
// has forged, gain 1 Æmber".
func (e GainAember) Text() string {
	phrase := fmt.Sprintf("gain %d Æmber", e.Amount)
	if e.Player == Opponent {
		phrase = fmt.Sprintf("your opponent gains %d Æmber", e.Amount)
	}
	return forEach(e.Per, phrase)
}

// Resolve adds the Æmber to the selected player's pool.
func (e GainAember) Resolve(ctx *EffectContext) {
	p := ctx.PlayerFor(e.Player)
	amount := e.Amount
	if e.Per != nil {
		amount *= e.Per.Value(ctx)
	}
	ctx.Resolver.SetAember(p, ctx.Resolver.Aember(p)+amount)
	ctx.Resolver.Logf("%s gains %d Æmber", ctx.Resolver.PlayerName(p), amount)
}

// To lose Æmber, a player returns that many Æmber from their pool to the common
// supply. A pool can never go below zero, so a player told to lose more Æmber
// than they have simply loses all of it.
//
//rulebook:effect Lose Æmber
type LoseAember struct {
	Player Player
	Amount int
}

// validate rejects a LoseAember whose player was left unset.
func (e LoseAember) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("LoseAember")
	}
	return nil
}

// Text renders the effect, e.g. "lose 1 Æmber" or "your opponent loses 4 Æmber".
func (e LoseAember) Text() string {
	if e.Player == Opponent {
		return fmt.Sprintf("your opponent loses %d Æmber", e.Amount)
	}
	return fmt.Sprintf("lose %d Æmber", e.Amount)
}

// Resolve removes up to Amount Æmber from the selected player's pool.
func (e LoseAember) Resolve(ctx *EffectContext) {
	p := ctx.PlayerFor(e.Player)
	lost := min(e.Amount, ctx.Resolver.Aember(p))
	ctx.Resolver.SetAember(p, ctx.Resolver.Aember(p)-lost)
	ctx.Resolver.Logf("%s loses %d Æmber", ctx.Resolver.PlayerName(p), lost)
}

// Stealing Æmber moves it from the opponent's pool into your own. You can only
// steal as much Æmber as the opponent actually has.
//
//rulebook:effect Steal Æmber
type StealAember struct {
	Amount int
}

// Text renders the effect, e.g. "steal 1 Æmber".
func (e StealAember) Text() string { return fmt.Sprintf("steal %d Æmber", e.Amount) }

// Resolve moves the Æmber from the opponent's pool to the controller's.
func (e StealAember) Resolve(ctx *EffectContext) {
	opp := ctx.Opponent()
	amt := min(e.Amount, ctx.Resolver.Aember(opp))
	ctx.Resolver.SetAember(opp, ctx.Resolver.Aember(opp)-amt)
	ctx.Resolver.SetAember(ctx.Controller, ctx.Resolver.Aember(ctx.Controller)+amt)
	ctx.Resolver.Logf("%s steals %d Æmber from %s", ctx.Resolver.PlayerName(ctx.Controller), amt, ctx.Resolver.PlayerName(opp))
}

// Capturing Æmber moves it from the opponent's pool onto the capturing creature.
// Captured Æmber sits on the creature and counts for no player — it is out of
// every pool until the creature leaves play, at which point it goes to the pool
// of the creature's controller's opponent. You can only capture what the opponent has.
//
//rulebook:effect Capture Æmber
type CaptureAember struct {
	Amount int
	// All captures the opponent's entire pool instead of a fixed Amount.
	All bool
}

// Text renders the effect, e.g. "{self} captures 1 Æmber" or "{self} captures all
// your opponent's Æmber"; the source card's name replaces the self placeholder
// when the card is rendered.
func (e CaptureAember) Text() string {
	if e.All {
		return fmt.Sprintf("%s captures all your opponent's Æmber", SelfName)
	}
	return fmt.Sprintf("%s captures %d Æmber", SelfName, e.Amount)
}

// Resolve moves the Æmber from the opponent's pool onto the source creature.
func (e CaptureAember) Resolve(ctx *EffectContext) {
	opp := ctx.Opponent()
	amt := e.Amount
	if e.All {
		amt = ctx.Resolver.Aember(opp)
	}
	amt = min(amt, ctx.Resolver.Aember(opp))
	ctx.Resolver.SetAember(opp, ctx.Resolver.Aember(opp)-amt)
	ctx.Resolver.AddAmberOn(ctx.Source, amt)
	ctx.Resolver.Logf("%s captures %d Æmber", ctx.Resolver.Name(ctx.Source), amt)
}

// To exalt a creature is to place 1 Æmber from the common supply onto a chosen
// friendly or enemy creature. The Æmber sits on the creature, belonging to no
// pool, until it leaves play, then goes to the owner's opponent's pool. Exalting
// N times places N Æmber.
//
//rulebook:effect Exalt
type Exalt struct {
	Target Target
	Times  int
}

// Text renders the effect, e.g. "exalt an enemy creature 2 times". A single
// exalt drops the count so it reads naturally.
func (e Exalt) Text() string {
	if e.Times == 1 {
		return "exalt " + e.Target.Text()
	}
	return fmt.Sprintf("exalt %s %d times", e.Target.Text(), e.Times)
}

// Resolve chooses a creature (through the Target) and places Times Æmber on it.
func (e Exalt) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddAmberOn(id, e.Times)
		ctx.Resolver.Logf("%s is exalted (%d Æmber placed)", ctx.Resolver.Name(id), e.Times)
	}
}

// Some effects cap both pools at once: every player holding more than Keep Æmber
// loses the excess and is left with exactly Keep, while a player already at or
// below Keep is untouched. This reins in a runaway leader without punishing a
// player who has been spending.
//
//rulebook:effect Reduce Æmber
type EachPlayerLosesAllBut struct {
	Keep int
}

// Text renders the effect, e.g. "each player with 6 Æmber or more loses all but
// 5 Æmber".
func (e EachPlayerLosesAllBut) Text() string {
	return fmt.Sprintf("each player with %d Æmber or more loses all but %d Æmber", e.Keep+1, e.Keep)
}

// Resolve reduces each player's pool that exceeds Keep down to Keep.
func (e EachPlayerLosesAllBut) Resolve(ctx *EffectContext) {
	for p := 0; p < 2; p++ {
		if ctx.Resolver.Aember(p) > e.Keep {
			ctx.Resolver.SetAember(p, e.Keep)
			ctx.Resolver.Logf("%s is reduced to %d Æmber", ctx.Resolver.PlayerName(p), e.Keep)
		}
	}
}
