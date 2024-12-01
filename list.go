package linked

import (
	"fmt"
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
		head: &Node[T]{
			mtx: &sync.RWMutex{},
			list: &List[T]{
				head:   &Node[T]{},
				tail:   &Node[T]{},
				length: 0,
				mtx:    &sync.RWMutex{},
			},
		},
		tail: &Node[T]{
			mtx: &sync.RWMutex{},
			list: &List[T]{
				head:   &Node[T]{},
				tail:   &Node[T]{},
				length: 0,
				mtx:    &sync.RWMutex{},
			},
		},
		length: 0,
		mtx:    &sync.RWMutex{},
	}
}

// InsertBefore adds a new node before a given node
func (list *List[T]) InsertBefore(data T, node *Node[T]) *Node[T] {
	list.mtx.RLock()
	defer list.mtx.RUnlock()

	newNode := list.makeNode(data)

	switch list.length {
	case 0:
		list.headUpdate(newNode)
		list.tailUpdate(newNode)
	default:
		node.previous.update(newNode)
		newNode.next.update(node)
	}

	list.length++
	return newNode
}

func (list *List[T]) makeNode(data T) *Node[T] {
	return &Node[T]{
		D:    data,
		mtx:  &sync.RWMutex{},
		list: list,
	}
}

// InsertAfter adds a new node after a given node
func (list *List[T]) InsertAfter(data T, node *Node[T]) *Node[T] {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	newNode := list.makeNode(data)

	switch list.length {
	case 0:
		list.headUpdate(newNode)
		list.tailUpdate(newNode)
	default:
		if node.isTail() {
			newNode.makeTail()
		}
		newNode.previous.update(node)
		node.next.update(newNode)
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
	if node == nil {
		return false
	}
	return node.list.head == node
}

func (node *Node[T]) isTail() bool {
	if node == nil {
		return false
	}
	return node.list.tail == node
}

func (node *Node[T]) remove() {
	if node == nil {
		return
	}
	node.previous.next.update(node.next)
	node.next.previous.update(node.previous)
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

		node.list.tail.update(node.list.head)
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
	// list.mtx.RLock()
	// defer list.mtx.RUnlock()

	current := list.head

	for range list.length {
		if current.D.GetID() == id {
			return current
		}

		current.update(current.next)
	}
	return nil
}

func (list *List[T]) headUpdate(node *Node[T]) {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.head = node
}

func (list *List[T]) tailUpdate(node *Node[T]) {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.head = node
}

func (list *List[T]) clear() {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.length = 0
	list.head = nil
	list.tail = nil
}

// DeleteNode deletes a node from the list
func (list *List[T]) DeleteNode(node *Node[T]) {
	// list.mtx.Lock()
	// defer list.mtx.Unlock()

	switch list.length {
	case 0:
		return

	case 1:
		//  list is now empty (remove last element)
		list.clear()
		return

	case 2:
		if node.isHead() { // if node to delete is current head
			list.headUpdate(node.next)
			// list.head.update(node.next)
		} else if node.isTail() { // if node to delete is current tail
			list.headUpdate(node.previous)
			// list.head.update(node.previous)
		}

		list.tailUpdate(list.head)
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
	node.list.head.previous.update(node)
	node.list.headUpdate(node)
}

func (node *Node[T]) makeTail() {
	fmt.Println("make tail", node.list.tail.D, node.list.tail.next)
	node.list.tail.next.update(node)
	node.list.tailUpdate(node)
}

func (node *Node[T]) update(newNode *Node[T]) {
	if node != nil {
		node.mtx.Lock()
		defer node.mtx.Unlock()
	}
	// fmt.Println("node", node)
	// node.mtx.Lock()
	// defer node.mtx.Unlock()

	node = newNode
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
			current.update(current.next)
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
			current.update(current.next)
		}
	}
}
