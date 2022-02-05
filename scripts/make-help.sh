#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

echo "Available tasks:"
echo ""
pad=24

help_name="help"
help_pad=$(expr $pad - ${#help_name} - 1)
#echo "help_pad: $help_pad"

printf "  %s% *s%s\n"  "$help_name" $help_pad  "" "show this help"
echo ""

fgrep -h "##" $@ \
    | fgrep -v fgrep \
    | sed -e 's/^/  /' -e 's/:/ /' -e 's/   //g' \
    | sort -k 1 \
    | grep -v '^  #' \
    | awk -v pad=$pad -F "#" '{printf ("%s% *s%s\n", $1, pad-length($1), "", $3)}'
