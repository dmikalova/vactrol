package engine

import (
	"fmt"
	"math/rand"
)

// Win/economy constants.
const (
	// KeyCost is the amount of Æmber required to forge one key.
	KeyCost = 6
	// KeysToWin is the number of keys a player must forge to win.
	KeysToWin = 3
	// HandSize is the number of cards a player draws back up to at end of turn.
	HandSize = 6
)

// Chooser makes target decisions for a player. The engine calls it whenever an
// effect must pick a creature. Implementations must be deterministic so games can
// be reproduced from a seed.
type Chooser interface {
	// ChooseCreature returns one id from candidates and true, or false if none.
	ChooseCreature(prompt string, candidates []LocalID) (LocalID, bool)
}

// FirstChooser always picks the first available candidate. It is the default and
// keeps behavior deterministic for tests and simulation.
type FirstChooser struct{}

// ChooseCreature returns the first candidate, or false if the list is empty.
func (FirstChooser) ChooseCreature(_ string, candidates []LocalID) (LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	return candidates[0], true
}

// Game bundles the flat GameState with the read-only Catalog and the surrounding
// engine services (player names, choosers, RNG, log). Cloning a state for MCTS
// only needs GameState.FastCopy; this wrapper is the live match harness.
type Game struct {
	State   GameState
	Verbose bool
	Log     []string

	names    [2]string
	choosers [2]Chooser
	cat      *catalog
	rng      *rand.Rand
}

// NewGame creates a new two-player game seeded for deterministic play.
func NewGame(p0Name, p1Name string, seed int64) *Game {
	g := &Game{
		names:    [2]string{p0Name, p1Name},
		choosers: [2]Chooser{FirstChooser{}, FirstChooser{}},
		cat:      &catalog{},
		rng:      rand.New(rand.NewSource(seed)),
	}
	g.State.Winner = -1
	return g
}

// SetChooser installs a custom chooser for a player (nil resets to the default).
func (g *Game) SetChooser(player int, c Chooser) { g.choosers[player] = c }

// PlayerName returns a player's display name.
func (g *Game) PlayerName(player int) string { return g.names[player] }

// logf appends a line to the game log (and prints it when Verbose).
func (g *Game) logf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	g.Log = append(g.Log, line)
	if g.Verbose {
		fmt.Println(line)
	}
}

// chooserFor returns the chooser for a player, defaulting to FirstChooser.
func (g *Game) chooserFor(player int) Chooser {
	if ch := g.choosers[player]; ch != nil {
		return ch
	}
	return FirstChooser{}
}

// orderByChoice asks controller to arrange ids into a resolution order by picking
// the next one repeatedly (the final id is forced, so it is never prompted). With
// 0 or 1 ids there is nothing to order and ids is returned unchanged; a rejected
// pick falls back to the remaining order. The default FirstChooser keeps the
// original order, so ordering only becomes interactive under a real UI.
func (g *Game) orderByChoice(controller int, prompt string, ids []LocalID) []LocalID {
	if len(ids) <= 1 {
		return ids
	}
	remaining := append([]LocalID(nil), ids...)
	ordered := make([]LocalID, 0, len(ids))
	for len(remaining) > 1 {
		chosen, ok := g.chooserFor(controller).ChooseCreature(prompt, remaining)
		if !ok {
			break
		}
		ordered = append(ordered, chosen)
		for i, id := range remaining {
			if id == chosen {
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	return append(ordered, remaining...)
}

// ---- registration & setup ----

// Register adds a definition to the catalog for an owner and returns its id.
func (g *Game) Register(def CardDefinition, owner int) LocalID {
	d := def
	return g.cat.add(&d, owner)
}

// AddToHand registers a card and places it in a player's hand.
func (g *Game) AddToHand(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Hand[owner].add(id)
	return id
}

// AddToDeck registers a card and places it on the bottom of a player's deck.
func (g *Game) AddToDeck(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Deck[owner].add(id)
	return id
}

// AddToBattleline registers a creature and places it on a player's battleline.
func (g *Game) AddToBattleline(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Cards[id].ArmorRemaining = int16(def.Armor)
	g.State.Battleline[owner].add(id)
	return id
}

// AddArtifact registers an artifact and places it in a player's artifact row.
func (g *Game) AddArtifact(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Artifacts[owner].add(id)
	return id
}

// Shuffle randomizes a player's deck using the game's seeded RNG.
func (g *Game) Shuffle(player int) {
	d := &g.State.Deck[player]
	for i := int(d.Count) - 1; i > 0; i-- {
		j := g.rng.Intn(i + 1)
		d.IDs[i], d.IDs[j] = d.IDs[j], d.IDs[i]
	}
}

// ---- read accessors (used by callers, effects, and tests) ----

// Def returns the read-only definition for an id.
func (g *Game) Def(id LocalID) *CardDefinition { return g.cat.def(id) }

// owner returns the owning player index for an id.
func (g *Game) owner(id LocalID) int { return g.cat.owner(id) }

// Name returns a card's printed name.
func (g *Game) Name(id LocalID) string { return g.cat.def(id).Name }

// Power returns a creature's current power including attached upgrades.
func (g *Game) Power(id LocalID) int {
	core := &g.State.Cards[id]
	p := g.cat.def(id).Power
	for i := 0; i < int(core.UpgradeCount); i++ {
		p += g.cat.def(core.Upgrades[i]).Static.PowerBonus
	}
	return p
}

// armor returns a creature's armor value including attached upgrades.
func (g *Game) armor(id LocalID) int {
	core := &g.State.Cards[id]
	a := g.cat.def(id).Armor
	for i := 0; i < int(core.UpgradeCount); i++ {
		a += g.cat.def(core.Upgrades[i]).Static.ArmorBonus
	}
	return a
}

// assault returns a creature's Assault value including attached upgrades.
func (g *Game) assault(id LocalID) int {
	core := &g.State.Cards[id]
	a := g.cat.def(id).Assault
	for i := 0; i < int(core.UpgradeCount); i++ {
		a += g.cat.def(core.Upgrades[i]).Static.AssaultBonus
	}
	return a
}

// hazardous returns a creature's Hazardous value including attached upgrades.
func (g *Game) hazardous(id LocalID) int {
	core := &g.State.Cards[id]
	h := g.cat.def(id).Hazardous
	for i := 0; i < int(core.UpgradeCount); i++ {
		h += g.cat.def(core.Upgrades[i]).Static.HazardousBonus
	}
	return h
}

// Damage returns the damage currently on a creature.
func (g *Game) Damage(id LocalID) int { return int(g.State.Cards[id].Damage) }

// AmberOn returns the Æmber sitting on a card (placed by exalt, capture, etc.).
func (g *Game) AmberOn(id LocalID) int { return int(g.State.Cards[id].Amber) }

// Exhausted reports whether a card is exhausted.
func (g *Game) Exhausted(id LocalID) bool { return g.State.Cards[id].Exhausted }

// Aember returns a player's Æmber pool.
func (g *Game) Aember(player int) int { return g.State.Aember[player] }

// Keys returns a player's forged key count.
func (g *Game) Keys(player int) int { return g.State.Keys[player] }

// Winner returns the winning player index, or -1 if the game is ongoing.
func (g *Game) Winner() int { return g.State.Winner }

// Hand returns a copy of the ids in a player's hand.
func (g *Game) Hand(player int) []LocalID { return cloneIDs(g.State.Hand[player].slice()) }

// Battleline returns a copy of the ids on a player's battleline.
func (g *Game) Battleline(player int) []LocalID {
	return cloneIDs(g.State.Battleline[player].slice())
}

// Discard returns a copy of the ids in a player's discard pile.
func (g *Game) Discard(player int) []LocalID { return cloneIDs(g.State.Discard[player].slice()) }

// Artifacts returns a copy of the ids in a player's artifact row.
func (g *Game) Artifacts(player int) []LocalID { return cloneIDs(g.State.Artifacts[player].slice()) }

// Upgrades returns the ids of upgrades attached to a creature, in attach order.
func (g *Game) Upgrades(id LocalID) []LocalID {
	core := &g.State.Cards[id]
	out := make([]LocalID, core.UpgradeCount)
	for i := range out {
		out[i] = core.Upgrades[i]
	}
	return out
}

// inPlay reports whether an id is on its owner's battleline or artifact row.
func (g *Game) inPlay(id LocalID) bool {
	o := g.owner(id)
	return g.State.Battleline[o].contains(id) || g.State.Artifacts[o].contains(id)
}

// battlelineCopy returns a fresh slice of a player's battleline ids, safe to hold
// across state mutations (e.g. while dealing damage to each creature).
func (g *Game) battlelineCopy(player int) []LocalID {
	return cloneIDs(g.State.Battleline[player].slice())
}

// allInPlay returns a fresh slice of a player's creatures and artifacts.
func (g *Game) allInPlay(player int) []LocalID {
	b := g.State.Battleline[player].slice()
	a := g.State.Artifacts[player].slice()
	out := make([]LocalID, 0, len(b)+len(a))
	out = append(out, b...)
	out = append(out, a...)
	return out
}

// cloneIDs copies a slice of ids so callers cannot alias the state arrays.
func cloneIDs(src []LocalID) []LocalID {
	out := make([]LocalID, len(src))
	copy(out, src)
	return out
}
