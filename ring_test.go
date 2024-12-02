package linked

import (
	"fmt"
	"testing"
)

type testRingItem struct {
	ID     string
	number int
	node   *testRingItem
}

func (i *testRingItem) GetID() string {
	return i.ID
}

func (i *testRingItem) SetNode(r *testRingItem) {
	i.node = r
}

func TestRing(t *testing.T) {
	r := InitRing[*testRingItem](1024)
	n := r.current
	for i := range r.length {
		t.Logf("%d: node: %p", i, n)
		t.Logf("previous: %p, next: %p", n.previous, n.next)
		n = n.next
	}
}
func TestRingOverlap(t *testing.T) {
	r := InitRing[*testRingItem](200)

	for i := range 300 {
		r.Add(&testRingItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})

	}
	for r := range r.Get() {
		t.Log("res", r.D.ID)
	}
}

func BenchmarkRing(b *testing.B) {
	r := InitRing[*testRingItem](1024)
	for i := 0; i < b.N; i++ {
		r.Add(&testRingItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
}
