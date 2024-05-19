package cachery

// Modifiers

func WithLimit(lim int64) Modifier {
	return func(cache any) {
		switch cache.(type) {
		case *lru:
			cache.(*lru).limit = lim
		default:
			panic("not compatible modifier: WithLimit")
		}
	}
}
