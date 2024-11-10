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
	for d := range list.AllData() {
		b.Log(d)
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
