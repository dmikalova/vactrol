#!/usr/bin/env bash
# Convert KeyForge master-vault pack files into normalized per-set card lists.
#
# For every pack .json in the given directory, extracts the fields we care about
# (including the card number), normalizes the text, and writes
#   internal/cards/provenance/<set>.json
# where <set> is the pack's set name lowercased with all non-alphanumerics
# removed (e.g. "Call of the Archons" -> callofthearchons). Cards are sorted by
# number so the output is stable.
#
# Usage: scripts/cards-to-text.sh [packs-dir]
#   defaults to ~/Code/github.com/dmikalova/keyteki/master-vault-data/packs
set -euo pipefail

if [[ $# -gt 1 ]]; then
	echo "usage: $0 [packs-dir]" >&2
	exit 1
fi

packs_dir="${1:-$HOME/Code/github.com/dmikalova/keyteki/master-vault-data/packs}"
script_dir="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

shopt -s nullglob
packs=("$packs_dir"/*.json)
if [[ ${#packs[@]} -eq 0 ]]; then
	echo "no .json pack files found in $packs_dir" >&2
	exit 1
fi

for pack in "${packs[@]}"; do
	if ! jq -e 'has("cards") and (.name | type == "string")' "$pack" >/dev/null 2>&1; then
		echo "skipping $pack: not a pack file" >&2
		continue
	fi
	set_slug="$(jq -r '.name | gsub("\u00c6"; "Ae") | gsub("\u00e6"; "ae") | ascii_downcase | gsub("[^a-z0-9]"; "")' "$pack")"
	if [[ -z "$set_slug" || "$set_slug" == "null" ]]; then
		echo "skipping $pack: no set name" >&2
		continue
	fi
	out_dir="$repo_root/internal/cards/provenance"
	mkdir -p "$out_dir"
	jq '
		# Normalize every string value (card names and text): curly quotes ->
		# straight, dashes -> hyphen, Ae/amber/damage symbols, CR/vertical-tab ->
		# newline, collapse spaces, and trim.
		def norm:
			gsub("\u2019"; "\u0027")   # ’ right single quote -> apostrophe
			| gsub("\u2018"; "\u0027") # ‘ left single quote  -> apostrophe
			| gsub("\u201c"; "\"")     # “ left double quote  -> "
			| gsub("\u201d"; "\"")     # ” right double quote -> "
			| gsub("\u2013"; "-")      # – en dash -> hyphen
			| gsub("\u2014"; "-")      # — em dash -> hyphen
			| gsub("\u00c6"; "Ae")     # Æ -> Ae
			# Fold remaining non-ASCII letters and stray spaces/punctuation to ASCII.
			| gsub("\u00e6"; "ae")     # æ -> ae
			| gsub("\u0103"; "a") | gsub("\u0102"; "A")   # ă Ă
			| gsub("\u0115"; "e") | gsub("\u0114"; "E")   # ĕ Ĕ
			| gsub("\u012d"; "i") | gsub("\u012c"; "I")   # ĭ Ĭ
			| gsub("\u014f"; "o") | gsub("\u014e"; "O")   # ŏ Ŏ
			| gsub("\u016d"; "u") | gsub("\u016c"; "U")   # ŭ Ŭ
			| gsub("\u00e4"; "a") | gsub("\u00c4"; "A")   # ä Ä
			| gsub("\u00f6"; "o") | gsub("\u00d6"; "O")   # ö Ö
			| gsub("\u00fc"; "u") | gsub("\u00dc"; "U")   # ü Ü
			| gsub("\u00e9"; "e") | gsub("\u00c9"; "E")   # é É
			| gsub("\u00e8"; "e") | gsub("\u00c8"; "E")   # è È
			| gsub("\u00e2"; "a") | gsub("\u00c2"; "A")   # â Â
			| gsub("\ufeff"; "")       # BOM / zero-width no-break space -> drop
			| gsub("\u202f"; " ")      # narrow no-break space -> space
			| gsub("\u2011"; "-")      # non-breaking hyphen -> hyphen
			| gsub("\u2022"; "-")      # bullet -> hyphen
			| gsub("<A>"; " Aember")   # amber symbol (legacy tag)
			| gsub("<D>"; " Damage")   # damage symbol (legacy tag)
			| gsub("\uf360"; " Aember")   # Master Vault amber icon
			| gsub("\uf361"; " Damage")   # Master Vault damage icon
			# Master Vault bonus/enhance icons. The capture pip is a two-codepoint
			# glyph (F36F+F560); Dark Tidings uses its own capture (F565) + tide
			# (F566) codepoints. House-enhance pips map to the house name.
			| gsub("\uf36f\uf560"; " Capture")   # capture pip
			| gsub("\uf565"; " Capture")         # capture pip (Dark Tidings)
			| gsub("\uf36e"; " Draw")            # draw pip
			| gsub("\uf372"; " Discard")         # discard pip
			| gsub("\uf392"; " +1 power counter") # +1 power counter pip
			| gsub("\uf566"; "")                 # tide icon (Dark Tidings) -> drop
			| gsub("\uf379"; " Brobnar")
			| gsub("\uf37a"; " Dis")
			| gsub("\uf37b"; " Ekwidon")
			| gsub("\uf37c"; " Geistoid")
			| gsub("\uf37d"; " Logos")
			| gsub("\uf37e"; " Mars")
			| gsub("\uf37f"; " Skyborn")
			| gsub("\uf386"; " Redemption")
			| gsub("\uf387"; " Sanctum")
			| gsub("\uf388"; " Saurian")
			| gsub("\uf389"; " Shadows")
			| gsub("\uf38a"; " Star Alliance")
			| gsub("\uf38b"; " Untamed")
			| gsub("\uf390"; " Unfathomable")
			| gsub("\uf391"; " Ouboros")
			| gsub("\u00a0"; " ")      # no-break space -> space
			| gsub("\u2026"; "...")    # horizontal ellipsis -> ...
			| gsub("\r\n?"; "\n")      # carriage return (CRLF or CR) -> newline
			| gsub("\u000b"; "\n")     # vertical tab -> newline
			| gsub(" {2,}"; " ")       # collapse runs of spaces
			| gsub(" *\n"; "\n")       # drop spaces before newlines
			| sub("^\\s+"; "")         # trim leading whitespace
			| sub("\\s+$"; "");        # trim trailing whitespace
		.cards
		| sort_by(.number)
		| map({number, name, house, keywords, traits, type, rarity, amber, armor, power, text}
			| with_entries(select(.value != 0 and .value != null and .value != [])))
		| walk(if type == "string" then norm else . end)
	' "$pack" >"$out_dir/$set_slug.json"
	echo "wrote $out_dir/$set_slug.json" >&2
done
