package card

import "github.com/dmikalova/vactrol/internal/cards/provenance"

// Provenance tags a card as derived from an original source card (set + collector
// number), for coverage tracking. Optional and repeatable — a card may draw from
// more than one original — e.g. card.Provenance(card.CotA, 1).
func Provenance(set provenance.SourceSet, number int) Option {
	return func(b *builder) { b.prov = append(b.prov, provenance.Ref{Set: set, Number: number}) }
}

// Source sets to tag a card's Provenance with, e.g. card.Provenance(card.CotA, 1).
var (
	CotA = provenance.CallOfTheArchons
	AoA  = provenance.AgeOfAscension
	WC   = provenance.WorldsCollide
	MM   = provenance.MassMutation
	DT   = provenance.DarkTidings
	WoE  = provenance.WindsOfExchange
	GR   = provenance.GrimReminders
	AS   = provenance.AemberSkies
	ToC  = provenance.TokensOfChange
	MoM  = provenance.MoreMutation
	Men  = provenance.Menagerie
	VM   = provenance.VaultMasters2025
	PV   = provenance.PropheticVisions
	CC   = provenance.CrucibleClash
	DM   = provenance.DraconianMeasures
)
