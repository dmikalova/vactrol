package engine

// Playing from the top of your deck reveals that card first, then turns it into a
// normal card play only when it matches the required house.
//
//rulebook:effect Play Top of Deck
type PlayTopOfDeckOfChosenHouse struct{}

// Text renders the Chaos Portal effect after the house choice.
func (PlayTopOfDeckOfChosenHouse) Text() string {
	return "reveal the top card of your deck. If it is of the chosen house, play it"
}

// Resolve reveals the controller's top deck card and plays it if it belongs to
// the house most recently chosen by ChooseHouseThen.
func (PlayTopOfDeckOfChosenHouse) Resolve(ctx *EffectContext) {
	if ctx.ChosenHouse == HouseNone {
		return
	}
	ctx.Resolver.PlayTopOfDeckIfHouse(ctx.Controller, ctx.ChosenHouse)
}

// DiscardTopOfEachDeckAndDestroyByHouse turns the discarded deck cards into a
// pair of house-bound destruction choices. The controller discards their top card
// first, then their opponent's; for each discarded card, they choose one creature
// or artifact in play of that card's house to destroy. If those choices destroy
// fewer than two cards, the source artifact destroys itself.
type DiscardTopOfEachDeckAndDestroyByHouse struct{}

// Text renders the Bonkers Killing Machine action as its three printed steps.
func (DiscardTopOfEachDeckAndDestroyByHouse) Text() string {
	return "discard the top card of each player's deck. For each card discarded this way, destroy a creature or artifact of that card's house. If fewer than 2 cards are destroyed this way, destroy " + SelfName
}

// Resolve discards the controller's top deck card, then the opponent's, using
// each discarded card's house to choose a matching creature or artifact in play.
func (DiscardTopOfEachDeckAndDestroyByHouse) Resolve(ctx *EffectContext) {
	destroyed := 0
	for _, player := range []int{ctx.Controller, ctx.Opponent()} {
		discarded, ok := ctx.Resolver.DiscardTopOfDeck(player)
		if !ok {
			continue
		}
		house := ctx.Resolver.House(discarded)
		chosen, ok := ctx.Resolver.ChooseCreature(ctx.Controller, ctx.Source, "Choose a "+house.String()+" creature or artifact", inPlayOfHouse(ctx, house))
		if !ok {
			continue
		}
		if destroyAndReport(ctx, chosen) {
			destroyed++
		}
	}
	if destroyed < 2 && resolverInPlay(ctx, ctx.Source) {
		ctx.Resolver.DestroyEach(ctx.Controller, []LocalID{ctx.Source})
	}
}

// DiscardTopOfDeckAndCancelFightIfActiveHouse is a "Before Fight" effect:
// discard the attacker's controller's top deck card, and if that card belongs to
// the active house, the fight does not occur. The attacker was still used to
// fight, so it remains exhausted; combat reads the cancellation and skips Assault,
// Hazardous, fight damage, and Fight: abilities.
type DiscardTopOfDeckAndCancelFightIfActiveHouse struct{}

// Text renders Evasion Sigil's granted "Before Fight" ability.
func (DiscardTopOfDeckAndCancelFightIfActiveHouse) Text() string {
	return "discard the top card of its controller's deck. If the discarded card is of the active house, the fight does not occur"
}

// Resolve discards the source creature controller's top deck card, then cancels
// the current fight when the discarded card's house is the active house.
func (DiscardTopOfDeckAndCancelFightIfActiveHouse) Resolve(ctx *EffectContext) {
	discarded, ok := ctx.Resolver.DiscardTopOfDeck(ctx.Controller)
	if ok && ctx.Resolver.House(discarded) == ctx.Resolver.ActiveHouse() {
		ctx.Resolver.CancelCurrentFight()
	}
}

// inPlayOfHouse returns every creature and artifact in play of house, from the
// controller's point of view: friendly battleline, enemy battleline, friendly
// artifacts, then enemy artifacts.
func inPlayOfHouse(ctx *EffectContext, house House) []LocalID {
	ids := ctx.Resolver.Battleline(ctx.Controller)
	ids = append(ids, ctx.Resolver.Battleline(ctx.Opponent())...)
	ids = append(ids, ctx.Resolver.Artifacts(ctx.Controller)...)
	ids = append(ids, ctx.Resolver.Artifacts(ctx.Opponent())...)
	var out []LocalID
	for _, id := range ids {
		if ctx.Resolver.House(id) == house {
			out = append(out, id)
		}
	}
	return out
}

// destroyAndReport destroys id and reports whether it actually left play.
func destroyAndReport(ctx *EffectContext, id LocalID) bool {
	before := resolverInPlay(ctx, id)
	ctx.Resolver.DestroyEach(ctx.Controller, []LocalID{id})
	return before && !resolverInPlay(ctx, id)
}

// resolverInPlay reports whether id appears in either player's battleline or
// artifact row using only Resolver reads.
func resolverInPlay(ctx *EffectContext, id LocalID) bool {
	for p := 0; p < 2; p++ {
		for _, candidate := range ctx.Resolver.Battleline(p) {
			if candidate == id {
				return true
			}
		}
		for _, candidate := range ctx.Resolver.Artifacts(p) {
			if candidate == id {
				return true
			}
		}
	}
	return false
}
