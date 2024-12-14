package linked

import (
	"iter"
	"sync"
)

type RingData[T any] interface {
	any
	GetID() string
	SetNode(T)
}

type RingNode[T RingData[T]] struct {
	D        *T
	next     *RingNode[T]
	previous *RingNode[T]
}

type Ring[T RingData[T]] struct {
	current *RingNode[T]
	head    *RingNode[T]

	length uint
	fill   uint
	mtx    *sync.RWMutex
}

// InitRing creates a new ring buffer
func InitRing[T RingData[T]](length uint) *Ring[T] {
	// create head node (first element)
	var head = &RingNode[T]{}
	var current = head

	// initialise ring
	ring := &Ring[T]{
		length:  length,
		fill:    0,
		mtx:     &sync.RWMutex{},
		current: head,
		head:    head,
	}

	// current = head
	for range length - 1 {
		// create new element and point it at current
		new := &RingNode[T]{
			D:        new(T),
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

	ring.current.D = &data
	ring.current = ring.current.next
	ring.inc()
}

func (ring *Ring[_]) inc() {
	if ring.fill < ring.length {
		ring.fill++
	}
}

func (ring *Ring[_]) Len() uint {
	return min(ring.length, ring.fill)
}

func (ring *Ring[T]) Get() iter.Seq[RingNode[T]] {
	if ring.length == 0 {
		return nil
	}

	ring.mtx.RLock()
	defer ring.mtx.RUnlock()

	var current *RingNode[T]

	return func(yield func(RingNode[T]) bool) {
		current = ring.head
		for range ring.Len() {
			if !yield(*current) {
				return
			}
			current = current.next
		}
	}
}
