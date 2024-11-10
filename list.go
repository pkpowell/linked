package list

import "iter"

type Node[T comparable] struct {
	data     *T
	next     *Node[T]
	previous *Node[T]
}

type List[T comparable] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

func NewList[T comparable]() *List[T] {
	return &List[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

func (list *List[T]) Append(data *T) {
	node := &Node[T]{data: data}
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
}

func (list *List[T]) DeleteNode(node *Node[T]) {
	if list.head == nil || list.length == 0 {
		return
	}

	switch list.length {
	case 1:
		// list is now empty
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

func (list *List[T]) AllNodes() iter.Seq[*Node[T]] {
	return func(yield func(*Node[T]) bool) {
		if list.head == nil {
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

func (list *List[T]) AllData() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		if list.head == nil {
			return
		}
		current := list.head

		for {
			if current == nil {
				return
			}

			if !yield(current.data) {
				return
			}
			current = current.next
		}

	}
}
