package cachery

import (
	"context"
	"sync"
)

const DefaultLimit = 6

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
	v, ok := l.vals[key]
	l.mu.RUnlock()
	if ok {
		l.updateCh <- key
	}
	return true, v
}

// Set pushes a value into cache to the first position by key
func (l *lru) Set(ctx context.Context, key string, val any) {
	l.mu.Lock()
	_, ok := l.vals[key]
	if !ok {
		l.vals[key] = val
		l.pushCh <- newNode(key, val)
	}
	l.mu.Unlock()
}

type lruNode struct {
	left  *lruNode
	right *lruNode
	key   string
	val   any
}

func NewLRU() Cache {
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

				c.right = l.head
				c.left = nil
				l.head.left = c
				l.head = c

				break
			}
			c = c.right

		}
	}
	for {
		select {
		case key := <-l.updateCh:
			upd(key)
		case <-l.ctx.Done():
			return
		}
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
	for {
		select {
		case <-l.delCh:
			del()
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *lru) push() {
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
			l.head.left = n
			n.right = l.head
			l.head = n
		}
		l.len++
	}
	for {
		select {
		case n := <-l.pushCh:
			insert(n)
		case <-l.ctx.Done():
			return
		}
	}
}
