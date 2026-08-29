// Command genprovenance converts KeyForge master-vault pack files into the
// normalized per-set card catalogs under internal/cards/provenance. It is the Go
// port of the old scripts/generate-provenance.sh and the source of
// `mage generateProvenance` (deliberately NOT part of `mage gen`, since it needs
// the external master-vault data checkout to run).
//
// For every pack .json in the packs directory it extracts the fields vactrol
// cares about (number, name, house, keywords, traits, type, rarity, amber, armor,
// power, text), drops zero/empty values, normalizes every string (curly quotes,
// dashes, Æ/icon glyphs, whitespace) to plain ASCII, sorts the cards by number,
// and writes internal/cards/provenance/<slug>.json — where <slug> is the pack's
// set name lowercased with all non-alphanumerics removed (e.g. "Call of the
// Archons" -> callofthearchons).
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// outDir is where the normalized per-set catalogs are written, relative to the
// repository root (the working directory when run via `go run`/mage).
const outDir = "internal/cards/provenance"

// defaultPacksSubdir is the master-vault packs checkout, relative to $HOME.
const defaultPacksSubdir = "Code/github.com/dmikalova/keyteki/master-vault-data/packs"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: genprovenance [packs-dir]")
	}

	packsDir := ""
	if len(args) == 1 {
		packsDir = args[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		packsDir = filepath.Join(home, defaultPacksSubdir)
	}

	packs, err := filepath.Glob(filepath.Join(packsDir, "*.json"))
	if err != nil {
		return err
	}
	if len(packs) == 0 {
		return fmt.Errorf("no .json pack files found in %s", packsDir)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	for _, pack := range packs {
		if err := convert(pack); err != nil {
			return fmt.Errorf("%s: %w", pack, err)
		}
	}
	return nil
}

// card is both the subset of a master-vault card we read and the normalized shape
// we write. Text is a pointer so an empty string (a printed vanilla creature) is
// kept while an absent/null text is dropped, matching the original jq filter.
type card struct {
	Number   string   `json:"number,omitempty"`
	Name     string   `json:"name,omitempty"`
	House    string   `json:"house,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
	Traits   []string `json:"traits,omitempty"`
	Type     string   `json:"type,omitempty"`
	Rarity   string   `json:"rarity,omitempty"`
	Amber    int      `json:"amber,omitempty"`
	Armor    int      `json:"armor,omitempty"`
	Power    int      `json:"power,omitempty"`
	Text     *string  `json:"text,omitempty"`
}

// convert reads one master-vault pack and writes its normalized catalog. Non-pack
// JSON (no "cards" array, or a non-string name) is skipped with a warning, as the
// original script did.
func convert(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Fprintf(os.Stderr, "skipping %s: not a JSON object\n", path)
		return nil
	}
	nameRaw, hasName := raw["name"]
	cardsRaw, hasCards := raw["cards"]
	if !hasName || !hasCards {
		fmt.Fprintf(os.Stderr, "skipping %s: not a pack file\n", path)
		return nil
	}
	var name string
	if err := json.Unmarshal(nameRaw, &name); err != nil {
		fmt.Fprintf(os.Stderr, "skipping %s: no set name\n", path)
		return nil
	}
	slug := setSlug(name)
	if slug == "" {
		fmt.Fprintf(os.Stderr, "skipping %s: no set name\n", path)
		return nil
	}

	var cards []card
	if err := json.Unmarshal(cardsRaw, &cards); err != nil {
		return err
	}
	for i := range cards {
		cards[i].normalize()
	}
	sort.SliceStable(cards, func(i, j int) bool { return cards[i].Number < cards[j].Number })

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cards); err != nil {
		return err
	}

	out := filepath.Join(outDir, slug+".json")
	if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", out)
	return nil
}

// normalize folds every string field (and slice element) to plain ASCII.
func (c *card) normalize() {
	c.Number = norm(c.Number)
	c.Name = norm(c.Name)
	c.House = norm(c.House)
	c.Type = norm(c.Type)
	c.Rarity = norm(c.Rarity)
	for i := range c.Keywords {
		c.Keywords[i] = norm(c.Keywords[i])
	}
	for i := range c.Traits {
		c.Traits[i] = norm(c.Traits[i])
	}
	if c.Text != nil {
		t := norm(*c.Text)
		c.Text = &t
	}
}

// setSlug derives the catalog filename stem from a set name: Æ/æ folded to Ae/ae,
// lowercased, then all non-alphanumeric characters removed.
func setSlug(name string) string {
	name = strings.ReplaceAll(name, "\u00c6", "Ae")
	name = strings.ReplaceAll(name, "\u00e6", "ae")
	name = strings.ToLower(name)
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// literals maps every fixed source glyph to its ASCII replacement: curly quotes,
// dashes, the Æ/damage/amber tags, the Master Vault icon-font Private Use Area
// glyphs (amber/damage/capture/draw/discard/counter pips and the per-house enhance
// pips), and stray non-breaking spaces. The two-codepoint capture pip is listed
// before the single-codepoint icons so it is matched whole.
var literals = strings.NewReplacer(
	"\u2019", "'", // ’ right single quote -> apostrophe
	"\u2018", "'", // ‘ left single quote  -> apostrophe
	"\u201c", "\"", // “ left double quote  -> "
	"\u201d", "\"", // ” right double quote -> "
	"\u2013", "-", // – en dash -> hyphen
	"\u2014", "-", // — em dash -> hyphen
	"\u00c6", "Ae", // Æ -> Ae
	"\u00e6", "ae", // æ -> ae
	"\u0103", "a", "\u0102", "A", // ă Ă
	"\u0115", "e", "\u0114", "E", // ĕ Ĕ
	"\u012d", "i", "\u012c", "I", // ĭ Ĭ
	"\u014f", "o", "\u014e", "O", // ŏ Ŏ
	"\u016d", "u", "\u016c", "U", // ŭ Ŭ
	"\u00e4", "a", "\u00c4", "A", // ä Ä
	"\u00f6", "o", "\u00d6", "O", // ö Ö
	"\u00fc", "u", "\u00dc", "U", // ü Ü
	"\u00e9", "e", "\u00c9", "E", // é É
	"\u00e8", "e", "\u00c8", "E", // è È
	"\u00e2", "a", "\u00c2", "A", // â Â
	"\ufeff", "", // BOM / zero-width no-break space -> drop
	"\u202f", " ", // narrow no-break space -> space
	"\u2011", "-", // non-breaking hyphen -> hyphen
	"\u2022", "-", // bullet -> hyphen
	"<A>", " Aember", // amber symbol (legacy tag)
	"<D>", " Damage", // damage symbol (legacy tag)
	"\uf360", " Aember", // Master Vault amber icon
	"\uf361", " Damage", // Master Vault damage icon
	"\uf36f\uf560", " Capture", // capture pip (two-codepoint glyph)
	"\uf565", " Capture", // capture pip (Dark Tidings)
	"\uf36e", " Draw", // draw pip
	"\uf372", " Discard", // discard pip
	"\uf392", " +1 power counter", // +1 power counter pip
	"\uf566", "", // tide icon (Dark Tidings) -> drop
	"\uf379", " Brobnar",
	"\uf37a", " Dis",
	"\uf37b", " Ekwidon",
	"\uf37c", " Geistoid",
	"\uf37d", " Logos",
	"\uf37e", " Mars",
	"\uf37f", " Skyborn",
	"\uf386", " Redemption",
	"\uf387", " Sanctum",
	"\uf388", " Saurian",
	"\uf389", " Shadows",
	"\uf38a", " Star Alliance",
	"\uf38b", " Untamed",
	"\uf390", " Unfathomable",
	"\uf391", " Ouboros",
	"\u00a0", " ", // no-break space -> space
	"\u2026", "...", // horizontal ellipsis -> ...
	"\u000b", "\n", // vertical tab -> newline
)

var (
	reCarriageReturn = regexp.MustCompile(`\r\n?`) // CRLF or bare CR -> newline
	reSpaces         = regexp.MustCompile(` {2,}`) // runs of spaces -> one
	reSpaceNewline   = regexp.MustCompile(` *\n`)  // spaces before a newline -> drop
)

// norm normalizes a single string value the way the old jq `norm` filter did:
// fold fixed glyphs to ASCII, turn carriage returns and vertical tabs into
// newlines, collapse space runs, drop spaces before newlines, and trim.
func norm(s string) string {
	s = literals.Replace(s)
	s = reCarriageReturn.ReplaceAllString(s, "\n")
	s = reSpaces.ReplaceAllString(s, " ")
	s = reSpaceNewline.ReplaceAllString(s, "\n")
	return strings.TrimSpace(s)
}
