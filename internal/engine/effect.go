package engine

// This file and its effect_*.go / target.go siblings make up the card effect
// "AST": the small tree of nodes that both prints a card's rules text and
// carries it out. Each node type lives with related nodes in a file grouped by
// mechanic (effect_aember.go, effect_stun.go, ...), and every type's doc comment
// explains the mechanic in rulebook terms before the code shows how it is
// modelled. They all live in package engine because effects reach deep into the
// engine (dealing damage, destroying, drawing, choosing creatures); splitting
// them into a separate package would force the whole engine to be exported.

// Effect is one node in a card's effect tree. Each node knows how to render
// itself to English (Text) and how to carry itself out against a live game
// (Resolve). Because the same node does both, a card's printed rules text can
// never drift from what the card actually does.
type Effect interface {
	Text() string
	Resolve(ctx *EffectContext)
}

// validator is implemented by effects that can be misconfigured in a card
// definition. NewCard checks it when a card is built so a bad definition fails at
// startup instead of silently misbehaving when the ability resolves.
type validator interface {
	validate() error
}

// validateEffect returns an effect's configuration error, or nil if it has none.
// Composite effects (Sequence, Conditional) implement validator by descending
// into their children through this helper.
func validateEffect(e Effect) error {
	if v, ok := e.(validator); ok {
		return v.validate()
	}
	return nil
}

// EffectContext carries the state an effect needs while resolving. It exposes the
// game only through a Resolver, so an effect can inspect and change the game only
// via that interface — never by reaching into the state directly. Cards are
// referenced by LocalID, keeping the context flat.
type EffectContext struct {
	Resolver   Resolver
	Source     LocalID // the card whose ability is resolving
	Controller int     // the player who controls the ability
	It         LocalID // the triggering creature, for "it"-style targets
	HasIt      bool    // whether It is set
}

// PlayerFor resolves a relative Player (Controller/Opponent) to an absolute
// player index, relative to the ability's controller.
func (ctx *EffectContext) PlayerFor(p Player) int {
	if p == Opponent {
		return 1 - ctx.Controller
	}
	return ctx.Controller
}

// Player selects which player an effect targets, relative to the card's
// controller: Controller is the player who controls the card, Opponent is their
// opponent. Cards are written from the controller's point of view, so most
// effects default to the controller (the zero value).
type Player int

const (
	// Controller is the player who controls the card/ability.
	Controller Player = iota
	// Opponent is the controller's opponent.
	Opponent
)

// SelfName is a placeholder an effect's text uses to refer to its own source
// card; RenderCardText and the game log substitute it with the card's name so
// text like "{self} captures 1 Æmber" prints as "Charette captures 1 Æmber".
const SelfName = "{self}"
