// Command tui runs the Vactrol terminal UI: a menu to browse the card database
// (card explorer) or play an interactive two-player hotseat match.
package main

import (
	"fmt"
	"os"

	"github.com/dmikalova/vactrol/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "vactrol:", err)
		os.Exit(1)
	}
}
