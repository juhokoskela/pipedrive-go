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

awk -v ver="$version" '
	$0 ~ "^## \\[" ver "\\]" {found=1; print; next}
	found && $0 ~ "^## \\[" {exit}
	found {print}
	END {
		if (!found) {
			exit 1
		}
	}
' "$file"
