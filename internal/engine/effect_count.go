package engine

// Count computes a number from live game state, for effects whose magnitude
// scales with the board — e.g. "gain 1 Æmber for each key your opponent has
// forged". It both yields the value (Value) and renders the "for each ..." tail
// of the effect's text (CountText), so the printed card and its behavior share
// one source, exactly like Effect and Condition.
type Count interface {
	Value(ctx *EffectContext) int
	CountText() string
}

// OpponentForgedKeys counts the keys the controller's opponent has forged.
type OpponentForgedKeys struct{}

// Value returns the opponent's forged-key count.
func (OpponentForgedKeys) Value(ctx *EffectContext) int {
	return ctx.Resolver.Keys(1 - ctx.Controller)
}

// CountText renders the singular noun the "for each" clause repeats.
func (OpponentForgedKeys) CountText() string { return "key your opponent has forged" }

// CardsInArchives counts the cards in the controller's archives.
type CardsInArchives struct{}

// Value returns how many cards the controller has archived.
func (CardsInArchives) Value(ctx *EffectContext) int {
	return len(ctx.Resolver.Archives(ctx.Controller))
}

// CountText renders the singular noun the "for each" clause repeats.
func (CardsInArchives) CountText() string { return "card in your archives" }
