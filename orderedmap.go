package xcontainer

import "iter"

type entry[K comparable, V any] struct {
	key K
	val V
}

func enter[K comparable, V any](key K, val V) entry[K, V] {
	return entry[K, V]{key: key, val: val}
}

// OrderedMap is similar to Go's built-in map type, but maintains
// insertion order for iteration. In other words, if element b is
// inserted after element a is already in the map, b will always be
// yieled after a during iteration.
//
// Unlike a built-in Go map, a zero-value OrderedMap is empty and
// ready to use.
type OrderedMap[K comparable, V any] struct {
	m    map[K]*List[entry[K, V]]
	head *List[entry[K, V]]
}

func (m *OrderedMap[K, V]) init() {
	if m.m != nil {
		return
	}

	m.m = make(map[K]*List[entry[K, V]])
}

// Set sets the provided key to val in m.
func (m *OrderedMap[K, V]) Set(key K, val V) {
	m.init()

	node, ok := m.m[key]
	if ok {
		node.Val.val = val
		return
	}

	m.head = m.head.InsertBefore(enter(key, val))
	m.m[key] = m.head.Prev()
}

// Lookup returns the value associated with the key and a boolean
// indicating if there was one or not.
func (m *OrderedMap[K, V]) Lookup(key K) (val V, ok bool) {
	node, ok := m.m[key]
	if !ok {
		return val, false
	}
	return node.Val.val, true
}

// Delete removes the value associated with key from the map. If no
// such value exists, it does nothing.
func (m *OrderedMap[K, V]) Delete(key K) {
	node, ok := m.m[key]
	if !ok {
		return
	}

	prev := node.Remove()
	if node == m.head {
		m.head = prev.Next()
	}

	delete(m.m, key)
}

// Clear deletes everything from m, resulting in an empty map.
func (m *OrderedMap[K, V]) Clear() {
	clear(m.m)
	m.head = nil
}

// All returns an iter.Seq that yields key value pairs in insertion
// order.
func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for entry := range m.head.Values() {
			if !yield(entry.key, entry.val) {
				return
			}
		}
	}
}

// Keys returns an iter.Seq that yields keys in insertion order.
func (m *OrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iter.Seq that yields values in insertion order.
func (m *OrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}
