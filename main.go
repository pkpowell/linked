package list

import "iter"

type Node[T any] struct {
	data     *T
	next     *Node[T]
	previous *Node[T]
}

type List[T any] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) Update(data *T) {

}

func (l *List[T]) Append(data *T) {
	node := &Node[T]{data: data}
	switch l.length {
	case 0:
		// init new list, head and tail point to new node
		l.head = node
		l.tail = node

	default:
		// point tail.next at new node
		l.tail.next = node

		// point new node.previous at tail
		node.previous = l.tail

		// point tail at new node
		l.tail = node
	}

	l.length++
}

func (l *List[T]) Remove(data *T) {
	if l.head == nil || l.length == 0 {
		return
	}

	if l.head.data == data {
		l.head = l.head.next
		l.length--
		return
	}

	current := l.head
	for current.next != nil {
		if current.next.data == data {
			current.next = current.next.next
			l.length--

			return
		}
		current = current.next
	}
}

func (l *List[T]) All() iter.Seq[*T] {
	return func(yield func(*T) bool) {

		for l.head.next != nil {
			if !yield(l.head.next.data) {
				return
			}
		}
	}
}
