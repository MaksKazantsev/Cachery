package cache

type ModifierFunc func(cache Cache)

func WithCapacity(cap int) ModifierFunc {
	return func(cache Cache) {
		switch cache.(type) {
		case *lru:
			cache.(*lru).cacheCapacity = cap
		default:
			panic("incompatible modifier: WithLimit")
		}
	}
}
