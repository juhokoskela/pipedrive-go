#!/bin/sh
set -eu

if [ "$#" -ne 1 ]; then
	echo "usage: $0 <tag>" >&2
	exit 2
fi

tag="$1"
version="${tag#v}"
file="CHANGELOG.md"

if [ ! -f "$file" ]; then
	echo "CHANGELOG.md not found" >&2
	exit 1
fi

# Match headings literally: version strings contain dots, which would act
# as regex wildcards and could select a lookalike heading.
awk -v ver="$version" '
	BEGIN {heading = "## [" ver "]"}
	index($0, heading) == 1 {found=1; print; next}
	found && index($0, "## [") == 1 {exit}
	found {print}
	END {
		if (!found) {
			exit 1
		}
	}
' "$file"
