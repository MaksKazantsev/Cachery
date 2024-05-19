package cachery

import (
	"context"
	"fmt"
)

type Modifier func(cache any)

type Cache interface {
	Get(ctx context.Context, key string) (bool, any)
	Set(ctx context.Context, key string, val any)
	Stop()
}

// NewCache creates a new example of cache
func NewCache(t cacheType, mods ...Modifier) Cache {
	switch t {
	case LRU:
		return NewLRU(mods...)
	case LFU:
		return nil
	default:
		_ = fmt.Errorf("error wrong cache type: %s", t)
		return nil
	}
}

type cacheType string

const (
	LFU          = "LFU"
	LRU          = "LRU"
	DefaultLimit = 10
)
