# AGENTS.md

Repo-wide guidance for agents. See `internal/cards/AGENTS.md` for
card-authoring specifics.

## Interpreting requests

Read every request through the lens of **idiomatic, maintainable, composable,
well-written Go that will keep being extended** as more of the game is
implemented. When a request is ambiguous or underspecified, choose the option a
senior Go engineer would find clearest and easiest to build on later — not merely
the shortest path to a passing build.

## Refactoring is welcome

Do **not** contort a new feature to fit the current implementation. If a mechanic
lands more cleanly after reshaping existing code — renaming, splitting, or
generalizing a type, effect, or seam — prefer the refactor. A good fit for the
new feature (and the features that will follow it) matters more than preserving
today's shape. Keep such refactors focused, keep everything green (including 100%
`internal/engine` coverage), and leave the design better than you found it.
