package engine

import (
	"fmt"
	"strings"
)

// ReturnFromDiscard returns a card the controller chooses from their own discard
// zone to a Destination — their hand (the default) or the top of their deck. Type
// restricts the choice to cards of that type; the zero value allows any card. This
// is how cards recur from the discard zone, e.g. "Put a creature from your discard
// zone on top of your deck."
type ReturnFromDiscard struct {
	// Type restricts the choice to cards of that type; the zero value (an unset
	// CardType) allows any card.
	Type CardType
	// Destination is where the chosen card goes: ToHand (the zero value/default) or
	// ToTopOfDeck. The other deck destinations are not supported yet.
	Destination Destination
}

// noun renders the kind of card the effect returns — the lowercased card type when
// Type is set (e.g. "creature"), otherwise the generic "card".
func (e ReturnFromDiscard) noun() string {
	if e.Type != "" {
		return strings.ToLower(string(e.Type))
	}
	return "card"
}

// validate rejects a Destination this effect cannot move a card to; only the hand
// and the top of the deck are supported.
func (e ReturnFromDiscard) validate() error {
	if e.Destination != ToHand && e.Destination != ToTopOfDeck {
		return fmt.Errorf("ReturnFromDiscard: unsupported destination %d", e.Destination)
	}
	return nil
}

// Text renders the effect, e.g. "put a card from your discard zone into your
// hand" or "put a creature from your discard zone on top of your deck".
func (e ReturnFromDiscard) Text() string {
	dest := "into your hand"
	if e.Destination == ToTopOfDeck {
		dest = "on top of your deck"
	}
	return "put a " + e.noun() + " from your discard zone " + dest
}

// Resolve lets the controller choose an eligible card from their discard zone and
// moves it to the chosen destination. It does nothing if there is no candidate or
// the choice is declined.
func (e ReturnFromDiscard) Resolve(ctx *EffectContext) {
	discard := ctx.Resolver.Discard(ctx.Controller)
	candidates := discard
	if e.Type != "" {
		candidates = nil
		for _, id := range discard {
			if ctx.Resolver.TypeOf(id) == e.Type {
				candidates = append(candidates, id)
			}
		}
	}
	id, ok := ctx.ChooseCreature("Choose a "+e.noun()+" from your discard zone", candidates)
	if !ok {
		return
	}
	if e.Destination == ToTopOfDeck {
		ctx.Resolver.ReturnFromDiscardToTopOfDeck(id)
	} else {
		ctx.Resolver.ReturnFromDiscardToHand(id)
	}
}

// DiscardHand discards cards from a player's hand: the chosen player's cards,
// optionally only creatures and only those of the house picked by an enclosing
// ChooseHouseThen. It models "discard each creature of the chosen house from your
// opponent's hand."
type DiscardHand struct {
	Player        Player
	CreaturesOnly bool
	OfChosenHouse bool
}

// validate rejects a DiscardHand whose player was left unset.
func (e DiscardHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardHand")
	}
	return nil
}

// Text renders the effect, e.g. "discard each creature of the chosen house from
// your opponent's hand".
func (e DiscardHand) Text() string {
	what := "each card"
	if e.CreaturesOnly {
		what = "each creature"
	}
	if e.OfChosenHouse {
		what += " of the chosen house"
	}
	whose := "your hand"
	if e.Player == Opponent {
		whose = "your opponent's hand"
	}
	return "discard " + what + " from " + whose
}

// Resolve discards every matching card from the chosen player's hand.
func (e DiscardHand) Resolve(ctx *EffectContext) {
	owner := ctx.PlayerFor(e.Player)
	for _, id := range ctx.Resolver.Hand(owner) {
		if e.CreaturesOnly && !ctx.Resolver.IsCreature(id) {
			continue
		}
		if e.OfChosenHouse && ctx.Resolver.House(id) != ctx.ChosenHouse {
			continue
		}
		ctx.Resolver.DiscardCardFromHand(owner, id)
	}
}
