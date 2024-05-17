package cachery

import (
	"context"
	"testing"
)

func BenchmarkLRUSet(b *testing.B) {
	c := NewCache()

	for i := 0; i < b.N; i++ {
		c.Set(context.Background(), "key", "value")
	}
}

func BenchmarkLRUGet(b *testing.B) {
	c := NewCache()

	for i := 0; i < b.N; i++ {
		c.Get(context.Background(), "key")
	}
}
