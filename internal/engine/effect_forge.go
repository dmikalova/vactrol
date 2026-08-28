package engine

// ForgeKey has the controller forge a key by paying its current cost, if they can
// afford it — firing "after you forge a key" abilities and, on the final key,
// winning the game. Cards use it to forge outside the normal start-of-turn step,
// e.g. "Lose 1 Æmber, and forge a key at its current cost."
type ForgeKey struct{}

// Text renders the effect.
func (ForgeKey) Text() string { return "forge a key at its current cost" }

// Resolve forges one key for the controller if affordable.
func (ForgeKey) Resolve(ctx *EffectContext) { ctx.Resolver.ForgeKey(ctx.Controller) }
