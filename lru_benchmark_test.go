package cachery

import (
	"context"
	"testing"
)

func BenchmarkLRUSet(b *testing.B) {
	c := NewCache(LRU)
	ctx := context.Background()

	for i := 0; i <= b.N; i++ {
		c.Set(ctx, "key", "test")
	}
}

func BenchmarkLRUGet(b *testing.B) {
	c := NewCache(LRU)
	ctx := context.Background()

	for i := 0; i <= b.N; i++ {
		c.Get(ctx, "key")
	}
}
