package engine

import "fmt"

// This file gathers the Æmber-economy effects. Æmber is the resource players
// spend to forge keys (the way to win). It lives in a player's pool, except when
// it is captured or placed on a creature, where it belongs to no one until that
// creature leaves play.

// GainAember adds Æmber to a player's pool. Player picks whose pool grows,
// relative to the ability's controller.
//
// Gaining Æmber simply moves it from the common supply into a pool.
type GainAember struct {
	Player Player
	Amount int
}

// Text renders the effect, e.g. "gain 1 Æmber" or "your opponent gains 2 Æmber".
func (e GainAember) Text() string {
	if e.Player == Opponent {
		return fmt.Sprintf("your opponent gains %d Æmber", e.Amount)
	}
	return fmt.Sprintf("gain %d Æmber", e.Amount)
}

// Resolve adds the Æmber to the selected player's pool.
func (e GainAember) Resolve(ctx *EffectContext) {
	p := ctx.PlayerFor(e.Player)
	ctx.Resolver.SetAember(p, ctx.Resolver.Aember(p)+e.Amount)
	ctx.Resolver.Logf("%s gains %d Æmber", ctx.Resolver.PlayerName(p), e.Amount)
}

// LoseAember removes Æmber from a player's pool. Player picks whose pool shrinks.
//
// A pool can never go below zero, so a player told to lose more Æmber than they
// have simply loses all of it.
type LoseAember struct {
	Player Player
	Amount int
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

// StealAember moves Æmber from the opponent's pool into the controller's pool.
//
// Stealing takes Æmber from the opponent and gives it to you; you can only steal
// as much as the opponent actually has.
type StealAember struct {
	Amount int
}

// Text renders the effect, e.g. "steal 1 Æmber".
func (e StealAember) Text() string { return fmt.Sprintf("steal %d Æmber", e.Amount) }

// Resolve moves the Æmber from the opponent's pool to the controller's.
func (e StealAember) Resolve(ctx *EffectContext) {
	opp := 1 - ctx.Controller
	amt := min(e.Amount, ctx.Resolver.Aember(opp))
	ctx.Resolver.SetAember(opp, ctx.Resolver.Aember(opp)-amt)
	ctx.Resolver.SetAember(ctx.Controller, ctx.Resolver.Aember(ctx.Controller)+amt)
	ctx.Resolver.Logf("%s steals %d Æmber from %s", ctx.Resolver.PlayerName(ctx.Controller), amt, ctx.Resolver.PlayerName(opp))
}

// CaptureAember moves Æmber from the opponent's pool onto the source creature.
//
// Captured Æmber sits on the creature and counts for no player — it is out of
// every pool until the creature leaves play, at which point it goes to the pool
// of the creature's owner's opponent. You can only capture what the opponent
// has.
type CaptureAember struct {
	Amount int
}

// Text renders the effect, e.g. "{self} captures 1 Æmber"; the source card's
// name is substituted for the self placeholder when the card is rendered.
func (e CaptureAember) Text() string {
	return fmt.Sprintf("%s captures %d Æmber", SelfName, e.Amount)
}

// Resolve moves the Æmber from the opponent's pool onto the source creature.
func (e CaptureAember) Resolve(ctx *EffectContext) {
	opp := 1 - ctx.Controller
	amt := min(e.Amount, ctx.Resolver.Aember(opp))
	ctx.Resolver.SetAember(opp, ctx.Resolver.Aember(opp)-amt)
	ctx.Resolver.AddAmberOn(ctx.Source, amt)
	ctx.Resolver.Logf("%s captures %d Æmber", ctx.Resolver.Name(ctx.Source), amt)
}

// Exalt places Æmber from the common supply onto a chosen creature. Player picks
// whether the creature is friendly or enemy.
//
// To exalt a creature is to put 1 Æmber on it; the Æmber sits on the creature
// (belonging to no pool) until it leaves play, then goes to the owner's
// opponent's pool. Exalting N times places N Æmber.
type Exalt struct {
	Player Player
	Times  int
}

// noun renders the target noun phrase.
func (e Exalt) noun() string {
	if e.Player == Opponent {
		return "an enemy creature"
	}
	return "a friendly creature"
}

// Text renders the effect, e.g. "exalt an enemy creature 2 times". A single
// exalt drops the count so it reads naturally.
func (e Exalt) Text() string {
	if e.Times == 1 {
		return "exalt " + e.noun()
	}
	return fmt.Sprintf("exalt %s %d times", e.noun(), e.Times)
}

// Resolve chooses a creature and places Times Æmber on it.
func (e Exalt) Resolve(ctx *EffectContext) {
	candidates := ctx.Resolver.Battleline(ctx.PlayerFor(e.Player))
	chosen, ok := ctx.Resolver.ChooseCreature(ctx.Controller, "Choose "+e.noun()+" to exalt", candidates)
	if !ok {
		ctx.Resolver.Logf("no legal target for %q", e.Text())
		return
	}
	ctx.Resolver.AddAmberOn(chosen, e.Times)
	ctx.Resolver.Logf("%s is exalted (%d Æmber placed)", ctx.Resolver.Name(chosen), e.Times)
}
