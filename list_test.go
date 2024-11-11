package linked

import (
	"testing"
)

type testItem struct {
	text   string
	number int
}

func BenchmarkNewList(b *testing.B) {
	list := NewList[testItem]()
	for i := range b.N {
		list.Append(&testItem{
			text:   "test",
			number: i,
		})
	}
	b.Logf("length list %d", list.length)
}

func BenchmarkAllList(b *testing.B) {
	list := NewList[testItem]()
	for i := range b.N {
		list.Append(&testItem{
			text:   "test",
			number: i,
		})
	}
	b.Logf("length list %d", list.length)
	for _ = range list.AllData() {
		// b.Log(d)
	}
}

func BenchmarkDelete(b *testing.B) {
	list := NewList[testItem]()
	for i := range b.N {
		list.Append(&testItem{
			text:   "test",
			number: i,
		})
	}
	for d := range list.AllNodes() {
		if d.D.number%200 == 0 {
			list.DeleteNode(d)
		}
		// b.Log(d)
	}
	b.Logf("length list %d", list.length)
}

func BenchmarkNewSlice(b *testing.B) {
	list := []*testItem{}
	for i := range b.N {
		list = append(list, &testItem{
			text:   "test",
			number: i,
		})
	}
	b.Logf("length array %d", len(list))
}

func TestList(t *testing.T) {
	l := NewList[testItem]()
	for i := range 10 {
		l.Append(&testItem{
			text:   "test",
			number: i,
		})
	}
	for d := range l.AllData() {
		t.Log(d)
	}
}

func TestEmptyList(t *testing.T) {
	l := NewList[testItem]()
	if l.length != 0 {
		t.Errorf("Expected empty list length 0, got %d", l.length)
	}
}

func TestAppendAndLength(t *testing.T) {
	l := NewList[testItem]()
	l.Append(&testItem{text: "first", number: 1})
	l.Append(&testItem{text: "second", number: 2})

	if l.length != 2 {
		t.Errorf("Expected length 2, got %d", l.length)
	}
}

func TestDeleteFirstNode(t *testing.T) {
	l := NewList[testItem]()
	firstNode := l.Append(&testItem{text: "first", number: 1})
	secondNode := l.Append(&testItem{text: "second", number: 2})
	thirdNode := l.Append(&testItem{text: "third", number: 3})
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
	l := NewList[testItem]()
	l.Append(&testItem{text: "first", number: 1})
	l.Append(&testItem{text: "second", number: 2})

	var lastNode *Node[testItem]
	for node := range l.AllNodes() {
		lastNode = node
	}
	l.DeleteNode(lastNode)

	if l.length != 1 {
		t.Errorf("Expected length 1 after deletion, got %d", l.length)
	}

	remaining := l.Length()
	if remaining != 1 {
		t.Errorf("Expected remaining number 1, got %d", remaining)
	}
}

func TestAllDataEmptyList(t *testing.T) {
	l := NewList[testItem]()
	count := 0
	for range l.AllData() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected no items in empty list, got %d", count)
	}
}
