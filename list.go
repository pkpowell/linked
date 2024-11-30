package linked

import (
	"iter"
	"sync"
)

type NodeData interface {
	any
	GetID() string
}

type Node[T NodeData] struct {
	D        T
	next     *Node[T]
	previous *Node[T]
	mtx      *sync.RWMutex
	list     *List[T]
}

type List[T NodeData] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
	mtx    *sync.RWMutex
}

// NewList creates a new list
func NewList[T NodeData]() *List[T] {
	return &List[T]{
		head:   nil,
		tail:   nil,
		length: 0,
		mtx:    &sync.RWMutex{},
	}
}

// InsertBefore adds a new node before a given node
func (list *List[T]) InsertBefore(data T, node *Node[T]) *Node[T] {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	newNode := list.newNode(data)

	switch list.length {
	case 0:
		list.head = newNode
		list.tail = newNode
	default:
		node.previous = newNode
		newNode.next = node
	}

	list.length++
	return newNode
}

func (list *List[T]) newNode(data T) *Node[T] {
	return &Node[T]{D: data, mtx: &sync.RWMutex{}, list: list}
}

// InsertAfter adds a new node after a given node
func (list *List[T]) InsertAfter(data T, node *Node[T]) *Node[T] {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	newNode := list.newNode(data)

	switch list.length {
	case 0:
		list.head = newNode
		list.tail = newNode
	default:
		if node.isTail() {
			newNode.makeTail()
		}
		newNode.previous = node
		node.next = newNode
	}

	list.length++
	return newNode
}

// Prepend adds a new node to the beginning of the list
func (list *List[T]) Prepend(data T) *Node[T] {
	// list.mtx.RLock()
	// defer list.mtx.RUnlock()

	return list.InsertBefore(data, list.head)
}

// Append adds a new node to the end of the list
func (list *List[T]) Append(data T) *Node[T] {
	// list.mtx.RLock()
	// defer list.mtx.RUnlock()

	return list.InsertAfter(data, list.tail)
}

func (node *Node[T]) isHead() bool {
	return node.list.head == node
}

func (node *Node[T]) isTail() bool {
	return node.list.tail == node
}

func (node *Node[T]) remove() {
	node.previous.next = node.next
	node.next.previous = node.previous
}

// removes node from list
func (node *Node[T]) Delete() {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	switch node.list.length {
	case 0:
		return

	case 1:
		//  list is now empty
		node.list.clear()
		return

	case 2:
		if node.isHead() { // if node to delete is current head
			node.next.makeHead()
		} else if node.isTail() { // if node to delete is current tail
			node.previous.makeTail()
		}

		node.list.tail = node.list.head
		node.list.length = 1

		return

	// list length 3 and longer
	default:
		if node.isHead() { // if node to delete is current head
			node.next.makeHead()
		} else if node.isTail() { // if node to delete is current tail
			node.previous.makeTail()
		} else { // if node to delete is in the middle
			node.remove()
		}

		node.list.length--
		return
	}
}

// returns list length
func (list *List[T]) Length() int {
	list.mtx.RLock()
	defer list.mtx.RUnlock()

	return list.length
}

// returns node with given id
func (list *List[T]) Get(id string) *Node[T] {
	list.mtx.RLock()
	defer list.mtx.RUnlock()

	current := list.head

	for range list.length {
		if current.D.GetID() == id {
			return current
		}

		current = current.next
	}
	return nil
}

func (list *List[T]) clear() {
	list.length = 0
	list.head = nil
	list.tail = nil
}

// DeleteNode deletes a node from the list
func (list *List[T]) DeleteNode(node *Node[T]) {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	switch list.length {
	case 0:
		return

	case 1:
		//  list is now empty (remove last element)
		list.clear()
		return

	case 2:
		if node.isHead() { // if node to delete is current head
			list.head = node.next
		} else if node.isTail() { // if node to delete is current tail
			list.head = node.previous
		}

		list.tail = list.head
		list.length = 1
		return

	default:
		if node.isHead() { // if node to delete is current head
			node.next.makeHead()
		} else if node.isTail() { // if node to delete is current tail
			node.previous.makeTail()
		} else { // if node to delete is in the middle
			node.remove()
		}

		list.length--
		return
	}
}

func (node *Node[T]) makeHead() {
	node.list.head.previous = node
	node.list.head = node
}

func (node *Node[T]) makeTail() {
	node.list.tail.next = node
	node.list.tail = node
}

// AllNodes returns all nodes in the list
func (list *List[T]) AllNodes() iter.Seq[*Node[T]] {
	list.mtx.RLock()
	defer list.mtx.RUnlock()

	return func(yield func(*Node[T]) bool) {
		if list.length == 0 {
			return
		}

		current := list.head

		for range list.length {
			if !yield(current) {
				return
			}
			current = current.next
		}
	}
}

// AllData returns all data in the list (without nodes)
func (list *List[T]) AllData() iter.Seq[T] {
	list.mtx.RLock()
	defer list.mtx.RUnlock()

	return func(yield func(T) bool) {
		if list.length == 0 {
			return
		}
		current := list.head

		for range list.length {
			if !yield(current.D) {
				return
			}
			current = current.next
		}
	}
}
