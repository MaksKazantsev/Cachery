package cachery

import (
	"context"
)

type Modifier func(cache any)

type Cache interface {
	Get(ctx context.Context, key string) (bool, any)
	Set(ctx context.Context, key string, val any)
	Stop()
}

// NewCache creates a new example of cache
func NewCache() Cache {
	return &lru{}
}
