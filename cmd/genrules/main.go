// Command genrules assembles the Vactrol rulebook (docs/rulebook.md) from the
// doc comments on the engine's implementation. It is the source of
// `make generate-rules` (and half of `make gen`).
//
// The idea mirrors gencomments: keep the rules text next to the code that
// enforces them so the two can never drift. A declaration opts into the rulebook
// by carrying a `//rulebook:<section> <Title>` directive in its doc comment; the
// rest of that comment is the entry's body. Entries are grouped by section and
// sorted alphabetically by title, so adding a new keyword or effect just drops it
// into its section without renumbering anything.
//
// Foundational prose that isn't tied to any one declaration — the overview and
// each section's intro — lives as Markdown under docs/rulebook/ and is spliced in
// around the generated entries.
package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// defaultRoot is the package tree scanned for rulebook directives.
	defaultRoot = "internal/engine"
	// proseDir holds the hand-written Markdown fragments (overview + section
	// intros) that frame the generated entries.
	proseDir = "docs/rulebook"
	// outputFile is the assembled rulebook.
	outputFile = "docs/rulebook.md"
	// directive is the doc-comment marker that files a declaration into the
	// rulebook. gofmt treats it as a directive and keeps it verbatim.
	directive = "rulebook:"
)

// section is one top-level rulebook section. Sections print in this fixed order;
// the entries within each are sorted alphabetically by title, so the order of
// terms is data-driven and never has to be maintained by hand.
type section struct {
	key   string // directive key that files entries here, e.g. "keyword"
	title string // heading printed for the section, e.g. "Keywords"
}

// sections is the ordered spine of the rulebook. Add a section by adding a row
// here (and, optionally, a docs/rulebook/<key>.md intro).
var sections = []section{
	{"cardtype", "Card Types"},
	{"keyword", "Keywords"},
	{"ability", "Abilities"},
	{"effect", "Effects"},
}

// entry is one harvested rulebook term: its title and its body Markdown.
type entry struct {
	title string
	body  string
}

func main() {
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{defaultRoot}
	}

	byKey := map[string][]entry{}
	for _, root := range roots {
		if err := harvest(root, byKey); err != nil {
			fmt.Fprintln(os.Stderr, "genrules:", err)
			os.Exit(1)
		}
	}

	doc, err := render(byKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, "genrules:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputFile, []byte(doc), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "genrules:", err)
		os.Exit(1)
	}

	var n int
	for _, es := range byKey {
		n += len(es)
	}
	fmt.Printf("genrules: wrote %s (%d entries)\n", outputFile, n)
}

// harvest walks root, parsing every non-test .go file and collecting the
// rulebook entry from each declaration, spec, or struct field that carries a
// directive.
func harvest(root string, byKey map[string][]entry) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".go" || isTestFile(path) {
			return nil
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		for _, decl := range file.Decls {
			switch dcl := decl.(type) {
			case *ast.GenDecl:
				collect(dcl.Doc, byKey)
				for _, spec := range dcl.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						collect(s.Doc, byKey)
						if st, ok := s.Type.(*ast.StructType); ok && st.Fields != nil {
							for _, f := range st.Fields.List {
								collect(f.Doc, byKey)
							}
						}
					case *ast.ValueSpec:
						collect(s.Doc, byKey)
					}
				}
			case *ast.FuncDecl:
				collect(dcl.Doc, byKey)
			}
		}
		return nil
	})
}

// collect files the entry described by doc, if any, under its section key.
func collect(doc *ast.CommentGroup, byKey map[string][]entry) {
	key, title, body, ok := parseDirective(doc)
	if !ok {
		return
	}
	byKey[key] = append(byKey[key], entry{title: title, body: body})
}

// parseDirective extracts the section key, title, and body from a doc comment.
// The body is every non-directive comment line, so it is unaffected by where
// gofmt places the directive within the comment. ok is false when the comment is
// absent or carries no rulebook directive.
func parseDirective(doc *ast.CommentGroup) (key, title, body string, ok bool) {
	if doc == nil {
		return "", "", "", false
	}
	var bodyLines []string
	for _, c := range doc.List {
		line := commentText(c.Text)
		if k, t, isDir := directiveOf(line); isDir {
			key, title, ok = k, t, true
			continue
		}
		bodyLines = append(bodyLines, strings.TrimPrefix(line, " "))
	}
	if !ok {
		return "", "", "", false
	}
	return key, title, joinBody(bodyLines), true
}

// commentText strips the comment markers from a single comment token, leaving the
// raw line (leading space, if any, preserved so paragraph text keeps its shape).
func commentText(raw string) string {
	if s, ok := strings.CutPrefix(raw, "//"); ok {
		return s
	}
	s := strings.TrimPrefix(raw, "/*")
	s = strings.TrimSuffix(s, "*/")
	return s
}

// directiveOf reports whether line is a rulebook directive and, if so, returns
// its section key and title.
func directiveOf(line string) (key, title string, ok bool) {
	trimmed := strings.TrimSpace(line)
	rest, ok := strings.CutPrefix(trimmed, directive)
	if !ok {
		return "", "", false
	}
	rest = strings.TrimSpace(rest)
	key, title, _ = strings.Cut(rest, " ")
	return key, strings.TrimSpace(title), true
}

// joinBody joins body lines, trims surrounding blank lines, and collapses runs of
// blank lines to a single paragraph break.
func joinBody(lines []string) string {
	var out []string
	prevBlank := true // drop leading blanks
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			if prevBlank {
				continue
			}
			prevBlank = true
			out = append(out, "")
			continue
		}
		prevBlank = false
		out = append(out, l)
	}
	return strings.TrimRight(strings.Join(out, "\n"), "\n")
}

// render assembles the full rulebook Markdown from the harvested entries and the
// prose fragments.
func render(byKey map[string][]entry) (string, error) {
	var b strings.Builder
	b.WriteString("<!-- Code generated by cmd/genrules; DO NOT EDIT. -->\n")
	b.WriteString("<!-- Term entries come from doc comments under " + defaultRoot + "; -->\n")
	b.WriteString("<!-- prose lives in " + proseDir + "/. Run `make gen` to regenerate. -->\n\n")

	overview, err := readProse("overview")
	if err != nil {
		return "", err
	}
	if overview != "" {
		b.WriteString(overview)
		b.WriteString("\n\n")
	}

	for _, sec := range sections {
		entries := byKey[sec.key]
		intro, err := readProse(sec.key)
		if err != nil {
			return "", err
		}
		if len(entries) == 0 && intro == "" {
			continue
		}
		fmt.Fprintf(&b, "## %s\n\n", sec.title)
		if intro != "" {
			b.WriteString(intro)
			b.WriteString("\n\n")
		}
		sort.SliceStable(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].title) < strings.ToLower(entries[j].title)
		})
		for _, e := range entries {
			fmt.Fprintf(&b, "### %s\n\n", e.title)
			if e.body != "" {
				b.WriteString(e.body)
				b.WriteString("\n\n")
			}
		}
	}

	return strings.TrimRight(b.String(), "\n") + "\n", nil
}

// readProse returns the trimmed contents of docs/rulebook/<name>.md, or "" when
// that fragment does not exist.
func readProse(name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(proseDir, name+".md"))
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// isTestFile reports whether path is a _test.go file.
func isTestFile(path string) bool {
	return strings.HasSuffix(filepath.Base(path), "_test.go")
}
