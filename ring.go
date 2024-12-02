package linked

import (
	"iter"
	"sync"
)

type RingNode[T NodeData] struct {
	D        T
	next     *RingNode[T]
	previous *RingNode[T]
}

type Ring[T NodeData] struct {
	current *RingNode[T]
	head    *RingNode[T]

	length int
	fill   int
	mtx    *sync.RWMutex
}

// InitRing creates a new ring buffer
func InitRing[T NodeData](length int) *Ring[T] {
	// create head node (first element)
	var head = &RingNode[T]{}

	// initialise ring
	ring := &Ring[T]{
		length:  length,
		fill:    0,
		mtx:     &sync.RWMutex{},
		current: head,
		head:    head,
	}

	// current is a temp variable
	current := head

	for range length - 1 {
		var d T
		// create new element and point it at current
		new := &RingNode[T]{
			D:        d,
			previous: current,
		}
		// point current at new element
		current.next = new

		// set current to new element
		current = new
	}
	// point last element at head and .v.v
	head.previous = current
	current.next = head

	return ring
}

func (ring *Ring[T]) Add(data T) {
	ring.mtx.Lock()
	defer ring.mtx.Unlock()

	ring.current.D = data
	ring.current = ring.current.next
}

func (ring *Ring[T]) Get() iter.Seq[*RingNode[T]] {
	ring.mtx.RLock()
	defer ring.mtx.RUnlock()

	if ring.length == 0 {
		return nil
	}

	return func(yield func(*RingNode[T]) bool) {

		current := ring.head

		for range ring.length {
			if !yield(current) {
				return
			}
			current = current.next
		}
	}
}
