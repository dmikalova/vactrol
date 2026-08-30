//go:build mage

package main

import (
	"os"
	"path/filepath"

	"github.com/dmikalova/vactrol/internal/hotreload"
	"github.com/magefile/mage/sh"
)

// Tui launches the vactrol TUI (card explorer / play a hotseat game).
func Tui() error {
	return sh.RunV("go", "run", "./cmd/tui")
}

// WebWasm builds the web client to WebAssembly (web/app.wasm).
func WebWasm() error {
	return sh.RunWithV(
		map[string]string{"GOARCH": "wasm", "GOOS": "js"},
		"go", "build", "-o", "web/app.wasm", "./cmd/web",
	)
}

// Web builds the wasm client and serves it at http://localhost:8000, rebuilding
// and restarting on any Go or CSS change so edits show up live. Each restart
// bumps go-app's version; the browser polls for it (see cmd/web devReload),
// reloads, and OnMount resumes the in-progress match. No external watcher needed;
// press Ctrl-C to stop.
func Web() error {
	bin := filepath.Join(os.TempDir(), "vactrol-web-dev")
	return hotreload.Serve(hotreload.Config{
		Build: func() error {
			if err := WebWasm(); err != nil {
				return err
			}
			return sh.Run("go", "build", "-o", bin, "./cmd/web")
		},
		Command:    bin,
		Extensions: []string{".go", ".css"},
	})
}
