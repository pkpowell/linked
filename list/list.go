package linked

import (
	"fmt"
	"iter"
)

type NodeData interface {
	any
	GetID() string
}

type Node[T NodeData] struct {
	D        T
	next     *Node[T]
	previous *Node[T]
}

type List[T NodeData] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

// NewList creates a new empty list
func NewList[T NodeData]() *List[T] {
	return &List[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

// InsertBefore adds a new node before a given node
func (list *List[T]) InsertBefore(data T, n *Node[T]) *Node[T] {
	node := &Node[T]{D: data}
	if list.length == 0 {
		list.head = node
		list.tail = node
		return node
	}
	n.previous = node
	node.next = n
	list.length++
	return node
}

// InsertAfter adds a new node after a given node
func (list *List[T]) InsertAfter(data T, n *Node[T]) *Node[T] {
	node := &Node[T]{D: data}
	if list.length == 0 {
		list.head = node
		list.tail = node
		return node
	}
	node.previous = n
	n.next = node
	list.length++

	return node
}

// Prepend adds a new node to the beginning of the list
func (list *List[T]) Prepend(data T) *Node[T] {
	node := &Node[T]{D: data}
	switch list.length {
	case 0:
		// init new list, head and tail point to new node
		list.head = node
		list.tail = node

	default:
		// point previous head at new node
		list.head.previous = node
		// point new node at previous head
		node.next = list.head

		// point tail at new node
		list.head = node
	}

	list.length++
	return node
}

// Append adds a new node to the end of the list
func (list *List[T]) Append(data T) *Node[T] {
	node := &Node[T]{D: data}
	switch list.length {
	case 0:
		// init new list, head and tail point to new node
		list.head = node
		list.tail = node

	default:
		// point tail.next at new node
		list.tail.next = node

		// point new node.previous at tail
		node.previous = list.tail

		// point tail at new node
		list.tail = node
	}

	list.length++
	return node
}

// type RunFunc[T any] func(*Node[T])

// func (node *Node[T]) Run(fu RunFunc[T]) {
// 	fu(node)
// }

// returns list length
func (list *List[T]) Length() int {
	return list.length
}

// returns node with given id
func (list *List[T]) Get(id string) *Node[T] {
	current := list.head
	data := current.D
	fmt.Printf("data %v\n", data)
	for {
		if current == nil {
			break
		}
		if data.GetID() == id {
			fmt.Printf("found %s id %s\n", data.GetID(), id)
			return current
		}
		current = current.next
	}
	return nil
}

// DeleteNode deletes a node from the list
func (list *List[T]) DeleteNode(node *Node[T]) {
	switch list.length {
	case 0:
		return
	case 1:
		//  list is now empty
		list.length = 0
		list.head = nil
		list.tail = nil
		return

	case 2:
		if node == list.head { // if node to delete is current head
			list.head = node.next
		} else if node.next == list.tail { // if node to delete is current tail
			list.head = node.previous
		}
		list.head.previous = nil
		list.head.next = nil
		list.tail = list.head
		list.length = 1
		return

	default:
		if node == list.head { // if node to delete is current head
			node.next.previous = nil
			list.head = node.next

		} else if node.next == list.tail { // if node to delete is current tail
			list.tail = node.previous
			list.tail.next = nil

		} else { // if node to delete is in the middle
			node.previous.next = node.next
			node.next.previous = node.previous
		}
		list.length--
		return
	}
}

// AllNodes returns all nodes in the list
func (list *List[T]) AllNodes() iter.Seq[*Node[T]] {
	return func(yield func(*Node[T]) bool) {
		if list.length == 0 {
			return
		}
		current := list.head

		for {
			if current == nil {
				return
			}

			if !yield(current) {
				return
			}
			current = current.next
		}
	}
}

// AllData returns all data in the list (without nodes)
func (list *List[T]) AllData() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		if list.length == 0 {
			return
		}
		current := list.head

		for {
			if current == nil {
				return
			}

			if !yield(&current.D) {
				return
			}
			current = current.next
		}
	}
}
