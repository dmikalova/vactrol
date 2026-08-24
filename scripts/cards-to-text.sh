#!/usr/bin/env bash
# Convert a card set JSON (e.g. docs/cota.json) into a list of card objects
# with just the fields we care about.
#
# Usage: scripts/cards-to-text.sh <file.json>
set -euo pipefail

if [[ $# -ne 1 ]]; then
	echo "usage: $0 <file.json>" >&2
	exit 1
fi

jq '.cards | map({name, house, keywords, traits, type, rarity, amber, armor, power, text}
	| with_entries(select(.value != 0 and .value != null and .value != []))
)' "$1" | "$(dirname "$0")/normalize-cards.sh"
