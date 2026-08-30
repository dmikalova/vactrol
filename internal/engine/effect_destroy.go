package engine

// Destroying a creature removes it from play. When an effect destroys several
// creatures they are destroyed simultaneously: every one is tagged for
// destruction and stays in play while their "Destroyed:" abilities resolve, in an
// order the controller chooses, so each ability sees the others still present;
// only then does each creature still in play move to the discard pile, along with
// its upgrades. A destroy effect can target every creature or only those matching
// a filter, such as "each creature with power 3 or lower".
//
//rulebook:effect Destroy
type Destroy struct {
	Target Target
}

// validate requires an explicit target.
func (e Destroy) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Destroy")
	}
	return nil
}

func (e Destroy) verb() string       { return "destroy" }
func (e Destroy) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "destroy each creature with power 3 or lower".
func (e Destroy) Text() string {
	return e.verb() + " " + e.targetText()
}

// Resolve destroys each selected creature, letting the controller order them.
func (e Destroy) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate destroys the selected creatures and reports whether any were, so
// Destroy can be the first half of a Then ("destroy a creature -> ...").
func (e Destroy) resolveGate(ctx *EffectContext) bool {
	ids := e.Target.Select(ctx)
	ctx.Resolver.DestroyEach(ctx.Controller, ids)
	return len(ids) > 0
}

// Choosing a creature and then destroying every creature that shares its power
// wipes out an entire power bracket at once — the chosen creature included. The
// choice fixes the power to match; the destruction then reaches both sides of the
// battle.
//
//rulebook:effect Destroy by Matching Power
type DestroySamePower struct{}

// Text renders the effect, binding the choice to its consequence with a dash.
func (DestroySamePower) Text() string {
	return "choose a creature - destroy each creature with the same power as the chosen creature"
}

// Resolve picks a creature, then destroys every creature whose power matches it.
func (DestroySamePower) Resolve(ctx *EffectContext) {
	all := append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Battleline(ctx.Opponent())...)
	chosen, ok := ctx.ChooseCreature("Choose a creature", all)
	if !ok {
		return
	}
	power := ctx.Resolver.Power(chosen)
	dying := make([]LocalID, 0, len(all))
	for _, id := range all {
		if ctx.Resolver.Power(id) == power {
			dying = append(dying, id)
		}
	}
	ctx.Resolver.DestroyEach(ctx.Controller, dying)
}
