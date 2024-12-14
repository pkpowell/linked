package linked

import (
	"iter"
	"strconv"
	"sync"
)

type Data interface {
	comparable
	GetID() string
}

type Node[T Data] struct {
	D        T
	next     *Node[T]
	previous *Node[T]
	list     *List[T]

	mtx *sync.RWMutex
}

type List[T Data] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
	mtx    *sync.RWMutex
}

// returns a new list
func NewList[T Data]() *List[T] {
	return &List[T]{
		head:   nil,
		tail:   nil,
		length: 0,
		mtx:    &sync.RWMutex{},
	}
}

// inserts a new node before a given node
func (list *List[T]) InsertBefore(data T, node *Node[T]) *Node[T] {
	newNode := list.newNode(data)

	switch list.length {
	case 0:
		list.setHead(newNode)
		list.setTail(newNode)

	default:
		node.setPrevious(newNode)
		newNode.next = node
	}

	list.inc()
	return newNode
}

// returns a new node
func (list *List[T]) newNode(data T) *Node[T] {
	return &Node[T]{
		D:    data,
		mtx:  &sync.RWMutex{},
		list: list,
	}
}

// inserts a new node after a given node
func (list *List[T]) InsertAfter(data T, node *Node[T]) *Node[T] {
	newNode := list.newNode(data)

	switch list.length {
	case 0:
		list.setHead(newNode)
		list.setTail(newNode)
	default:
		if node.isTail() {
			newNode.makeTail()
		}
		newNode.previous = node
		node.setNext(newNode)
	}

	list.inc()
	return newNode
}

// sets the next node
func (node *Node[T]) setNext(newNode *Node[T]) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.next = newNode
}

// sets the previous node
func (node *Node[T]) setPrevious(newNode *Node[T]) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.previous = newNode
}

// Prepend adds a new node to the beginning of the list
func (list *List[T]) Prepend(data T) *Node[T] {
	return list.InsertBefore(data, list.head)
}

// Append adds a new node to the end of the list
func (list *List[T]) Append(data T) *Node[T] {
	return list.InsertAfter(data, list.tail)
}

// sets the head node
func (list *List[T]) setHead(node *Node[T]) {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.head = node
}

// sets the tail node
func (list *List[T]) setTail(node *Node[T]) {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.tail = node
}

// sets the length of the list
func (list *List[_]) setLength(l int) {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.length = l
}

// increments the the list length
func (list *List[_]) inc() {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.length++
}

// decrements the the list length
func (list *List[_]) dec() {
	list.mtx.Lock()
	defer list.mtx.Unlock()

	list.length--
}

// returns true if the node is the current head
func (node *Node[_]) isHead() bool {
	switch true {
	case node == nil, node.list == nil:
		return false

	case node.list.length == 1:
		return true

	default:
		return node.list.head == node
	}
}

// returns true if the node is the current tail
func (node *Node[_]) isTail() bool {
	return node.list.tail == node
}

// remove removes itself from the list
func (node *Node[_]) remove() {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.previous.setNext(node.next)
	node.next.setPrevious(node.previous)
}

// removes node from list and sets new head or tail
func (node *Node[_]) Delete() {
	switch node.list.length {
	case 0:
		return

	case 1:
		//  list is now empty
		node.list.setLength(0)
		return

	case 2:
		if node.isHead() { // if node to delete is current head
			node.next.makeHead()
		} else if node.isTail() { // if node to delete is current tail
			node.previous.makeTail()
		}

		node.list.setTail(node.list.head)
		node.list.setLength(1)

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

		node.list.dec()
		return
	}
}

// returns list length
func (list *List[_]) Len() int {
	list.mtx.RLock()
	defer list.mtx.RUnlock()

	return list.length
}

// returns list length as string
func (list *List[_]) LenStr() string {
	list.mtx.RLock()
	defer list.mtx.RUnlock()

	return strconv.Itoa(list.length)
}

// returns node with given id (T must implement GetID() string).
// returns nil if node not found
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

// deletes a node from the list
func (list *List[T]) DeleteNode(node *Node[T]) {
	switch list.length {
	case 0:
		return

	case 1:
		//  list is now empty (remove last element)
		list.setLength(0)
		return

	case 2:
		if node.isHead() { // if node to delete is current head
			list.setHead(node.next)
		} else if node.isTail() { // if node to delete is current tail
			list.setTail(node.previous)
		}

		list.setLength(1)
		return

	default:
		if node.isHead() { // if node to delete is current head
			node.next.makeHead()
		} else if node.isTail() { // if node to delete is current tail
			node.previous.makeTail()
		} else { // if node to delete is in the middle
			node.remove()
		}

		list.dec()
		return
	}
}

// makeHead promotes the node to head
func (node *Node[_]) makeHead() {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	// point current head at new head
	node.list.head.previous = node
	// set new head
	node.list.head = node
}

// makeTail promotes the node to tail
func (node *Node[_]) makeTail() {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	// point current tail at new tail
	node.list.tail.next = node
	// set new tail
	node.list.tail = node
}

// returns all nodes in the list
func (list *List[T]) AllNodes() iter.Seq[*Node[T]] {
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

// returns all data in the list (without nodes)
func (list *List[T]) AllData() iter.Seq[T] {
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
