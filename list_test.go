package linked

import (
	"fmt"
	"testing"
	"time"
)

type testItem struct {
	ID     string
	number int
}

func (i *testItem) GetID() string {
	return i.ID
}

func BenchmarkNewList(b *testing.B) {
	list := NewList[*testItem]()
	for i := range b.N {
		list.Append(&testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
	b.Logf("length list %d", list.length)
}

func BenchmarkAllList(b *testing.B) {
	list := NewList[*testItem]()
	for i := range b.N {
		list.Append(&testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
	b.Logf("length list %d", list.length)
	for _ = range list.AllData() {
		// b.Log(d)
	}
}

func BenchmarkDelete(b *testing.B) {
	list := NewList[*testItem]()
	for i := range b.N {
		list.Append(&testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
	for d := range list.AllNodes() {
		if d.D.number%200 == 0 {
			list.DeleteNode(d)
			// d.Delete()
		}
		// b.Log(d)
	}
	b.Logf("length list %d", list.length)
}

func BenchmarkNewSlice(b *testing.B) {
	arr := []*testItem{}
	for i := range b.N {
		arr = append(arr, &testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
	b.Logf("length array %d", len(arr))
}

func TestList(t *testing.T) {
	l := NewList[*testItem]()
	for i := range 10 {
		l.Append(&testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
	for d := range l.AllData() {
		t.Log(*d)
	}
}

func checkList(l *List[*testItem], t *testing.T) {
	for d := range l.AllData() {
		d = &testItem{
			ID:     fmt.Sprintf("%s-test-", d.ID),
			number: d.number * 2,
		}
		t.Log(*d)
	}
}

func updateList(l *List[*testItem], t *testing.T) {
	fmt.Println("updating...")
	for d := range l.AllData() {
		d = &testItem{
			ID:     fmt.Sprintf("%s-test-", d.ID),
			number: d.number * 2,
		}
		t.Log(*d)
	}
}
func TestListConcurrency(t *testing.T) {
	update := time.NewTicker(time.Microsecond * 1000)
	check := time.NewTicker(time.Microsecond * 1100)
	done := time.NewTimer(time.Second * 60)
	l := NewList[*testItem]()
	for i := range 10 {
		fmt.Println("appending...")
		l.Append(&testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
	for {
		select {
		case <-done.C:
			fmt.Println("Done")
			return
		case <-update.C:
			go updateList(l, t)
		case <-check.C:
			go checkList(l, t)
		}
	}

}

func TestEmptyList(t *testing.T) {
	l := NewList[*testItem]()
	if l.length != 0 {
		t.Errorf("Expected empty list length 0, got %d", l.length)
	}
}

func TestAppendAndLength(t *testing.T) {
	l := NewList[*testItem]()
	l.Append(&testItem{ID: "first", number: 1})
	l.Append(&testItem{ID: "second", number: 2})

	if l.length != 2 {
		t.Errorf("Expected length 2, got %d", l.length)
	}
}
func TestGet(t *testing.T) {
	l := NewList[*testItem]()
	ids := []string{"first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth", "ninth", "tenth"}
	for idx, id := range ids {
		node := l.Append(&testItem{ID: id, number: idx})

		get := l.Get(id)
		if get == node {
			t.Logf("Found node %v", get.D)
		}
	}
}

func TestDeleteFirstNode(t *testing.T) {
	l := NewList[*testItem]()
	firstNode := l.Append(&testItem{ID: "first", number: 1})
	secondNode := l.Append(&testItem{ID: "second", number: 2})
	thirdNode := l.Append(&testItem{ID: "third", number: 3})
	length := l.length
	// firstNode := <-l.AllNodes()
	l.DeleteNode(firstNode)

	if l.length != length-1 {
		t.Errorf("Expected length 1 after deletion, got %d", l.length)
	}
	t.Logf("length list %d", l.length)
	t.Logf("secondNode %v", secondNode.D)
	t.Logf("thirdNode %v", thirdNode.D)
	// remaining := l.length
	// if remaining != 2 {
	// 	t.Errorf("Expected remaining number 2, got %d", remaining.number)
	// }
}

func TestDeleteLastNode(t *testing.T) {
	l := NewList[*testItem]()
	l.Append(&testItem{ID: "first", number: 1})
	l.Append(&testItem{ID: "second", number: 2})
	l.Append(&testItem{ID: "third", number: 3})
	l.Append(&testItem{ID: "forth", number: 4})

	var lastNode *Node[*testItem]
	for node := range l.AllNodes() {
		lastNode = node
	}
	l.DeleteNode(lastNode)

	if l.length != 3 {
		t.Errorf("Expected length 1 after deletion, got %d", l.length)
	}

	remaining := l.Length()
	if remaining != 3 {
		t.Errorf("Expected remaining number 1, got %d", remaining)
	}
	for n := range l.AllNodes() {
		t.Log("list", n.D.ID)
	}
}
func TestRing(t *testing.T) {
	r := InitRing[*testItem](256)
	n := r.current
	for i := range r.length {
		t.Logf("%d: node: %p", i, n)
		t.Logf("previous: %p, next:%p", n.previous, n.next)
		n = n.next
	}
}

func TestGetID(t *testing.T) {
	l := NewList[*testItem]()
	for i := range 10 {
		l.Append(&testItem{
			ID:     fmt.Sprintf("%d-test-", i),
			number: i,
		})
	}
	// l.Append(&testItem{ID: "first", number: 1})
	// l.Append(&testItem{ID: "second", number: 2})
	for n := range l.AllNodes() {
		t.Logf("%s, %s", n.D.ID, n.D.GetID())
	}
}

func TestAllDataEmptyList(t *testing.T) {
	l := NewList[*testItem]()
	count := 0
	for range l.AllData() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected no items in empty list, got %d", count)
	}
}
