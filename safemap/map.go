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
	dirty map[string]T
	rw    sync.RWMutex
}

// NewSafeMap 创建一个新的 SafeMap
func NewSafeMap[T any]() *SafeMap[T] {
	return &SafeMap[T]{
		dirty: make(map[string]T),
		rw:    sync.RWMutex{},
	}
}

// Load 返回给定键的值
func (m *SafeMap[T]) Load(key string) (T, bool) {
	m.rw.RLock()
	value, ok := m.dirty[key]
	m.rw.RUnlock()
	return value, ok
}

// Store 设置给定键的值
func (m *SafeMap[T]) Store(key string, value T) {
	m.rw.Lock()
	if m.dirty == nil {
		m.dirty = make(map[string]T)
	}
	m.dirty[key] = value
	m.rw.Unlock()
}

// LoadOrStore 返回给定键的值, 如果不存在则存储给定值
func (m *SafeMap[T]) LoadOrStore(key string, value T) (actual T, loaded bool) {
	m.rw.Lock()
	actual, loaded = m.dirty[key]
	if !loaded {
		actual = value
		if m.dirty == nil {
			m.dirty = make(map[string]T)
		}
		m.dirty[key] = value
	}
	m.rw.Unlock()
	return actual, loaded
}

// Delete 删除给定键的值
func (m *SafeMap[T]) Delete(key string) bool {
	m.rw.Lock()
	delete(m.dirty, key)
	m.rw.Unlock()
	return true
}

// LoadAndDelete 返回给定键的值, 并删除该键
func (m *SafeMap[T]) LoadAndDelete(key string) (T, bool) {
	m.rw.Lock()
	value, loaded := m.dirty[key]
	if loaded {
		delete(m.dirty, key)
	}
	m.rw.Unlock()
	return value, loaded
}

// Keys 返回所有键
func (m *SafeMap[T]) Keys() []string {
	m.rw.RLock()
	keys := make([]string, 0, len(m.dirty))
	for key := range m.dirty {
		keys = append(keys, key)
	}
	m.rw.RUnlock()
	return keys
}

// Range 对 Map 中的每个键值对调用给定的函数
func (m *SafeMap[T]) Range(f func(key string, value T) bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	for key, value := range m.dirty {
		if !f(key, value) {
			return
		}
	}
}

// Len 返回 Map 中的元素数量
func (m *SafeMap[T]) Len() int {
	m.rw.RLock()
	defer m.rw.RUnlock()

	return len(m.dirty)

}

// SharedSafeMap 是一个可以安全地由多个 goroutine 共享的 map. 使用分片思路来实现.
type SharedSafeMap[T any] struct {
	buckets []*SafeMap[T]
}

// NewSharedSafeMap 创建一个新的 SharedSafeMap. 根据当前系统的 CPU 核心数创建对应数量的分片.
func NewSharedSafeMap[T any]() *SharedSafeMap[T] {
	n := runtime.GOMAXPROCS(0)
	buckets := make([]*SafeMap[T], n)
	for i := range buckets {
		buckets[i] = NewSafeMap[T]()
	}
	return &SharedSafeMap[T]{buckets: buckets}

}

// Load 返回给定键的值
func (sm *SharedSafeMap[T]) Load(key string) (T, bool) {
	i := share(key, len(sm.buckets))
	return sm.buckets[i].Load(key)
}

// Store 设置给定键的值
func (sm *SharedSafeMap[T]) Store(key string, value T) {
	i := share(key, len(sm.buckets))
	sm.buckets[i].Store(key, value)
}

// LoadOrStore 返回给定键的值, 如果不存在则存储给定值
func (sm *SharedSafeMap[T]) LoadOrStore(key string, value T) (actual T, loaded bool) {
	i := share(key, len(sm.buckets))
	return sm.buckets[i].LoadOrStore(key, value)
}

// Delete 删除给定键的值
func (sm *SharedSafeMap[T]) Delete(key string) bool {
	i := share(key, len(sm.buckets))
	return sm.buckets[i].Delete(key)
}

// LoadAndDelete 返回给定键的值, 并删除该键
func (sm *SharedSafeMap[T]) LoadAndDelete(key string) (T, bool) {
	i := share(key, len(sm.buckets))
	return sm.buckets[i].LoadAndDelete(key)
}

// Keys 返回所有键
func (sm *SharedSafeMap[T]) Keys() []string {
	keys := make([]string, 0)
	for _, bucket := range sm.buckets {
		keys = append(keys, bucket.Keys()...)
	}
	return keys
}

// Range 对 Map 中的每个键值对调用给定的函数
func (sm *SharedSafeMap[T]) Range(f func(key string, value T) bool) {
	for _, bucket := range sm.buckets {
		bucket.Range(f)
	}
}

// Len 返回 Map 中的元素数量
func (sm *SharedSafeMap[T]) Len() int {
	n := 0
	for _, bucket := range sm.buckets {
		n += bucket.Len()
	}
	return n
}
