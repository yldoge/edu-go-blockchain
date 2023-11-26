package types

import (
	"fmt"
	"reflect"
)

type List[T any] struct {
	Data []T
}

func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

func (l *List[T]) Get(i int) T {
	if i > len(l.Data)-1 {
		err := fmt.Errorf("the given index (%d) is higher than the length (%d)", i, len(l.Data))
		panic(err)
	}
	return l.Data[i]
}

func (l *List[T]) Insert(v T) {
	l.Data = append(l.Data, v)
}

func (l *List[T]) Clear() {
	l.Data = []T{}
}

// Return -1 is v is not in the list
func (l *List[T]) GetIndex(v T) int {
	for i, val := range l.Data {
		if reflect.DeepEqual(v, val) {
			return i
		}
	}
	return -1
}

func (l *List[T]) Remove(v T) {
	i := l.GetIndex(v)
	if i == -1 {
		return
	}
	l.Pop(i)
}

func (l *List[T]) Pop(i int) {
	l.Data = append(l.Data[:i], l.Data[i+1:]...)
}

func (l *List[T]) Contains(v T) bool {
	for _, val := range l.Data {
		if reflect.DeepEqual(v, val) {
			return true
		}
	}
	return false
}

func (l *List[T]) Last() T {
	return l.Data[len(l.Data)-1]
}

func (l *List[T]) Len() int {
	return len(l.Data)
}
