package engine

import "fmt"

// ReadyAndBelongToHouseAfterYouPlayCreature is the effect granted by Brain Stem
// Antenna. The generic after-creature-enters trigger fires for both players and
// every house; this effect narrows it to a creature of House played by the
// controller, then readies the source creature and changes the house it belongs to
// for the rest of that controller's turn.
//
//rulebook:effect Belong to House
type ReadyAndBelongToHouseAfterYouPlayCreature struct {
	// House is both the house the triggering creature must belong to and the house
	// the source creature belongs to for the remainder of the turn.
	House House
}

// validate requires the house named by the ability.
func (e ReadyAndBelongToHouseAfterYouPlayCreature) validate() error {
	if e.House == HouseNone {
		return fmt.Errorf("ReadyAndBelongToHouseAfterYouPlayCreature: house must be set")
	}
	return nil
}

// Text renders Brain Stem Antenna's granted ability text.
func (e ReadyAndBelongToHouseAfterYouPlayCreature) Text() string {
	return fmt.Sprintf("after you play a %s creature, ready %s and for the remainder of the turn it belongs to house %s", e.House, SelfName, e.House)
}

// abilityText lets this effect supply the whole triggered ability line instead of
// inheriting TriggerAfterCreatureEnters' broader "After a creature enters play"
// prefix.
func (e ReadyAndBelongToHouseAfterYouPlayCreature) abilityText(trigger Trigger) (string, bool) {
	if trigger != TriggerAfterCreatureEnters {
		return "", false
	}
	return e.Text(), true
}

// Resolve ignores off-house creatures and creatures the opponent played. When the
// trigger matches, it readies the source creature and changes the house it belongs
// to until its controller's turn ends.
func (e ReadyAndBelongToHouseAfterYouPlayCreature) Resolve(ctx *EffectContext) {
	if !ctx.HasIt || ctx.Resolver.Owner(ctx.It) != ctx.Controller || ctx.Resolver.House(ctx.It) != e.House {
		return
	}
	ctx.Resolver.SetExhausted(ctx.Source, false)
	ctx.Resolver.BelongToHouseForRemainderOfTurn(ctx.Source, e.House)
}
