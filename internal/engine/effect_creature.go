package engine

import "strings"

// Some abilities have you choose a single creature and then do one or more things
// to it — "Ready and fight with a friendly creature." OnChosenCreature models
// that: it picks one creature and applies an ordered list of CreatureVerbs to
// that same creature, so the verbs read as one sentence sharing a target.

// CreatureVerb is a single action applied to a chosen creature. Verbs are the
// building blocks of OnChosenCreature and compose into natural card text.
type CreatureVerb interface {
	VerbText() string
	Apply(ctx *EffectContext, target LocalID)
}

// ReadyVerb readies the chosen creature.
//
// Readying an exhausted creature stands it back up so it can be used again this
// turn (to reap, fight, or use an action).
type ReadyVerb struct{}

// VerbText returns the verb phrase.
func (ReadyVerb) VerbText() string { return "ready" }

// Apply readies the creature.
func (ReadyVerb) Apply(ctx *EffectContext, target LocalID) {
	ctx.Resolver.SetExhausted(target, false)
	ctx.Resolver.Logf("%s is readied", ctx.Resolver.Name(target))
}

// FightVerb makes the chosen creature fight an enemy creature.
//
// Using a creature to fight has it attack an enemy creature: both deal their
// power as damage to each other simultaneously (see Game.fight). The creature's
// controller chooses which enemy it fights.
type FightVerb struct{}

// VerbText returns the verb phrase.
func (FightVerb) VerbText() string { return "fight with" }

// Apply has the creature fight a chosen enemy creature.
func (FightVerb) Apply(ctx *EffectContext, target LocalID) {
	owner := ctx.Resolver.Owner(target)
	enemies := ctx.Resolver.Battleline(1 - owner)
	victim, ok := ctx.Resolver.ChooseCreature(owner, "Choose a creature to fight", enemies)
	if !ok {
		ctx.Resolver.Logf("%s has no creature to fight", ctx.Resolver.Name(target))
		return
	}
	ctx.Resolver.FightWith(target, victim)
}

// OnChosenCreature picks a single friendly or enemy creature and applies one or
// more verbs to it. The verbs share the single chosen target, which lets card
// text read naturally, e.g. "Ready and fight with a friendly creature."
type OnChosenCreature struct {
	Player Player
	// Neighbors restricts the choice to the source creature's battleline
	// neighbors instead of any friendly/enemy creature ("a neighboring creature").
	Neighbors bool
	// ExcludeHouse, when set, drops candidates of that house ("a friendly
	// non-Sanctum creature").
	ExcludeHouse House
	Verbs        []CreatureVerb
}

// noun renders the target noun phrase.
func (e OnChosenCreature) noun() string {
	if e.Neighbors {
		return "a neighboring creature"
	}
	if e.Player == Opponent {
		return "an enemy creature"
	}
	if e.ExcludeHouse != HouseNone {
		return "a friendly non-" + e.ExcludeHouse.String() + " creature"
	}
	return "a friendly creature"
}

// Text joins the verbs and the shared target, e.g.
// "ready and fight with a friendly creature".
func (e OnChosenCreature) Text() string {
	parts := make([]string, 0, len(e.Verbs))
	for _, v := range e.Verbs {
		parts = append(parts, v.VerbText())
	}
	return strings.Join(parts, " and ") + " " + e.noun()
}

// Resolve chooses the creature once, then applies each verb to it.
func (e OnChosenCreature) Resolve(ctx *EffectContext) {
	candidates := ctx.Resolver.Battleline(ctx.PlayerFor(e.Player))
	if e.Neighbors {
		candidates = neighbors(ctx, ctx.Source)
	}
	if e.ExcludeHouse != HouseNone {
		kept := make([]LocalID, 0, len(candidates))
		for _, id := range candidates {
			if ctx.Resolver.House(id) != e.ExcludeHouse {
				kept = append(kept, id)
			}
		}
		candidates = kept
	}
	chosen, ok := ctx.Resolver.ChooseCreature(ctx.Controller, "Choose "+e.noun(), candidates)
	if !ok {
		ctx.Resolver.Logf("no legal target for %q", e.Text())
		return
	}
	for _, v := range e.Verbs {
		v.Apply(ctx, chosen)
	}
}

// UseVerb uses the chosen creature. The controller picks how to use it — reap,
// fight, or its "Action:" ability — and that use resolves completely, nesting any
// further uses it triggers, before control returns. This is "Use a friendly
// creature": a creature can only be used while ready, so an already-exhausted
// creature may be chosen but nothing happens when it is used.
type UseVerb struct{}

// VerbText returns the verb phrase.
func (UseVerb) VerbText() string { return "use" }

// Apply offers the controller the target's available uses and resolves the chosen
// one.
func (UseVerb) Apply(ctx *EffectContext, target LocalID) {
	owner := ctx.Resolver.Owner(target)
	labels := []string{"reap"}
	uses := []func(){func() { ctx.Resolver.ReapWith(target) }}
	if len(ctx.Resolver.Battleline(1-owner)) > 0 {
		labels = append(labels, "fight")
		uses = append(uses, func() { FightVerb{}.Apply(ctx, target) })
	}
	if ctx.Resolver.HasAction(target) {
		labels = append(labels, "use its action")
		uses = append(uses, func() { ctx.Resolver.UseActionOf(target) })
	}
	idx := ctx.Resolver.ChooseOption(owner, "Choose how to use "+ctx.Resolver.Name(target), labels)
	if idx < 0 || idx >= len(uses) {
		return
	}
	uses[idx]()
}
