package cachery

// WithLimit adds limit to a cache
func WithLimit(lim int64) Modifier {
	return func(cache any) {
		switch cache.(type) {
		case *lru:
			cache.(*lru).limit = lim
		default:
			panic("incompatible")
		}
	}
}
