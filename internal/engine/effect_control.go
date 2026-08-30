package engine

// TakeControl moves this creature to the ability controller's battleline and
// makes that player its controller until the control-granting Upgrade leaves play.
// In KeyForge terms, control decides whose battleline the creature is in and who
// may use it; ownership stays fixed and still decides which discard, hand, deck,
// archives, or purge pile the card returns to when it leaves play.
type TakeControl struct{}

// Text renders Collar of Subordination's control-changing Play ability.
func (TakeControl) Text() string {
	return "take control of this creature until " + UpgradeName + " leaves play"
}

// Resolve changes control of this creature to the player resolving the ability.
func (TakeControl) Resolve(ctx *EffectContext) {
	ctx.Resolver.TakeControl(ctx.Source, ctx.Controller)
}
