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
)

// Target describes which cards an effect applies to. Kind picks the base set;
// the optional filters added by WithTrait and PowerAtMost narrow that set and
// extend the rendered text.
type Target struct {
	Kind        TargetKind
	trait       Trait
	maxPower    int
	hasMaxPower bool
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
	default:
		phrase = "a " + noun
	}
	if t.hasMaxPower {
		phrase += fmt.Sprintf(" with power %d or lower", t.maxPower)
	}
	return phrase
}

// Select resolves the target into concrete card ids in the current game state,
// applying any trait or power filters on top of the base set chosen by Kind.
func (t Target) Select(ctx *EffectContext) []LocalID {
	ids := t.selectBase(ctx)
	if t.trait == "" && !t.hasMaxPower {
		return ids
	}
	filtered := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if t.trait != "" && !ctx.Resolver.HasTrait(id, t.trait) {
			continue
		}
		if t.hasMaxPower && ctx.Resolver.Power(id) > t.maxPower {
			continue
		}
		filtered = append(filtered, id)
	}
	return filtered
}

// selectBase resolves the unfiltered base set chosen by Kind.
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
	case TargetEachCreature:
		return append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Battleline(1-ctx.Controller)...)
	case TargetEachFriendlyCreature:
		return ctx.Resolver.Battleline(ctx.Controller)
	case TargetEachEnemyCreature:
		return ctx.Resolver.Battleline(1 - ctx.Controller)
	case TargetEachOtherFriendlyCreature:
		out := make([]LocalID, 0)
		for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
			if id != ctx.Source {
				out = append(out, id)
			}
		}
		return out
	case TargetChosenCreature:
		cands := append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Battleline(1-ctx.Controller)...)
		if len(cands) == 0 {
			return nil
		}
		id, ok := ctx.Resolver.ChooseCreature(ctx.Controller, "Choose "+t.Text(), cands)
		if !ok {
			return nil
		}
		return []LocalID{id}
	default:
		return nil
	}
}
