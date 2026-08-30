//go:build mage

package main

import "github.com/magefile/mage/sh"

// Lookup prints every source card whose name contains the query, with its set
// code, collector number, house/type/rarity, printed text, and a ready-made
// card.Provenance(...) call — the details a card.New(...) definition needs. Run
// `mage lookup "ether spider"`.
func Lookup(query string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "lookup", query)
}

// Missing lists the source cards in a set that no implemented card tags with a
// provenance Ref yet — the cards still to implement. Pass a set slug, e.g.
// `mage missing callofthearchons` (slugs match the files in
// internal/cards/provenance, minus the .json).
func Missing(setSlug string) error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "missing", setSlug)
}

// Coverage prints, per source set, how many of its cards are covered by an
// implemented card's provenance Ref.
func Coverage() error {
	return sh.RunV("go", "run", "./magefiles/cardlookup", "coverage")
}
