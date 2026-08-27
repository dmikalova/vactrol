package engine

// Destroy destroys each creature its Target selects.
//
// Destroying a creature removes it from play. When an effect destroys several
// creatures they are destroyed simultaneously: every one is tagged for
// destruction and stays in play while their "Destroyed:" abilities resolve (in an
// order the controller chooses), so each ability sees the others still present;
// only then does each creature still in play move to the discard, along with its
// upgrades (see Game.destroyEach). Pairing Destroy with a filtered Target
// expresses many printed cards, e.g. "destroy each creature with power 3 or
// lower" or "destroy each Scientist trait creature".
type Destroy struct {
	Target Target
}

// Text renders the effect, e.g. "destroy each creature with power 3 or lower".
func (e Destroy) Text() string {
	return "destroy " + e.Target.Text()
}

// Resolve destroys each selected creature, letting the controller order them.
func (e Destroy) Resolve(ctx *EffectContext) {
	ctx.Resolver.DestroyEach(ctx.Controller, e.Target.Select(ctx))
}
