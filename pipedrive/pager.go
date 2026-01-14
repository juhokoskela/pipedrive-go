package pipedrive

import "context"

type CursorPager[T any] struct {
	fetch func(ctx context.Context, cursor *string) ([]T, *string, error)

	cursor  *string
	started bool
	items   []T
	err     error
}

func NewCursorPager[T any](fetch func(ctx context.Context, cursor *string) ([]T, *string, error)) *CursorPager[T] {
	return &CursorPager[T]{fetch: fetch}
}

func (p *CursorPager[T]) Next(ctx context.Context) bool {
	if p.err != nil {
		return false
	}
	if p.started && p.cursor == nil {
		return false
	}
	p.started = true

	items, next, err := p.fetch(ctx, p.cursor)
	if err != nil {
		p.err = err
		return false
	}
	p.items = items
	p.cursor = next
	return true
}

func (p *CursorPager[T]) Items() []T { return p.items }

func (p *CursorPager[T]) Err() error { return p.err }

