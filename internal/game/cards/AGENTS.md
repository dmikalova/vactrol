# Card authoring guide

This directory holds the card database. Package `cards` is the aggregator; each
set lives in its own subpackage (e.g. `callofthearchons/`) and every card is one
self-registering file.

## File layout

- One card per file, named `snake_case.go` after the card (e.g. `dust_imp.go`,
  `ammonia_clouds.go`). The matching test is `snake_case_test.go`.
- Author cards through the `card` facade only
  (`github.com/dmikalova/vactrol/internal/game/card`). Use the grouped
  namespaces — `card.House.X`, `card.Type.X`, `card.Rarity.X`,
  `card.Keyword.X`, `card.Trigger.X`, `card.Controller.X`, `card.Target.X` —
  never the raw engine package.
- `card.New(...)` self-registers the card (via `init`); there is no central
  list to update. Adding a card is just adding a file.

## Style

Each card is a `var` whose doc comment is exactly three kinds of line, and whose
`card.New(...)` call puts **every parameter on its own line**:

```go
// Ammonia Clouds
//
//	Mars / Action / Common
//	Play: Deal 3 damage to each creature.
var AmmoniaClouds = card.New(
	"Ammonia Clouds",
	card.House.Mars,
	card.Type.Action,
	card.Rarity.Common,
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Target: card.Target.EachCreature,
		}),
)
```

Doc comment (tab-indented so godoc renders the stat/text lines as preformatted;
`gofmt` requires the blank `//` line after the title before that block):

1. `// <Card Name>` — the card's display name.
2. `//\t<House> / <Type> / <Rarity>` followed by any extra stats in this order:
   `Power`, `Armor`, `Æmber`, traits, keywords
   (e.g. `Brobnar / Creature / Rare / 5 Power / 0 Armor / Giant`).
3. `//\t<rules text>` — the printed card text. Omit this line for vanilla cards
   that have no rules text.

`card.New(...)`:

- The four positional arguments — name, `card.House.X`, `card.Type.X`,
  `card.Rarity.X` — each on their own line.
- Each `card.With*` option on its own line, in this order: `WithPower`,
  `WithArmor`, `WithAemberBonus`, `WithTraits`, `WithKeywords`, `WithStatic`,
  `WithAbility`.
- `card.WithAbility(` breaks onto the next line; the trigger and effect share a
  line (`card.Trigger.Play, card.DealDamage{`), and effect struct fields each go
  on their own line when there is more than one. Single-field effects stay
  inline (e.g. `card.GainAember{Amount: 1}`).

Run `gofmt -w .` after editing.

## Tests

- Every card has its own `snake_case_test.go` in the same set package — one test
  per card, even for vanilla creatures (a small `PlayCreature` + `Power` check
  that duplicates a little coverage is fine; it documents the card and guards
  its stats). Card `var`s initialize at package load, so an untested card still
  reports as "covered" — a dedicated test is about intent, not the number.
- Shared setup lives in the `cardtest` package
  (`internal/game/cards/cardtest`), imported by every set's tests so helpers are
  written once rather than copied per set:
  - `cardtest.Started(t, house)` — a new game with player 0 active and `house`
    chosen, ready to play that player's cards of that house.
  - `cardtest.Vanilla(name, house, power)` — a plain no-ability creature to use
    as a supporting body or opponent.
- Set test packages import only `game` + `cardtest`; they use the public engine
  API and exported card `var`s (no reaching into engine internals).
