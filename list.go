package xcontainer

import (
	"iter"
	"slices"
)

// List is a doubly-linked circular list. Any element in the list is a
// reference to the whole list. A zero-value List is not valid, but a
// nil *List is and is the correct representation of an empty list.
type List[T any] struct {
	prev, next *List[T]
	Val        T
}

// NewList returns a List containing the provided values.
func NewList[T any](vals ...T) *List[T] {
	return NewListFromSeq(slices.Values(vals))
}

// NewListFromSeq returns a List containing the values yielded by seq.
func NewListFromSeq[T any](seq iter.Seq[T]) (list *List[T]) {
	return list.InsertSeqBefore(seq)
}

// Prev returns the node before list.
func (list *List[T]) Prev() *List[T] {
	if list == nil {
		return nil
	}
	return list.prev
}

// Next returns the node after list.
func (list *List[T]) Next() *List[T] {
	if list == nil {
		return nil
	}
	return list.next
}

// InsertBefore inserts v as a new node in between list and the node
// before it. It returns list, possibly modified. Similar to [append],
// the existing list should always be replaced with the result of this
// function.
func (list *List[T]) InsertBefore(v T) *List[T] {
	node := &List[T]{next: list, Val: v}
	if list != nil {
		node.prev, list.prev = list.prev, node
		return list
	}
	node.prev, node.next = node, node
	return node
}

// InsertAfter inserts v as a new node in between list and the node
// after it. It returns list, possibly modified. Similar to [append],
// the existing list should always be replaced with the result of this
// function.
func (list *List[T]) InsertAfter(v T) *List[T] {
	node := &List[T]{prev: list, Val: v}
	if list != nil {
		node.next, list.next = list.next, node
		return list
	}
	node.prev, node.next = node, node
	return node
}

// Remove removes list from the List that it represents. It returns
// the node before list. list is no longer valid after a call to
// Remove.
func (list *List[T]) Remove() *List[T] {
	if list == nil {
		return nil
	}

	list.prev.next, list.next.prev = list.next, list.prev
	return list.prev
}

// InsertSeqBefore inserts the values yielded by seq as new nodes
// before list. The order of the elements in the List will be the same
// as they were in seq. The behavior of this function is otherwise the
// same as that of [InsertBefore].
func (list *List[T]) InsertSeqBefore(seq iter.Seq[T]) *List[T] {
	for v := range seq {
		list = list.InsertBefore(v)
	}
	return list
}

// InsertSeqAfter inserts the values yielded by seq as new nodes
// after list. The order of the elements in the List will be the same
// as they were in seq. The behavior of this function is otherwise the
// same as that of [InsertBefore].
func (list *List[T]) InsertSeqAfter(seq iter.Seq[T]) *List[T] {
	tail := list
	for v := range seq {
		tail = tail.InsertAfter(v).Next()
		if list == nil {
			list = tail
		}
	}
	return list
}

func (list *List[T]) seq(advance func(*List[T]) *List[T]) iter.Seq[*List[T]] {
	return func(yield func(*List[T]) bool) {
		if list == nil {
			return
		}

		cur := list
		for {
			if !yield(cur) {
				return
			}
			cur = advance(cur)
			if cur == list {
				return
			}
		}
	}
}

// All returns an iter.Seq that yields the nodes of the List in next
// order, starting with list itself, once each.
//
// For most situations, [Values] is more convenient.
func (list *List[T]) All() iter.Seq[*List[T]] {
	return list.seq((*List[T]).Next)
}

// Backward returns an iter.Seq that yields the nodes of the List in
// prev order, starting with list itself, once each.
//
// For most situations, [ValuesBackward] is more convenient.
func (list *List[T]) Backward() iter.Seq[*List[T]] {
	return list.seq((*List[T]).Prev)
}

// Values returns an iter.Seq that yields the values of each node in
// the List in next order, starting with that of list itself, once
// each.
func (list *List[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := range list.All() {
			if !yield(node.Val) {
				return
			}
		}
	}
}

// ValuesBackward returns an iter.Seq that yields the values of each
// node in the List in prev order, starting with that of list itself,
// once each.
func (list *List[T]) ValuesBackward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := range list.Backward() {
			if !yield(node.Val) {
				return
			}
		}
	}
}
