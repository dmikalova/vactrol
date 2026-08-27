# Card authoring guide

This directory holds the card database. Package `cards` (`cards.go`) is the
aggregator; each set lives in its own subpackage under `sets/`
(e.g. `sets/callofthearchons/`) and every card is one self-registering file.
Shared test helpers live in the sibling `cardtest/` package — it is a test
helper (stdlib `httptest`/`iotest` convention), not a set, so it sits outside
`sets/`.

## File layout

- One card per file, named `snake_case.go` after the card (e.g. `dust_imp.go`,
  `ammonia_clouds.go`). The matching test is `snake_case_test.go`.
- Author cards through the `card` facade only
  (`github.com/dmikalova/vactrol/internal/card`). Use the grouped
  namespaces — `card.House.X`, `card.Type.X`, `card.Rarity.X`,
  `card.Keyword.X`, `card.Trigger.X`, `card.Target.X` (and the relative-player
  values `card.Controller` / `card.Opponent`) —
  never the raw engine package.
- `card.New(...)` self-registers the card (via `init`); there is no central
  list to update. Adding a card is just adding a file.
- Tag each card's origin with `card.Provenance(card.<Set>, <number>)` (e.g.
  `card.Provenance(card.CotA, 1)`) as the first option — this links it to the
  original KeyForge card it derives from and drives the TUI Provenance coverage
  view. Source-set aliases (`card.CotA`, `card.DT`, …) and the catalogs live in
  `internal/cards/provenance`. Optional and repeatable; omit for wholly original
  cards.

## Style

The doc comment above each card is **generated, not hand-written** — run
`make generate-comments` (it rewrites every card's comment from its definition,
the same details the TUI card box shows via `engine.RenderCardText`). Write the
`card.New(...)` call and let the generator fill in the comment:

```go
// Ammonia Clouds
//
//	House:  Mars
//	Type:   Action
//	Rarity: Common
//
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

The generated comment (tab-indented so godoc renders the block preformatted;
`gofmt` requires the blank `//` line after the title) is:

1. `// <Card Name>` — the card's display name.
2. A labeled, colon-aligned block: `House`, `Type`, `Rarity`, then `Power` and
   `Armor` for creatures, `Æmber` when the card has an Æmber bonus, and `Traits`
   when it has traits.
3. After a blank `//`, the printed rules text (keywords, upgrade modifier, and
   ability lines) — omitted entirely for vanilla cards with no rules text.

So authoring a card is just adding the `card.New(...)` file and running
`make generate-comments`; a placeholder `// <Name>` line above the var is enough
to seed it. `card.New(...)`:

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
  (`internal/cards/cardtest`), imported by every set's tests so helpers are
  written once rather than copied per set:
  - `cardtest.Started(t, house)` — a new game with player 0 active and `house`
    chosen, ready to play that player's cards of that house.
  - `cardtest.Vanilla(name, house, power)` — a plain no-ability creature to use
    as a supporting body or opponent.
- Set test packages import only `engine` + `cardtest`; they use the public engine
  API and exported card `var`s (no reaching into engine internals).
