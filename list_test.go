package list

import (
	"testing"
)

type testItem struct {
	name string
	age  int
}

func BenchmarkNewList(b *testing.B) {
	list := NewList[testItem]()
	for i := range b.N {
		list.Append(&testItem{
			name: "test",
			age:  i,
		})
	}
	b.Logf("length list %d", list.length)
}
func BenchmarkAllList(b *testing.B) {
	list := NewList[testItem]()
	for i := range 10 {
		list.Append(&testItem{
			name: "test",
			age:  i,
		})
	}
	b.Logf("length list %d", list.length)
	for d := range list.All() {
		b.Log(d)
	}
}
func BenchmarkNewSlice(b *testing.B) {
	list := []*testItem{}
	for i := range b.N {
		list = append(list, &testItem{
			name: "test",
			age:  i,
		})
	}
	b.Logf("length array %d", len(list))
}
