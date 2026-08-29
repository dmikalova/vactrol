package card

import "github.com/dmikalova/vactrol/internal/engine"

// Duration groups the spans a timed effect can last, e.g. card.Duration.NextTurn
// (see card.CannotFight). It mirrors the engine's duration.go.
var Duration = durations{
	NextTurn: engine.NextTurn,
}

type durations struct {
	NextTurn engine.Duration
}
