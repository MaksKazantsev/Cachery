package Cachery

func newNode(key string, val any) *lruNode {
	return &lruNode{key: key, val: val}
}

func (l *lru) pushToStart(node *lruNode) {
	if node == l.head {
		return
	}

	if node.left != nil {
		node.left.right = node.right
	}

	if node.right != nil {
		node.right.left = node.left
	}

	if node == l.tail {
		l.tail = node.left
	}

	node.left = nil
	node.right = l.head
	l.head.left = node
	l.head = node

	if l.tail == nil {
		l.tail = node
	}
}
