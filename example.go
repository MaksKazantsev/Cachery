package cachery

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (any, bool)
	Set(ctx context.Context, key string, val any)
}

// New - cache constructor, requires cache type, additionally modifiers
func New() Cache {
	return NewLRU()
}
