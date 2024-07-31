/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"runtime"
	"sync"
)

// SafeMap 多线程安全的 Map
type SafeMap[T any] struct {
	m map[string]T
	l sync.RWMutex
}

// NewSafeMap 创建一个新的 SafeMap
func NewSafeMap[T any]() *SafeMap[T] {
	return &SafeMap[T]{
		m: make(map[string]T),
		l: sync.RWMutex{},
	}
}

// Get returns the value stored under the given key.
func (sm *SafeMap[T]) Get(key string) (T, bool) {
	sm.l.RLock()
	value, ok := sm.m[key]
	sm.l.RUnlock()
	return value, ok
}

// Set stores the given value under the given key.
func (sm *SafeMap[T]) Set(key string, value T) {
	sm.l.Lock()
	sm.m[key] = value
	sm.l.Unlock()
}

// Del deletes the value stored under the given key.
func (sm *SafeMap[T]) Del(key string) {
	sm.l.Lock()
	delete(sm.m, key)
	sm.l.Unlock()
}

// Keys returns all keys in the map.
func (sm *SafeMap[T]) Keys() []string {
	sm.l.RLock()
	keys := make([]string, 0, len(sm.m))
	for key := range sm.m {
		keys = append(keys, key)
	}
	sm.l.RUnlock()
	return keys
}

// Range calls the given function for each key-value pair in the map.
func (sm *SafeMap[T]) Range(f func(key string, value T) bool) {
	sm.l.RLock()
	defer sm.l.RUnlock()

	for key, value := range sm.m {
		if !f(key, value) {
			return
		}
	}
}

func (sm *SafeMap[T]) Len() int {
	sm.l.RLock()
	defer sm.l.RUnlock()

	return len(sm.m)

}

// SharedSafeMap is a map that can be safely shared by multiple goroutines.
type SharedSafeMap[T any] struct {
	buckets []*SafeMap[T]
}

// NewSharedSafeMap creates a new SharedSafeMap.
func NewSharedSafeMap[T any]() *SharedSafeMap[T] {
	n := runtime.GOMAXPROCS(0)
	buckets := make([]*SafeMap[T], n)
	for i := range buckets {
		buckets[i] = NewSafeMap[T]()
	}
	return &SharedSafeMap[T]{buckets: buckets}

}

// Load returns the value stored under the given key.
func (sm *SharedSafeMap[T]) Load(key string) (T, bool) {
	i := share(key, len(sm.buckets))
	return sm.buckets[i].Get(key)
}

// Store stores the given value under the given key.
func (sm *SharedSafeMap[T]) Store(key string, value T) {
	i := share(key, len(sm.buckets))
	sm.buckets[i].Set(key, value)
}

// Del deletes the value stored under the given key.
func (sm *SharedSafeMap[T]) Del(key string) {
	i := share(key, len(sm.buckets))
	sm.buckets[i].Del(key)
}

// Keys returns all keys in the map.
func (sm *SharedSafeMap[T]) Keys() []string {
	keys := make([]string, 0)
	for _, bucket := range sm.buckets {
		keys = append(keys, bucket.Keys()...)
	}
	return keys
}

// Range calls the given function for each key-value pair in the map.
func (sm *SharedSafeMap[T]) Range(f func(key string, value T) bool) {
	for _, bucket := range sm.buckets {
		bucket.Range(f)
	}
}

// Len returns the number of elements in the map.
func (sm *SharedSafeMap[T]) Len() int {
	n := 0
	for _, bucket := range sm.buckets {
		n += bucket.Len()
	}
	return n
}
