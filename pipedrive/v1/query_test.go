package v1

import (
	"net/url"
	"testing"
)

func TestMergeQueryValues_LastWriterWins(t *testing.T) {
	t.Parallel()

	dst := mergeQueryValues(nil, url.Values{"limit": {"10"}})
	dst = mergeQueryValues(dst, url.Values{"limit": {"20"}})

	if got := dst["limit"]; len(got) != 1 || got[0] != "20" {
		t.Fatalf("expected the later value to replace the earlier one, got %v", got)
	}
}

func TestMergeQueryValues_PreservesMultiValueKeys(t *testing.T) {
	t.Parallel()

	dst := mergeQueryValues(nil, url.Values{"ids": {"1", "2"}})
	if got := dst["ids"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
		t.Fatalf("expected both values preserved, got %v", got)
	}
}

func TestMergeQueryValues_DoesNotAliasSource(t *testing.T) {
	t.Parallel()

	src := url.Values{"ids": {"1"}}
	dst := mergeQueryValues(nil, src)
	src["ids"][0] = "mutated"

	if got := dst["ids"][0]; got != "1" {
		t.Fatalf("merged values must not alias the source slice, got %q", got)
	}
}

func TestMergeQueryValues_EmptySourceKeepsDestination(t *testing.T) {
	t.Parallel()

	dst := url.Values{"limit": {"10"}}
	if got := mergeQueryValues(dst, nil); len(got) != 1 {
		t.Fatalf("expected destination untouched, got %v", got)
	}
}
