package Cachery

import (
	"context"
	"sync"
)

const DefaultCapacity = 10

type lru struct {
	cacheCapacity int

	mu sync.RWMutex

	vals map[string]*lruNode

	tail *lruNode
	head *lruNode
}

type lruNode struct {
	left  *lruNode
	right *lruNode

	key string
	val any
}

func NewLRU(m ...ModifierFunc) Cache {
	c := &lru{
		cacheCapacity: DefaultCapacity,
		vals:          make(map[string]*lruNode),
	}

	for _, mod := range m {
		mod(c)
	}
	return c
}

func (l *lru) Get(ctx context.Context, key string) (any, bool) {
	l.mu.RLock()
	val, ok := l.vals[key]
	l.mu.RUnlock()

	if ok {
		l.pushToStart(val)
		return val.val, true
	}
	return nil, false
}

func (l *lru) Set(ctx context.Context, key string, val any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node := newNode(key, val)

	existingNode, exists := l.vals[key]
	if exists {
		existingNode.val = val
		l.pushToStart(node)
		return
	}

	l.vals[key] = node

	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.right = l.head
		l.head.left = node
		l.head = node
	}

	if len(l.vals) > l.cacheCapacity {
		if l.tail != nil {
			delete(l.vals, l.tail.key)
			l.tail = l.tail.left
			if l.tail != nil {
				l.tail.right = nil
			} else {
				l.head = nil
			}

		}
	}

}

func (l *lru) Stop() {
	clear(l.vals)
	l.head = nil
	l.tail = nil
}
