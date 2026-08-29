package engine

// Duration says how long a timed effect lasts. Timed effects (see CannotFight)
// take a Duration, so one effect can serve several windows as more are added —
// the same way a Target says which cards an effect reaches.
type Duration int

const (
	// durationUnset is the invalid zero value: a timed effect must name its
	// duration rather than leave it unset.
	durationUnset Duration = iota
	// NextTurn lasts through the affected player's next turn, then lifts. It waits
	// for that player's next turn however many turns intervene, so an extra turn
	// taken by someone else never consumes it.
	NextTurn
)

// valid reports whether d names a real duration (not the unset zero value).
func (d Duration) valid() bool { return d != durationUnset }
