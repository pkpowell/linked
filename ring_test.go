package linked

import (
	"fmt"
	"testing"
)

func TestRing(t *testing.T) {
	r := InitRing[*testItem](256)
	n := r.current
	for i := range r.length {
		t.Logf("%d: node: %p", i, n)
		t.Logf("previous: %p, next: %p", n.previous, n.next)
		n = n.next
	}
}
func TestRingOverlap(t *testing.T) {
	r := InitRing[*testItem](200)

	for i := range 300 {
		r.Add(&testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})

	}
	for r := range r.Get() {
		t.Log("res", r.D.ID)
	}
}