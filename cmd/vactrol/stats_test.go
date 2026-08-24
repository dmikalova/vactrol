package main

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/game"
	"github.com/dmikalova/vactrol/internal/game/cards"
)

// fields finds the table row whose first column equals key and returns it as
// space-separated tokens, so assertions don't depend on the exact padding.
func fields(lines []string, key string) string {
	for _, l := range lines {
		f := strings.Fields(l)
		if len(f) > 0 && f[0] == key {
			return strings.Join(f, " ")
		}
	}
	return ""
}

func TestRenderSetStats(t *testing.T) {
	set := cards.Set{
		Name: "Test Set",
		Cards: []game.CardDefinition{
			game.NewCard("A", game.Brobnar, game.Creature, game.Common),
			game.NewCard("B", game.Brobnar, game.Action, game.Common),
			game.NewCard("C", game.Brobnar, game.Artifact, game.Common),
			game.NewCard("D", game.Logos, game.Creature, game.Common),
		},
	}
	out := renderSetStats(set)
	lines := strings.Split(out, "\n")

	if !strings.Contains(out, "Test Set") || !strings.Contains(out, "(4 cards)") {
		t.Errorf("missing set name/count:\n%s", out)
	}
	// Columns are: House Total Creature Action Artifact Upgrade.
	if got := fields(lines, "Brobnar"); got != "Brobnar 3 1 1 1 0" {
		t.Errorf("Brobnar row = %q", got)
	}
	if got := fields(lines, "Logos"); got != "Logos 1 1 0 0 0" {
		t.Errorf("Logos row = %q", got)
	}
	if got := fields(lines, "Total"); got != "Total 4 2 1 1 0" {
		t.Errorf("Total row = %q", got)
	}
}
