package cachery

import (
	"context"
	"sync"
)

type lru struct {
	head *lruNode
	tail *lruNode

	mu   sync.RWMutex
	vals map[string]any

	ctx       context.Context
	cancelCtx context.CancelFunc

	pushCh   chan *lruNode
	delCh    chan struct{}
	updateCh chan string

	len   int64
	limit int64
}

// Get takes the value from cache by key and updates it's position in cache
func (l *lru) Get(ctx context.Context, key string) (bool, any) {
	l.mu.RLock()
	val, ok := l.vals[key]
	l.mu.RUnlock()

	if ok {
		select {
		case l.updateCh <- key:
		case <-l.ctx.Done():
			close(l.updateCh)
		}
	}
	return ok, val
}

// Set pushes a value into cache to the first position by key
func (l *lru) Set(ctx context.Context, key string, val any) {
	l.mu.Lock()
	_, ok := l.vals[key]
	if ok {
		l.mu.Unlock()
		return
	}
	l.vals[key] = val
	l.mu.Unlock()

	select {
	case l.pushCh <- newNode(key, val):
	case <-l.ctx.Done():
		close(l.pushCh)
	}
}

// Stop stops the cache
func (l *lru) Stop() {
	l.cancelCtx()
}

type lruNode struct {
	left  *lruNode
	right *lruNode
	key   string
	val   any
}

func NewLRU(mods ...Modifier) Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &lru{
		vals:      make(map[string]any),
		ctx:       ctx,
		cancelCtx: cancel,
		pushCh:    make(chan *lruNode),
		delCh:     make(chan struct{}),
		updateCh:  make(chan string),
		limit:     DefaultLimit,
	}
	for _, mod := range mods {
		mod(&c)
	}

	go c.push()
	go c.pop()
	go c.update()
	return c
}

func newNode(key string, val any) *lruNode {
	return &lruNode{key: key, val: val}
}

func (l *lru) update() {
	upd := func(key string) {
		l.mu.Lock()
		defer l.mu.Unlock()

		c := l.head

		for c != nil {
			if c.key == key && c.left != nil {
				if c.right != nil {
					c.right.left = c.left
				} else {
					l.tail = c.left
				}
				c.left.right = c.right

				c.left = nil
				c.right = l.head
				l.head.left = c
				l.head = c

				break
			}
			c = c.right

		}
	}
	for key := range l.updateCh {
		upd(key)
	}
}

func (l *lru) pop() {
	del := func() {
		l.mu.Lock()
		defer l.mu.Unlock()

		delete(l.vals, l.tail.key)
		l.tail = l.tail.left
		l.tail.right = nil
	}
	for range l.delCh {
		del()
	}
}

func (l *lru) push() {
	defer func() {
		close(l.delCh)
	}()
	insert := func(n *lruNode) {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.head == nil {
			l.head = n
			l.tail = n
		} else {
			if l.len >= l.limit {
				l.delCh <- struct{}{}
			}

			n.right = l.head
			l.head.left = n
			l.head = n
		}
		l.len++
	}
	for n := range l.pushCh {
		insert(n)
	}
}
