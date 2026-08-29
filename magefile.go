//go:build mage

// Vactrol developer tasks. Run `mage` to list targets, `mage <target>` to run
// one, or `mage -l` for the same list with descriptions.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmikalova/vactrol/internal/hotreload"
	"github.com/magefile/mage/mg"
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

// Build builds all packages.
func Build() error {
	return sh.RunV("go", "build", "./...")
}

// Test runs all tests.
func Test() error {
	return sh.RunV("go", "test", "./...")
}

// Cover runs tests and reports engine coverage (kept at 100%).
func Cover() error {
	if err := sh.RunV("go", "test", "./internal/engine/", "-coverprofile=coverage.out"); err != nil {
		return err
	}
	out, err := sh.Output("go", "tool", "cover", "-func=coverage.out")
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	fmt.Println(lines[len(lines)-1])
	return nil
}

// Vet runs go vet.
func Vet() error {
	return sh.RunV("go", "vet", "./...")
}

// Fmt formats all Go files.
func Fmt() error {
	return sh.RunV("gofmt", "-w", ".")
}

// FmtCheck fails if any Go file is not gofmt-clean (does not modify files).
func FmtCheck() error {
	out, err := sh.Output("gofmt", "-l", ".")
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("gofmt needed:\n%s", out)
	}
	return nil
}

// Check is the full green gate: fmt-check, build, vet, test, coverage.
func Check() error {
	mg.Deps(FmtCheck, Build, Vet, Test, Cover)
	fmt.Println("ALL GREEN")
	return nil
}

// Tidy tidies module dependencies.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}

// GenerateComments rewrites each card's doc comment from its definition.
func GenerateComments() error {
	return sh.RunV("go", "run", "./cmd/gencomments")
}

// GenerateRules generates docs/rulebook.md from engine doc comments.
func GenerateRules() error {
	return sh.RunV("go", "run", "./cmd/genrules")
}

// GenerateProvenance rebuilds the provenance card catalogs from the master-vault
// pack data. It is intentionally NOT part of Gen: it needs the external
// master-vault-data checkout, so it is run by hand when that source data changes.
func GenerateProvenance() error {
	return sh.RunV("go", "run", "./cmd/genprovenance")
}

// Gen regenerates card comments and the rulebook.
func Gen() {
	mg.Deps(GenerateComments, GenerateRules)
}
