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

// Get 返回给定键的值
func (sm *SafeMap[T]) Get(key string) (T, bool) {
	sm.l.RLock()
	value, ok := sm.m[key]
	sm.l.RUnlock()
	return value, ok
}

// Set 设置给定键的值
func (sm *SafeMap[T]) Set(key string, value T) {
	sm.l.Lock()
	sm.m[key] = value
	sm.l.Unlock()
}

// Del 删除给定键的值
func (sm *SafeMap[T]) Del(key string) {
	sm.l.Lock()
	delete(sm.m, key)
	sm.l.Unlock()
}

// Keys 返回所有键
func (sm *SafeMap[T]) Keys() []string {
	sm.l.RLock()
	keys := make([]string, 0, len(sm.m))
	for key := range sm.m {
		keys = append(keys, key)
	}
	sm.l.RUnlock()
	return keys
}

// Range 对 Map 中的每个键值对调用给定的函数
func (sm *SafeMap[T]) Range(f func(key string, value T) bool) {
	sm.l.RLock()
	defer sm.l.RUnlock()

	for key, value := range sm.m {
		if !f(key, value) {
			return
		}
	}
}

// Len 返回 Map 中的元素数量
func (sm *SafeMap[T]) Len() int {
	sm.l.RLock()
	defer sm.l.RUnlock()

	return len(sm.m)

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
	return sm.buckets[i].Get(key)
}

// Store 设置给定键的值
func (sm *SharedSafeMap[T]) Store(key string, value T) {
	i := share(key, len(sm.buckets))
	sm.buckets[i].Set(key, value)
}

// Del 删除给定键的值
func (sm *SharedSafeMap[T]) Del(key string) {
	i := share(key, len(sm.buckets))
	sm.buckets[i].Del(key)
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
