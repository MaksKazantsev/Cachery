package cachery

import (
	"context"
	"sync"
)

const DefaultLimit = 10

type lru struct {
	len   uint64
	limit uint64

	ctx       context.Context
	cancelCtx context.CancelFunc

	mu   sync.RWMutex
	vals map[string]any

	head *lruNode
	tail *lruNode

	pushCh    chan *lruNode
	popCh     chan struct{}
	changePos chan string
}

func (l *lru) Get(ctx context.Context, key string) any {
	l.mu.RLock()
	v, ok := l.vals[key]
	l.mu.RUnlock()

	if ok {
		l.changePos <- key
	}
	return v
}

func (l *lru) Set(ctx context.Context, key string, val any) {
	l.mu.Lock()
	_, ok := l.vals[key]
	if !ok {
		l.vals[key] = val
		l.pushCh <- newNode(key, val)
	}
	l.mu.Unlock()
}

func (l *lru) push() {
	insertToStart := func(n *lruNode) {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.head == nil {
			l.head = n
			l.tail = n
		} else {
			if l.len >= l.limit {
				l.popCh <- struct{}{}
			}
			n.right = l.head
			l.head.left = n
			l.head = n
		}
		l.len++
	}

	for {
		select {
		case n := <-l.pushCh:
			go insertToStart(n)
		case <-l.ctx.Done():
			return
		}
	}

}
func (l *lru) pop() {
	deleteFromEnd := func() {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.tail = l.tail.left
		l.tail.right = nil
	}

	for {
		select {
		case <-l.popCh:
			go deleteFromEnd()
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *lru) updatePos() {
	removeAndInsertAtStart := func(key string) {
		l.mu.Lock()
		defer l.mu.Unlock()

		curr := l.head

		for curr != nil {
			if curr.key == key && curr.left != nil {
				if curr.right != nil {
					curr.right.left = curr.left
				} else {
					l.tail = curr.left
				}

				curr.left.right = curr.right

				curr.left = nil
				curr.right = l.head
				l.head.left = curr
				l.head = curr

				break
			}

			curr = curr.right
		}
	}

	for {
		select {
		case key := <-l.changePos:
			go removeAndInsertAtStart(key)
		case <-l.ctx.Done():
			return
		}
	}
}

type lruNode struct {
	key string
	val any

	right *lruNode
	left  *lruNode
}

func NewLRU() Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := lru{
		limit:     DefaultLimit,
		ctx:       ctx,
		cancelCtx: cancel,
		vals:      make(map[string]any),
		pushCh:    make(chan *lruNode),
		popCh:     make(chan struct{}),
		changePos: make(chan string),
	}
	go c.pop()
	go c.push()
	go c.updatePos()
	return &c
}

func newNode(key string, val any) *lruNode {
	return &lruNode{
		key: key,
		val: val,
	}
}
