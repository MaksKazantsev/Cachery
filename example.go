package cachery

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (bool, any)
	Set(ctx context.Context, key string, val any)
}

func NewCache() Cache {
	return NewLRU()
}
