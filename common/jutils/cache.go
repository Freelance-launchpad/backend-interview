package jutils

import "context"

type Cache[K comparable, V any] struct {
	m map[K]V

	store func(context.Context, K) (V, error)
}

func NewCache[K comparable, V any](store func(context.Context, K) (V, error)) *Cache[K, V] {
	return &Cache[K, V]{
		store: store,
		m:     make(map[K]V),
	}
}

func (c *Cache[K, V]) Get(ctx context.Context, key K) (V, error) {
	if value, ok := c.m[key]; ok {
		return value, nil
	}

	value, err := c.store(ctx, key)
	if err != nil {
		var zero V
		return zero, err
	}

	c.m[key] = value
	return value, nil
}
