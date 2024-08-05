// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safemap

import (
	"github.com/bytedance/gopkg/collection/skipmap"
	"runtime"
	"sync"
)

// This file contains reference map implementations for unit-tests.

// mapInterface is the interface Map implements.
type mapInterface[T any] interface {
	Load(key string) (T, bool)
	Store(key string, value T)
	LoadOrStore(key string, value T) (actual T, loaded bool)
	Delete(string) bool
	Range(func(key string, value T) (shouldContinue bool))
}

// SyncMap
type SyncMap[T any] struct {
	m sync.Map
}

func NewSyncMap[T any]() *SyncMap[T] {
	return &SyncMap[T]{}
}

func (m *SyncMap[T]) Load(key string) (T, bool) {
	if value, ok := m.m.Load(key); ok {
		return value.(T), ok
	}
	return *new(T), false
}

func (m *SyncMap[T]) Store(key string, value T) {
	m.m.Store(key, value)
}

func (m *SyncMap[T]) Delete(key string) bool {
	m.m.Delete(key)
	return true
}

func (m *SyncMap[T]) LoadOrStore(key string, value T) (T, bool) {
	actual, loaded := m.m.LoadOrStore(key, value)
	if loaded {
		return actual.(T), loaded
	} else {
		return *new(T), false
	}
}

func (m *SyncMap[T]) Range(f func(key string, value T) (shouldContinue bool)) {
	m.m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(T))
	})
}

// SkipMap
type SkipMap[T any] struct {
	m *skipmap.StringMap
}

func NewSkipMap[T any]() *SkipMap[T] {
	return &SkipMap[T]{
		m: skipmap.NewString(),
	}
}

func (m *SkipMap[T]) Load(key string) (T, bool) {
	if val, ok := m.m.Load(key); ok {
		return val.(T), ok
	}
	return *new(T), false
}

func (m *SkipMap[T]) Store(key string, value T) {
	m.m.Store(key, value)
}

func (m *SkipMap[T]) LoadOrStore(key string, value T) (T, bool) {
	actual, load := m.m.LoadOrStore(key, value)
	return actual.(T), load
}

func (m *SkipMap[T]) Range(f func(key string, value T) (shouldContinue bool)) {
	m.m.Range(func(key string, value interface{}) bool {
		return f(key, value.(T))
	})
}

func (m *SkipMap[T]) Delete(key string) bool {
	return m.m.Delete(key)
}

// SharedSyncMap
type SharedSyncMap[T any] struct {
	shares []sync.Map
}

func NewSyncMapShared[T any]() *SharedSyncMap[T] {
	return &SharedSyncMap[T]{
		shares: make([]sync.Map, runtime.GOMAXPROCS(0)),
	}
}

func (o *SharedSyncMap[T]) Load(key string) (T, bool) {
	i := share(key, len(o.shares))
	if val, ok := o.shares[i].Load(key); ok {
		return val.(T), ok
	} else {
		return *new(T), false
	}

}

func (o *SharedSyncMap[T]) Store(key string, value T) {
	i := share(key, len(o.shares))
	o.shares[i].Store(key, value)
}

func (o *SharedSyncMap[T]) LoadOrStore(key string, value T) (T, bool) {
	i := share(key, len(o.shares))
	actual, loaded := o.shares[i].LoadOrStore(key, value)
	if loaded {
		return actual.(T), loaded
	} else {
		return *new(T), false
	}
}

func (o *SharedSyncMap[T]) Delete(key string) bool {
	i := share(key, len(o.shares))
	o.shares[i].Delete(key)
	return true
}

func (o *SharedSyncMap[T]) Range(f func(key string, value T) bool) {
	for i := range o.shares {
		o.shares[i].Range(func(key, value interface{}) bool {
			return f(key.(string), value.(T))
		})
	}
}
