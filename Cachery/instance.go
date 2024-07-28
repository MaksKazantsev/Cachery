package Cachery

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (any, bool)
	Set(ctx context.Context, key string, val any)
	Stop()
}

type cacheType string

const (
	LRU cacheType = "LRU"
)

func NewCache(cacheT cacheType, m ...ModifierFunc) Cache {
	switch cacheT {
	case LRU:
		return NewLRU(m...)
	default:
		panic("unknown type of cache!")
	}
}
