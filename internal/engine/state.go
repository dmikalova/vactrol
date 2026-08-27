package engine

// LocalID identifies one physical card within a single match. During setup every
// card (including duplicates) is registered in the Catalog and assigned a stable
// LocalID; the flat GameState references cards only by this id.
type LocalID uint8

// Capacities for the fixed-size, pointerless state arrays. They are generous
// upper bounds — a real match uses far fewer.
const (
	// maxCards is the LocalID space for one match.
	maxCards = 128
	// zoneCap is the capacity of a single zone (deck is the largest zone).
	zoneCap = 40
	// maxUpgrades is the most upgrades that can be stacked on one creature.
	maxUpgrades = 6
)

// CardCore is the mutable per-match state of a single card, stored purely by
// value. It carries no pointers so the whole GameState copies flat.
type CardCore struct {
	Exhausted      bool
	Stunned        bool
	Damage         int16
	ArmorRemaining int16
	// Amber is Æmber sitting on the card (e.g. placed by exalt or capture). It
	// belongs to no player's pool while it stays here.
	Amber        int16
	UpgradeCount uint8
	Upgrades     [maxUpgrades]LocalID
}

// Zone is an ordered, fixed-capacity collection of card ids (hand, deck, battle
// line, etc.). It is a value type: copying a Zone copies its contents.
type Zone struct {
	IDs   [zoneCap]LocalID
	Count uint8
}

// slice returns the live ids as a slice header into the underlying array. The
// result must be treated as read-only and not retained across mutations.
func (z *Zone) slice() []LocalID { return z.IDs[:z.Count] }

// add appends an id to the end of the zone.
func (z *Zone) add(id LocalID) {
	z.IDs[z.Count] = id
	z.Count++
}

// addFront inserts an id at the front of the zone (the left flank / top of deck).
func (z *Zone) addFront(id LocalID) {
	copy(z.IDs[1:z.Count+1], z.IDs[:z.Count])
	z.IDs[0] = id
	z.Count++
}

// indexOf returns the position of id, or -1 if absent.
func (z *Zone) indexOf(id LocalID) int {
	for i := 0; i < int(z.Count); i++ {
		if z.IDs[i] == id {
			return i
		}
	}
	return -1
}

// contains reports whether the zone holds id.
func (z *Zone) contains(id LocalID) bool { return z.indexOf(id) >= 0 }

// removeAt removes the id at position i, preserving order, and returns it.
func (z *Zone) removeAt(i int) LocalID {
	id := z.IDs[i]
	copy(z.IDs[i:], z.IDs[i+1:z.Count])
	z.Count--
	z.IDs[z.Count] = 0
	return id
}

// remove deletes id from the zone if present, reporting whether it was found.
func (z *Zone) remove(id LocalID) bool {
	i := z.indexOf(id)
	if i < 0 {
		return false
	}
	z.removeAt(i)
	return true
}

// GameState is the complete mutable state of a match, laid out as a flat value.
// It contains no pointers, slices, or maps, so a copy is a pure value copy with
// no heap allocation or garbage-collector pressure — the property MCTS rollouts
// depend on. Read-only card definitions live in the separate catalog.
type GameState struct {
	Cards      [maxCards]CardCore
	Battleline [2]Zone
	Hand       [2]Zone
	Deck       [2]Zone
	Discard    [2]Zone
	Artifacts  [2]Zone

	Aember [2]int
	Keys   [2]int

	ActivePlayer int
	ActiveHouse  House
	Turn         int
	Winner       int // -1 while the game is ongoing
}

// FastCopy returns an independent copy of the state. Because every field is a
// value type this is a single flat copy; mutating the result never affects the
// original.
func (s GameState) FastCopy() GameState { return s }

// catalog is the read-only registry of card definitions for a match. It is held
// separately from GameState (by pointer) and never mutated during play, so it is
// shared freely across cloned states.
type catalog struct {
	defs   []*CardDefinition
	owners []uint8
}

// add registers a definition for an owner and returns its assigned LocalID.
func (c *catalog) add(def *CardDefinition, owner int) LocalID {
	id := LocalID(len(c.defs))
	c.defs = append(c.defs, def)
	c.owners = append(c.owners, uint8(owner))
	return id
}

// def returns the definition for an id.
func (c *catalog) def(id LocalID) *CardDefinition { return c.defs[id] }

// owner returns the owning player index for an id.
func (c *catalog) owner(id LocalID) int { return int(c.owners[id]) }
