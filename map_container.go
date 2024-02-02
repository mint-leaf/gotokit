package main

import "sort"

type Tuple[K comparable, V any] struct {
	Key   K
	Value V
}

type Cpr[T comparable] struct {
	data []T
	less func(i, j T) bool
}

func (c Cpr[T]) Len() int {
	return len(c.data)
}

func (c Cpr[T]) Less(i, j int) bool {
	return c.less(c.data[i], c.data[j])
}

func (c Cpr[T]) Swap(i, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
}

func MapFilter[M ~map[K]V, K comparable, V any](m M, fs ...func(ele K) bool) M {
	c := make(M)
	for k := range m {
		for i := range fs {
			if fs[i](k) {
				c[k] = m[k]
			}
		}
	}
	return c
}

func MapKeys[M ~map[K]V, S []K, K comparable, V any](m M) S {
	c := make([]K, len(m))
	index := 0
	for k := range m {
		c[index] = k
	}
	return c
}

func MapKeysSort[M ~map[K]V, S []K, K comparable, V any](m M, f func(i, j K) bool) S {
	keys := MapKeys(m)
	sort.Sort(Cpr[K]{data: keys, less: f})
	return keys
}

func Load[M ~map[K]V, K comparable, V any](m M, key K, d V) V {
	if m == nil || len(m) == 0 {
		return d
	}
	v, ok := m[key]
	if !ok {
		return d
	}
	return v
}

func MMapper[M ~map[K]V, S []N, K comparable, T any, V any, N any](m M, f func(key K, value V) N) S {
	c := make([]N, len(m))
	index := 0
	for k, v := range m {
		c[index] = f(k, v)
		index++
	}
	return c
}

func MapToSlice[M ~map[K]V, K comparable, V any](m M) []*Tuple[K, V] {
	s := make([]*Tuple[K, V], len(m))
	index := 0
	for k, v := range m {
		s[index] = &Tuple[K, V]{Key: k, Value: v}
		index++
	}
	return s
}
