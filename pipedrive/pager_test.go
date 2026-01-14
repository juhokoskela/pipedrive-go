package pipedrive

import (
	"context"
	"testing"
)

func TestCursorPager_IteratesPages(t *testing.T) {
	t.Parallel()

	type item struct{ ID int }

	var calls int
	pager := NewCursorPager(func(_ context.Context, cursor *string) ([]item, *string, error) {
		calls++
		switch calls {
		case 1:
			next := "c2"
			if cursor != nil {
				t.Fatalf("expected nil cursor on first call, got %q", *cursor)
			}
			return []item{{ID: 1}}, &next, nil
		case 2:
			if cursor == nil || *cursor != "c2" {
				t.Fatalf("expected cursor c2 on second call, got %v", cursor)
			}
			return []item{{ID: 2}}, nil, nil
		default:
			t.Fatalf("unexpected call count %d", calls)
			return nil, nil, nil
		}
	})

	var got []int
	for pager.Next(context.Background()) {
		for _, it := range pager.Items() {
			got = append(got, it.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("unexpected items: %v", got)
	}
}

