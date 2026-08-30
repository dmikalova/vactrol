package engine

// SwapWithFriendlyCreatureAndUse swaps this creature with another friendly
// creature in its controller's battleline, then uses that other creature.
//
// Battleline order matters for flanks and neighbors, so the swap is purely
// positional: no damage, upgrades, status, control, or other card state moves
// between the creatures. The chosen creature is then used normally, so if it is
// exhausted the use has no effect.
type SwapWithFriendlyCreatureAndUse struct{}

// Text renders Transposition Sandals' granted Action ability.
func (SwapWithFriendlyCreatureAndUse) Text() string {
	return "swap this creature with another friendly creature in your battleline. Use that other creature this turn"
}

// Resolve chooses another friendly creature, swaps its position with this
// creature, then uses the chosen creature.
func (SwapWithFriendlyCreatureAndUse) Resolve(ctx *EffectContext) {
	target := Target{Kind: TargetChosenOtherFriendlyCreature}
	for _, other := range target.Select(ctx) {
		ctx.Resolver.SwapBattlelinePositions(ctx.Source, other)
		UseVerb{}.Apply(ctx, other)
	}
}
