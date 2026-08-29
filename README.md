# Vactrol

Vactrol is a rules engine for a [KeyForge](https://keyforge.fandom.com/)-style
card game, written in Go. It models a full two-player match — houses, Æmber,
keys, combat, and card abilities — behind a small, pointerless state type that is
cheap to copy (for AI search), and ships with a terminal UI for browsing the card
database and playing hotseat games.

## Quick start

```sh
mage run      # launch the terminal UI (card explorer + hotseat game)
mage web      # build the wasm client and serve it at http://localhost:8000
mage test     # run the test suite
mage cover    # engine test coverage (kept at 100%)
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
| `internal/match` | `match` | Shared match setup (random decks, house list) used by every frontend. |
| `internal/web` | `web` | [go-app](https://github.com/maxence-charriere/go-app) WebAssembly client, Monokai-themed. |
| `cmd/web` | `main` | Serves the web client and, compiled to wasm, runs it in the browser. |

The web UI ships today; further frontends (e.g. an MCTS bot) and a lobby server
are planned as their own `cmd/…` binaries and `internal/…` packages on the same
engine.

## Development

`mage -l` lists the available targets: `run`, `web`, `webWasm`, `build`, `test`,
`cover`, `vet`, `fmt`, `tidy`, `gen`. Card-authoring conventions live in
[`internal/cards/AGENTS.md`](internal/cards/AGENTS.md).

## License

See [LICENSE](LICENSE).