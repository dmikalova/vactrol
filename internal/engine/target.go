package engine

import "fmt"

// A Target names the cards an ability acts on. KeyForge abilities are written in
// terms of noun phrases — "this creature", "each enemy creature", "a friendly
// creature", "each Scientist trait creature", "each creature with power 3 or
// lower" — and Target captures exactly that: a base set chosen by Kind, narrowed
// by optional filters. Because the same value both renders the phrase (Text) and
// selects the cards (Select), one effect such as Destroy can express many
// different printed cards just by changing its Target.

// TargetKind enumerates the base sets a Target can select before filtering.
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
	// TargetEachOtherFriendlyCreature selects the controller's creatures except
	// the source card.
	TargetEachOtherFriendlyCreature
	// TargetChosenCreature selects a single creature the controller chooses from
	// all creatures in play (either player's).
	TargetChosenCreature
	// TargetChosenEnemyCreature selects a single enemy creature the controller
	// chooses.
	TargetChosenEnemyCreature
)

// Target describes which cards an effect applies to. Kind picks the base set;
// the optional filters added by WithTrait and PowerAtMost narrow that set and
// extend the rendered text.
type Target struct {
	Kind        TargetKind
	trait       Trait
	maxPower    int
	hasMaxPower bool
	damaged     bool
	onFlank     bool
	notOnFlank  bool
}

// WithTrait narrows the target to cards that have the given trait, e.g.
// Target{Kind: TargetEachCreature}.WithTrait("Scientist").
func (t Target) WithTrait(trait Trait) Target {
	t.trait = trait
	return t
}

// PowerAtMost narrows the target to creatures whose power is max or lower, e.g.
// Target{Kind: TargetEachCreature}.PowerAtMost(3).
func (t Target) PowerAtMost(max int) Target {
	t.maxPower = max
	t.hasMaxPower = true
	return t
}

// Damaged narrows the target to creatures that currently have damage on them.
func (t Target) Damaged() Target {
	t.damaged = true
	return t
}

// OnFlank narrows the target to creatures on a flank of their battleline (its
// leftmost or rightmost creature).
func (t Target) OnFlank() Target {
	t.onFlank = true
	return t
}

// NotOnFlank narrows the target to creatures that are not on a flank of their
// battleline (neither its leftmost nor rightmost creature).
func (t Target) NotOnFlank() Target {
	t.notOnFlank = true
	return t
}

// Text renders the target as an English noun phrase, e.g. "each enemy creature",
// "each Scientist trait creature", or "each creature with power 3 or lower".
func (t Target) Text() string {
	if t.Kind == TargetTriggeringCreature {
		return "it"
	}
	noun := "creature"
	if t.Kind == TargetEachArtifact {
		noun = "artifact"
	}
	if t.trait != "" {
		noun = string(t.trait) + " trait " + noun
	}
	if t.onFlank {
		noun = "flank " + noun
	}
	if t.damaged {
		noun = "damaged " + noun
	}
	var phrase string
	switch t.Kind {
	case TargetThisCreature:
		phrase = "this " + noun
	case TargetEachCreature, TargetEachArtifact:
		phrase = "each " + noun
	case TargetEachFriendlyCreature:
		phrase = "each friendly " + noun
	case TargetEachEnemyCreature:
		phrase = "each enemy " + noun
	case TargetEachOtherFriendlyCreature:
		phrase = "each other friendly " + noun
	case TargetChosenEnemyCreature:
		phrase = "an enemy " + noun
	default:
		phrase = "a " + noun
	}
	if t.hasMaxPower {
		phrase += fmt.Sprintf(" with power %d or lower", t.maxPower)
	}
	if t.notOnFlank {
		phrase += " that is not on a flank"
	}
	return phrase
}

// Select resolves the target into concrete card ids, applying its filters. For a
// chosen kind it asks the controller to pick one of the filtered candidates
// (returning nil when there are none or the choice is declined).
func (t Target) Select(ctx *EffectContext) []LocalID {
	ids := t.filter(ctx, t.selectBase(ctx))
	if !t.isChosen() {
		return ids
	}
	if len(ids) == 0 {
		return nil
	}
	id, ok := ctx.Resolver.ChooseCreature(ctx.Controller, "Choose "+t.Text(), ids)
	if !ok {
		return nil
	}
	return []LocalID{id}
}

// isChosen reports whether the Kind resolves to a single player-chosen creature.
func (t Target) isChosen() bool {
	return t.Kind == TargetChosenCreature || t.Kind == TargetChosenEnemyCreature
}

// filter narrows ids to those matching the target's trait, power, damaged, and
// flank filters.
func (t Target) filter(ctx *EffectContext, ids []LocalID) []LocalID {
	if t.trait == "" && !t.hasMaxPower && !t.damaged && !t.onFlank && !t.notOnFlank {
		return ids
	}
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if t.trait != "" && !ctx.Resolver.HasTrait(id, t.trait) {
			continue
		}
		if t.hasMaxPower && ctx.Resolver.Power(id) > t.maxPower {
			continue
		}
		if t.damaged && ctx.Resolver.Damage(id) == 0 {
			continue
		}
		if t.onFlank && !onFlank(ctx, id) {
			continue
		}
		if t.notOnFlank && onFlank(ctx, id) {
			continue
		}
		out = append(out, id)
	}
	return out
}

// onFlank reports whether a creature is on a flank of its battleline (its
// leftmost or rightmost creature).
func onFlank(ctx *EffectContext, id LocalID) bool {
	bl := ctx.Resolver.Battleline(ctx.Resolver.Owner(id))
	return len(bl) > 0 && (bl[0] == id || bl[len(bl)-1] == id)
}

// selectBase resolves the unfiltered base set chosen by Kind. Chosen kinds return
// the pool of candidates; Select applies filters and prompts for the choice.
func (t Target) selectBase(ctx *EffectContext) []LocalID {
	switch t.Kind {
	case TargetThisCreature:
		return []LocalID{ctx.Source}
	case TargetTriggeringCreature:
		if ctx.HasIt {
			return []LocalID{ctx.It}
		}
		return nil
	case TargetEachArtifact:
		return append(ctx.Resolver.Artifacts(ctx.Controller), ctx.Resolver.Artifacts(1-ctx.Controller)...)
	case TargetEachCreature, TargetChosenCreature:
		return append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Battleline(1-ctx.Controller)...)
	case TargetEachFriendlyCreature:
		return ctx.Resolver.Battleline(ctx.Controller)
	case TargetEachEnemyCreature, TargetChosenEnemyCreature:
		return ctx.Resolver.Battleline(1 - ctx.Controller)
	case TargetEachOtherFriendlyCreature:
		out := make([]LocalID, 0)
		for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
			if id != ctx.Source {
				out = append(out, id)
			}
		}
		return out
	default:
		return nil
	}
}
