package card

import (
	"sort"

	"github.com/dmikalova/vactrol/internal/game"
)

// registry holds every card built with New. A set package declares each card as
// a package-level var, `var X = card.New(...)`, so simply importing the set
// enrolls its cards at package-initialization time — there is no per-card init
// and no central list to maintain (see package cards).
var registry []Definition

// New builds a card Definition and enrolls it in the global database, returning
// it so it reads as the initializer of a package-level var:
//
//	var Anger = card.New("Anger", card.House.Brobnar, card.Type.Action, card.Rarity.Common, ...)
//
// For a throwaway card that should not join the database (tests, scripted demos)
// use the engine's game.NewCard directly.
func New(name string, house game.House, ct game.CardType, rarity game.Rarity, opts ...Option) Definition {
	d := game.NewCard(name, house, ct, rarity, opts...)
	registry = append(registry, d)
	return d
}

// Registered returns a copy of every registered card, sorted by house then name
// so the database is deterministic regardless of package initialization order.
// The cards aggregator (package cards) imports the set packages that populate the
// registry, so callers should reach it through cards.All().
func Registered() []Definition {
	out := make([]Definition, len(registry))
	copy(out, registry)
	sort.Slice(out, func(i, j int) bool {
		if out[i].House != out[j].House {
			return out[i].House < out[j].House
		}
		return out[i].Name < out[j].Name
	})
	return out
}
