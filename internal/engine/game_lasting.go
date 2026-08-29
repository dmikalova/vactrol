package engine

// This file holds LASTING EFFECTS: "for the remainder of the turn" effects that
// attach to a later game event. Two flavors share the same flat registry:
//
//   - a REACTION runs after an event (Full Moon gains Æmber after you play a
//     creature, Charge! deals damage after you play a creature, Crystal Hive gains
//     Æmber after a creature reaps). Reactions are fired by fireLasting, which
//     gathers every reaction responding to an event and lets the controller order
//     them when several fire at once.
//   - a REPLACEMENT changes an event's own outcome before it happens (Dimension
//     Door makes reaping steal Æmber instead of gaining it). The event site queries
//     for a replacement (lastingReplacement) and applies it in place.
//
// Because the game state is a flat, pointerless value it cannot hold effect
// closures, so a lasting effect is a small enum-tagged record (LastingEffect)
// resolved by a central switch. Adding a new one is a new Event and/or
// lastingAction plus a case here — never a bespoke branch in the play/reap path.

// Event names a moment a lasting effect attaches to.
type Event uint8

const (
	// EventCreaturePlayed fires after the controller plays a creature (a reaction
	// point). Full Moon and Charge! attach here.
	EventCreaturePlayed Event = iota
	// EventReap fires after one of the controller's creatures reaps (a reaction
	// point). Crystal Hive attaches here.
	EventReap
	// EventReapAember is the Æmber a reap grants (a replacement point). Dimension
	// Door replaces it.
	EventReapAember
)

// isReaction reports whether the event is a reaction point (fired after) rather
// than a replacement point (queried during).
func (e Event) isReaction() bool { return e == EventCreaturePlayed || e == EventReap }

// clause renders the "when" phrase for a reaction, e.g. "each time you play a
// creature".
func (e Event) clause() string {
	if e == EventReap {
		return "after a creature reaps"
	}
	return "each time you play a creature"
}

// gerund renders the "instead of ..." phrase for a replacement, e.g. "gaining
// Æmber from reaping".
func (e Event) gerund() string {
	return "gaining Æmber from reaping"
}

// lastingAction is what a lasting effect does when it fires or replaces.
type lastingAction uint8

const (
	actGainAember lastingAction = iota
	actDealDamage
	actSteal
)

// describe is a short label for a reaction, used when the controller orders several
// that fire at once.
func (a lastingAction) describe() string {
	if a == actDealDamage {
		return "deal damage"
	}
	return "gain Æmber"
}

// LastingEffect is one registered lasting effect: the event it attaches to, what it
// does, the player it belongs to, and its magnitude. It is a plain value so the
// flat GameState can hold a fixed array of them.
type LastingEffect struct {
	On         Event
	Do         lastingAction
	Controller int8
	Amount     int8
}

// maxLasting bounds how many lasting effects can be active at once — generous for
// the handful of "remainder of the turn" cards, and a fixed size keeps the state a
// flat value.
const maxLasting = 8

// AddLasting registers a lasting effect owned by controller for the rest of their
// turn. It is the single seam a "for the remainder of the turn" effect uses instead
// of hardcoding itself into the play or reap path.
func (g *Game) AddLasting(on Event, do lastingAction, controller, amount int) {
	if int(g.State.LastingCount) >= maxLasting {
		return
	}
	g.State.Lasting[g.State.LastingCount] = LastingEffect{
		On:         on,
		Do:         do,
		Controller: int8(controller),
		Amount:     int8(amount),
	}
	g.State.LastingCount++
}

// clearLasting drops the lasting effects a player owns, called when their turn ends
// so the effects last only that turn.
func (g *Game) clearLasting(player int) {
	n := 0
	for i := 0; i < int(g.State.LastingCount); i++ {
		if int(g.State.Lasting[i].Controller) != player {
			g.State.Lasting[n] = g.State.Lasting[i]
			n++
		}
	}
	for i := n; i < int(g.State.LastingCount); i++ {
		g.State.Lasting[i] = LastingEffect{}
	}
	g.State.LastingCount = uint8(n)
}

// fireLasting resolves every reaction actor owns that responds to event. When
// several fire at once the controller chooses the order (KeyForge lets the active
// player order simultaneous triggers). subject is the card that caused the event —
// the played creature — for reactions that need a source.
func (g *Game) fireLasting(event Event, actor int, subject LocalID) {
	var pending []LastingEffect
	for i := 0; i < int(g.State.LastingCount); i++ {
		if le := g.State.Lasting[i]; int(le.Controller) == actor && le.On == event {
			pending = append(pending, le)
		}
	}
	for len(pending) > 0 {
		idx := 0
		if len(pending) > 1 {
			labels := make([]string, len(pending))
			for i, le := range pending {
				labels[i] = le.Do.describe()
			}
			idx = g.chooseOption(actor, "", "Choose the next effect to resolve", labels)
		}
		g.resolveReaction(pending[idx], actor, subject)
		pending = append(pending[:idx], pending[idx+1:]...)
	}
}

// resolveReaction resolves a single reaction for actor, using subject as the
// triggering card where one is needed.
func (g *Game) resolveReaction(le LastingEffect, actor int, subject LocalID) {
	switch le.Do {
	case actDealDamage:
		DealDamage{Amount: int(le.Amount), Target: Target{Kind: TargetChosenEnemyCreature}}.Resolve(&EffectContext{
			Resolver:   g,
			Source:     subject,
			Controller: actor,
		})
	default: // actGainAember
		g.State.Aember[actor] += int(le.Amount)
		g.logf("%s gains %d Æmber (%s)", g.names[actor], le.Amount, le.On.clause())
	}
}

// lastingReplacement returns the replacement a player has for a replacement event,
// or ok=false when none is active. It is how an event site (reapWith) asks whether
// its outcome is being replaced this turn.
func (g *Game) lastingReplacement(player int, event Event) (lastingAction, bool) {
	for i := 0; i < int(g.State.LastingCount); i++ {
		if le := g.State.Lasting[i]; int(le.Controller) == player && le.On == event {
			return le.Do, true
		}
	}
	return 0, false
}
