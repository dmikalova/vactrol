package engine

// Count computes a number from live game state, for effects whose magnitude
// scales with the board — e.g. "for each key your opponent has forged, gain 1
// Æmber". It both yields the value (Value) and renders the leading "for each ..."
// clause of the effect's text (CountText), so the printed card and its behavior
// share one source, exactly like Effect and Condition.
type Count interface {
	Value(ctx *EffectContext) int
	CountText() string
}

// forEach front-loads a Per count onto an effect's body as a leading "for each
// ..." clause, so the sentence reads subject-first and forward (card-wording rule
// 9). With no count it returns body unchanged.
func forEach(per Count, body string) string {
	if per == nil {
		return body
	}
	return "for each " + per.CountText() + ", " + body
}

// OpponentForgedKeys counts the keys the controller's opponent has forged.
type OpponentForgedKeys struct{}

// Value returns the opponent's forged-key count.
func (OpponentForgedKeys) Value(ctx *EffectContext) int {
	return ctx.Resolver.Keys(ctx.Opponent())
}

// CountText renders the singular noun the "for each" clause repeats.
func (OpponentForgedKeys) CountText() string { return "key your opponent has forged" }

// CardsInArchives counts the cards in a player's archives.
type CardsInArchives struct{ Player Player }

// Value returns how many cards the chosen player has archived.
func (e CardsInArchives) Value(ctx *EffectContext) int {
	return len(ctx.Resolver.Archives(ctx.PlayerFor(e.Player)))
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CardsInArchives) CountText() string {
	if e.Player == Opponent {
		return "card in your opponent's archives"
	}
	return "card in your archives"
}

// FriendlyCreaturesInPlay counts the creatures the controller has in play.
type FriendlyCreaturesInPlay struct{}

// Value returns how many creatures the controller has in play.
func (FriendlyCreaturesInPlay) Value(ctx *EffectContext) int {
	return len(ctx.Resolver.Battleline(ctx.Controller))
}

// CountText renders the singular noun the "for each" clause repeats.
func (FriendlyCreaturesInPlay) CountText() string { return "friendly creature in play" }
