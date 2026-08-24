package game

import (
	"fmt"
	"strings"
)

// Effect is a node in a card's effect tree (its "AST"). Each effect knows how to
// render itself to English (Text) and how to resolve itself against a live game
// (Resolve). Keeping both behaviors on the same node guarantees the printed card
// text can never drift from what the card actually does.
type Effect interface {
	Text() string
	Resolve(ctx *EffectContext)
}

// EffectContext carries the state an effect needs while resolving. Cards are
// referenced by LocalID, keeping the context flat.
type EffectContext struct {
	Game       *Game
	Source     LocalID // the card whose ability is resolving
	Controller int     // the player who controls the ability
	It         LocalID // the triggering creature, for "it"-style targets
	HasIt      bool    // whether It is set
}

// TargetKind enumerates the ways an effect can select creatures.
type TargetKind int

const (
	// TargetThisCreature selects the source card itself.
	TargetThisCreature TargetKind = iota
	// TargetTriggeringCreature selects the creature that caused the trigger ("it").
	TargetTriggeringCreature
	// TargetEachCreature selects every creature in play.
	TargetEachCreature
	// TargetEachFriendlyCreature selects every creature the controller controls.
	TargetEachFriendlyCreature
	// TargetEachEnemyCreature selects every creature the opponent controls.
	TargetEachEnemyCreature
	// TargetEachArtifact selects every artifact in play, both players'.
	TargetEachArtifact
)

// Target describes which creatures an effect applies to.
type Target struct {
	Kind TargetKind
}

// Text renders the target as an English noun phrase.
func (t Target) Text() string {
	switch t.Kind {
	case TargetThisCreature:
		return "this creature"
	case TargetTriggeringCreature:
		return "it"
	case TargetEachCreature:
		return "each creature"
	case TargetEachFriendlyCreature:
		return "each friendly creature"
	case TargetEachArtifact:
		return "each artifact"
	case TargetEachEnemyCreature:
		return "each enemy creature"
	default:
		return "a creature"
	}
}

// Select resolves the target into concrete card ids in the current game state.
func (t Target) Select(ctx *EffectContext) []LocalID {
	switch t.Kind {
	case TargetThisCreature:
		return []LocalID{ctx.Source}
	case TargetTriggeringCreature:
		if ctx.HasIt {
			return []LocalID{ctx.It}
		}
		return nil
	case TargetEachArtifact:
		return append(ctx.Game.artifactsCopy(ctx.Controller), ctx.Game.artifactsCopy(1-ctx.Controller)...)
	case TargetEachCreature:
		return append(ctx.Game.battlelineCopy(ctx.Controller), ctx.Game.battlelineCopy(1-ctx.Controller)...)
	case TargetEachEnemyCreature:
		return ctx.Game.battlelineCopy(1 - ctx.Controller)
	default:
		return nil
	}
}

// DealDamage deals a fixed amount of damage to each selected creature.
type DealDamage struct {
	Amount int
	Target Target
}

// Text renders the effect, e.g. "deal 2 damage to each enemy creature".
func (e DealDamage) Text() string {
	return fmt.Sprintf("deal %d damage to %s", e.Amount, e.Target.Text())
}

// Resolve applies the damage and checks for destruction.
func (e DealDamage) Resolve(ctx *EffectContext) {
	ids := ctx.Game.orderByChoice(ctx.Controller, "Choose the next creature to damage", e.Target.Select(ctx))
	for _, id := range ids {
		ctx.Game.applyDamage(id, e.Amount)
		ctx.Game.checkDestroyed(id)
	}
}

// ReturnToDeck puts each selected card on top of its owner's deck.
type ReturnToDeck struct {
	Target Target
}

// Text renders the effect, e.g. "put each artifact on top of its owner's deck".
func (e ReturnToDeck) Text() string {
	return fmt.Sprintf("put %s on top of its owner's deck", e.Target.Text())
}

// Resolve moves each selected card from play to the top of its owner's deck.
func (e ReturnToDeck) Resolve(ctx *EffectContext) {
	ids := ctx.Game.orderByChoice(ctx.Controller, "Choose the next card to put on top of the deck", e.Target.Select(ctx))
	for _, id := range ids {
		ctx.Game.returnToTopOfDeck(id)
	}
}

// GainAember adds Æmber to the controller's pool.
type GainAember struct {
	Amount int
}

// Text renders the effect, e.g. "gain 1 Æmber".
func (e GainAember) Text() string { return fmt.Sprintf("gain %d Æmber", e.Amount) }

// Resolve adds the Æmber.
func (e GainAember) Resolve(ctx *EffectContext) {
	ctx.Game.State.Aember[ctx.Controller] += e.Amount
	ctx.Game.logf("%s gains %d Æmber", ctx.Game.names[ctx.Controller], e.Amount)
}

// Controller selects whose creatures a chosen-creature effect targets, relative
// to the effect's controller.
type Controller int

const (
	// Friendly targets the controller's own creatures.
	Friendly Controller = iota
	// Enemy targets the opponent's creatures.
	Enemy
)

// CreatureVerb is a single action applied to a chosen creature. Verbs are the
// building blocks of OnChosenCreature and compose into natural card text.
type CreatureVerb interface {
	VerbText() string
	Apply(ctx *EffectContext, target LocalID)
}

// ReadyVerb readies (un-exhausts) the chosen creature.
type ReadyVerb struct{}

// VerbText returns the verb phrase.
func (ReadyVerb) VerbText() string { return "ready" }

// Apply readies the creature.
func (ReadyVerb) Apply(ctx *EffectContext, target LocalID) {
	ctx.Game.State.Cards[target].Exhausted = false
	ctx.Game.logf("%s is readied", ctx.Game.Name(target))
}

// FightVerb makes the chosen creature fight an enemy creature.
type FightVerb struct{}

// VerbText returns the verb phrase.
func (FightVerb) VerbText() string { return "fight with" }

// Apply has the creature fight a chosen enemy creature.
func (FightVerb) Apply(ctx *EffectContext, target LocalID) {
	owner := ctx.Game.owner(target)
	enemies := ctx.Game.battlelineCopy(1 - owner)
	victim, ok := ctx.Game.chooserFor(owner).ChooseCreature("Choose a creature to fight", enemies)
	if !ok {
		ctx.Game.logf("%s has no creature to fight", ctx.Game.Name(target))
		return
	}
	ctx.Game.fight(target, victim)
}

// OnChosenCreature picks a single friendly or enemy creature and applies one or
// more verbs to it. The verbs share the single chosen target, which lets card
// text read naturally, e.g. "Ready and fight with a friendly creature."
type OnChosenCreature struct {
	Controller Controller
	Verbs      []CreatureVerb
}

// noun renders the target noun phrase.
func (e OnChosenCreature) noun() string {
	if e.Controller == Enemy {
		return "an enemy creature"
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
	targetPlayer := ctx.Controller
	if e.Controller == Enemy {
		targetPlayer = 1 - ctx.Controller
	}
	candidates := ctx.Game.battlelineCopy(targetPlayer)
	chosen, ok := ctx.Game.chooserFor(ctx.Controller).ChooseCreature("Choose "+e.noun(), candidates)
	if !ok {
		ctx.Game.logf("no legal target for %q", e.Text())
		return
	}
	for _, v := range e.Verbs {
		v.Apply(ctx, chosen)
	}
}

// Exalt places Æmber on a chosen creature from the common supply. "To exalt a
// creature" means put 1 Æmber on it, so exalting N times places N Æmber. The
// Æmber sits on the creature (CardCore.Amber); it enters no player's pool.
type Exalt struct {
	Controller Controller
	Times      int
}

// noun renders the target noun phrase.
func (e Exalt) noun() string {
	if e.Controller == Enemy {
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
	targetPlayer := ctx.Controller
	if e.Controller == Enemy {
		targetPlayer = 1 - ctx.Controller
	}
	candidates := ctx.Game.battlelineCopy(targetPlayer)
	chosen, ok := ctx.Game.chooserFor(ctx.Controller).ChooseCreature("Choose "+e.noun()+" to exalt", candidates)
	if !ok {
		ctx.Game.logf("no legal target for %q", e.Text())
		return
	}
	ctx.Game.State.Cards[chosen].Amber += int16(e.Times)
	ctx.Game.logf("%s is exalted (%d Æmber placed)", ctx.Game.Name(chosen), e.Times)
}

// Sequence resolves several effects in order and renders them joined by ", and".
type Sequence struct {
	Effects []Effect
}

// Text joins the child effect texts.
func (e Sequence) Text() string {
	parts := make([]string, 0, len(e.Effects))
	for _, child := range e.Effects {
		parts = append(parts, child.Text())
	}
	return strings.Join(parts, ", and ")
}

// Resolve resolves each child effect in order.
func (e Sequence) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}
