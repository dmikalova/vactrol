package engine

import "errors"

// Engine action errors.
var (
	ErrNotActivePlayer       = errors.New("not the active player")
	ErrCardNotInHand         = errors.New("card index is not in hand")
	ErrWrongType             = errors.New("card is the wrong type for this action")
	ErrWrongHouse            = errors.New("card's house is not the active house")
	ErrCardExhausted         = errors.New("card is exhausted")
	ErrNoTarget              = errors.New("no legal target")
	ErrGameOver              = errors.New("the game is already over")
	ErrCannotFight           = errors.New("cannot use creatures to fight this turn")
	ErrCannotPlayCreature    = errors.New("cannot play creatures")
	ErrCardPlayLimit         = errors.New("card-play limit reached this turn")
	ErrMustChooseForcedHouse = errors.New("must choose the forced active house this turn")
)

// This file holds the turn lifecycle — begin (with the forge step), choose the
// active house, and end (ready cards, refresh armor, draw back up) — together
// with the key-forging step that decides the game.

// BeginTurn starts a player's turn: it becomes the active player, the active
// house is cleared, and the forge step runs (keys are forged if affordable).
func (g *Game) BeginTurn(player int) {
	if g.State.Winner >= 0 {
		return
	}
	g.State.ActivePlayer = player
	g.State.ActiveHouse = HouseNone
	g.State.Turn++
	g.State.CardsPlayedThisTurn[player] = 0
	// A fight bar armed on a previous turn becomes active for this player now.
	if g.State.CannotFightNext[player] {
		g.State.CannotFight[player] = true
		g.State.CannotFightNext[player] = false
	}
	// A forced active house armed on a previous turn takes effect for this player now.
	g.State.ForcedHouse[player] = g.State.ForcedHouseNext[player]
	g.State.ForcedHouseNext[player] = HouseNone
	g.logf("--- %s begins turn %d ---", g.names[player], g.State.Turn)
	g.forgeKey(player)
}

// Choose a house: pick one of your deck's three houses to be your active house
// for this turn. For the rest of the turn you may play from hand and use only
// cards of that house, except cards that ignore the restriction such as Versatile
// ones.
//
//rulebook:turn Turn structure / 2. Choose a house

// ChooseHouse sets ActiveHouse, then offers to draw the player's archived cards
// into hand now that a house is locked in.
func (g *Game) ChooseHouse(player int, house House) error {
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	if fh := g.State.ForcedHouse[player]; fh != HouseNone && house != fh {
		return ErrMustChooseForcedHouse
	}
	g.State.ActiveHouse = house
	g.logf("%s chooses house %s", g.names[player], house)
	g.offerArchives(player)
	return nil
}

// Ready and draw: at the end of your turn every card you control readies (turns
// back upright, ready to act next turn) and you draw back up to a full hand of
// six cards. Creatures and artifacts that entered play exhausted this turn ready
// here too.
//
//rulebook:turn Turn structure / 3. Ready and draw

// EndTurn clears Exhausted on the player's in-play cards, refreshes creature
// armor for the new turn, and draws up to HandSize.
func (g *Game) EndTurn(player int) {
	for _, id := range g.allInPlay(player) {
		core := &g.State.Cards[id]
		core.Exhausted = false
		if g.cat.def(id).Type == Creature {
			core.ArmorRemaining = int16(g.armor(id))
		}
	}
	g.State.CannotFight[player] = false
	g.State.MayFightHouse[player] = HouseNone
	g.clearLasting(player)
	g.drawTo(player, HandSize)
	g.logf("%s ends their turn", g.names[player])
}

// CannotFightNextTurn arms a fight bar on a player for their next turn.
func (g *Game) CannotFightNextTurn(player int) {
	g.State.CannotFightNext[player] = true
}

// GrantFightForHouse lets a player use creatures of house h to fight this turn
// even when h is not the active house. EndTurn clears the grant.
func (g *Game) GrantFightForHouse(player int, h House) {
	g.State.MayFightHouse[player] = h
	g.logf("%s's %s creatures may fight this turn", g.names[player], h)
}

// ForceActiveHouseNextTurn makes a player have to choose house h as their active
// house on their next turn (Control the Weak). BeginTurn promotes the armed house.
func (g *Game) ForceActiveHouseNextTurn(player int, h House) {
	g.State.ForcedHouseNext[player] = h
	g.logf("%s must choose house %s next turn", g.names[player], h)
}

// Forge a key: at the start of your turn you forge a single key if you can pay
// its current cost — 6 Æmber by default. A player forges at most one key per turn.
// Keys are the win condition — forge your third key and you win the game.
//
//rulebook:turn Turn structure / 1. Forge a key

// forgeKey forges one key when the player can afford the current key cost, paying
// it and firing "after you forge a key" abilities. BeginTurn forges at most one
// key at the start of a turn; cards may forge one more via the ForgeKey effect.
func (g *Game) forgeKey(player int) {
	cost := g.keyCost(player)
	if g.State.Aember[player] < cost {
		return
	}
	g.State.Aember[player] -= cost
	g.State.Keys[player]++
	g.logf("%s forges a key (%d/%d)", g.names[player], g.State.Keys[player], KeysToWin)
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerAfterForgeKey, 0, false)
	}
	if g.State.Keys[player] >= KeysToWin {
		g.State.Winner = player
		g.logf("%s wins the game!", g.names[player])
	}
}
