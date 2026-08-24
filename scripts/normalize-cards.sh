#!/usr/bin/env bash
# Normalize the text of card objects read as JSON on stdin, writing normalized
# JSON to stdout. Runs through jq so replacing curly quotes with straight ones
# stays valid JSON (they get re-escaped). Applied to every string value (card
# names and text):
#   ’ ‘   -> '        curly single quotes -> apostrophe
#   “ ”   -> "        curly double quotes -> straight
#   – —   -> -        en / em dash -> hyphen
#   Æ     -> Ae       so "Irradiated Æmber" matches the <A> spelling
#   <A>   -> Aember   KeyForge amber symbol
#   <D>   -> Damage   KeyForge damage symbol
#   U+000B (\u000b, vertical tab) -> newline   KeyForge line break
#   collapses runs of spaces and trims stray leading/trailing whitespace
#
# Usage: <json on stdin> | scripts/normalize-cards.sh
set -euo pipefail

jq '
  def norm:
    gsub("\u2019"; "\u0027")   # ’ right single quote -> apostrophe
    | gsub("\u2018"; "\u0027") # ‘ left single quote  -> apostrophe
    | gsub("\u201c"; "\"")     # “ left double quote  -> "
    | gsub("\u201d"; "\"")     # ” right double quote -> "
    | gsub("\u2013"; "-")      # – en dash -> hyphen
    | gsub("\u2014"; "-")      # — em dash -> hyphen
    | gsub("\u00c6"; "Ae")     # Æ -> Ae
    | gsub("<A>"; " Aember")   # amber symbol
    | gsub("<D>"; " Damage")   # damage symbol
    | gsub("\u000b"; "\n")     # vertical tab -> newline
    | gsub(" {2,}"; " ")       # collapse runs of spaces
    | gsub(" *\n"; "\n")       # drop spaces before newlines
    | sub("^\\s+"; "")         # trim leading whitespace
    | sub("\\s+$"; "");        # trim trailing whitespace
  walk(if type == "string" then norm else . end)
'
