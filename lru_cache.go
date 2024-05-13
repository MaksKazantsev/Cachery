package cachery

import (
	"context"
	"sync"
)

const Default = 10

type lru struct {
	tail *lruNode
	head *lruNode

	vals map[string]any

	stopCtx       context.Context
	stopCancelCtx context.CancelFunc

	mu sync.RWMutex

	delCh       chan struct{}
	pushCh      chan *lruNode
	updatePosCh chan string

	limit uint64
	len   uint64
}

type lruNode struct {
	left  *lruNode
	right *lruNode

	key string
	val any
}

func NewLRU() {
	ctx, cancel := context.WithCancel(context.Background())
	c := lru{
		limit:         Default,
		vals:          make(map[string]any),
		stopCtx:       ctx,
		stopCancelCtx: cancel,
		delCh:         make(chan struct{}),
		pushCh:        make(chan *lruNode),
		updatePosCh:   make(chan string),
	}

	go c.del()
	go c.push()
	go c.updatePos()
}

func newNode(key string, val any) *lruNode {
	return &lruNode{key: key, val: val}
}

func (c *lru) Set(ctx context.Context, key string, val any) {
	c.mu.Lock()
	_, ok := c.vals[key]
	if !ok {
		c.vals[key] = val
		c.pushCh <- newNode(key, val)
	}
	c.mu.Unlock()
}

func (c *lru) Get(ctx context.Context, key string) any {
	c.mu.RLock()
	val, ok := c.vals[key]
	if ok {
		c.updatePosCh <- key
	}
	c.mu.RUnlock()
	return val
}
func (c *lru) del() {
	deleteAtEnd := func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		delete(c.vals, c.tail.key)
		c.tail = c.tail.left
		c.tail.right = nil
	}

	for {
		select {
		case <-c.delCh:
			go deleteAtEnd()
		case <-c.stopCtx.Done():
			return
		}
	}
}
func (c *lru) push() {
	insertToStart := func(n *lruNode) {
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.head == nil {
			c.head = n
			c.tail = n
		} else {
			if c.len >= c.limit {
				c.delCh <- struct{}{}
			}
			n.left = c.head
			c.head.right = n
			c.head = n
		}
		c.len++
	}

	for {
		select {
		case n := <-c.pushCh:
			go insertToStart(n)
		case <-c.stopCtx.Done():
			return
		}
	}
}
func (c *lru) updatePos() {
	removeAndAdd := func(key string) {
		c.mu.Lock()
		defer c.mu.Unlock()

		curr := c.head
		for curr != nil {
			if curr.key == key && curr.right != nil {
				if curr.left != nil {
					curr.right.left = curr.left
				} else {
					c.tail = curr.right
				}
				curr.right.left = curr.left

				curr.right = nil
				curr.left = c.head
				c.head.right = curr
				c.head = curr

				break
			}
			curr = curr.left
		}
	}
	for {
		select {
		case key := <-c.updatePosCh:
			go removeAndAdd(key)
		case <-c.stopCtx.Done():
			return
		}
	}
}
