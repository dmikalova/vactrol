# Vactrol

Vactrol is a rules engine for a [KeyForge](https://keyforge.fandom.com/)-style
card game, written in Go. It models a full two-player match — houses, Æmber,
keys, combat, and card abilities — behind a small, pointerless state type that is
cheap to copy (for AI search), and ships with a terminal UI for browsing the card
database and playing hotseat games.

## Quick start

```sh
make run     # launch the terminal UI (card explorer + hotseat game)
make test    # run the test suite
make cover   # engine test coverage (kept at 100%)
```

Requires Go 1.26+.

## Layout

| Path | Package | Responsibility |
| ---- | ------- | -------------- |
| `internal/engine` | `engine` | Core rules engine: game state, combat, and the card-effect AST. Pointerless and clone-friendly. |
| `internal/card` | `card` | Authoring facade over the engine (grouped namespaces like `card.House.X`) plus the card registry. |
| `internal/cards` | `cards` | Card-database aggregator; blank-imports every set so its cards self-register. |
| `internal/cards/sets/<set>` | e.g. `callofthearchons` | One self-registering file per card. |
| `internal/cards/cardtest` | `cardtest` | Shared test harness for the set packages. |
| `internal/tui` | `tui` | [Bubble Tea](https://github.com/charmbracelet/bubbletea) terminal UI. |
| `cmd/tui` | `main` | Thin entry point that calls `tui.Run()`. |

Additional frontends (a web UI, an MCTS bot) are planned as their own `cmd/…`
binaries and `internal/…` packages built on the same engine.

## Development

`make help` lists the available targets: `run`, `build`, `test`, `cover`, `vet`,
`fmt`, `tidy`. Card-authoring conventions live in
[`internal/cards/AGENTS.md`](internal/cards/AGENTS.md).

## License

See [LICENSE](LICENSE).