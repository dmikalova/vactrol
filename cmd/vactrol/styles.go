package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	headerStyle   = lipgloss.NewStyle().Bold(true).Underline(true)
	faintStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))

	// activeHouseStyle colors cards of the turn's active house (green);
	// otherHouseStyle colors cards of every other house (blue).
	activeHouseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	otherHouseStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
)

// cursor returns the left-hand marker for a list row.
func cursor(selected bool) string {
	if selected {
		return "› "
	}
	return "  "
}

// clamp constrains v to [lo, hi]; if the range is empty (hi < lo) it returns lo.
func clamp(v, lo, hi int) int {
	switch {
	case hi < lo:
		return lo
	case v < lo:
		return lo
	case v > hi:
		return hi
	default:
		return v
	}
}

// wrap moves an index by delta within [0, n), wrapping around the ends. It
// returns 0 when the list is empty.
func wrap(cur, delta, n int) int {
	if n <= 0 {
		return 0
	}
	return ((cur+delta)%n + n) % n
}
