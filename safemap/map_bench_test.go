/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"github.com/bytedance/gopkg/collection/skipmap"
	"sync"
	"testing"
)

type OriginWithRWLock struct {
	m map[int]struct{}
	l sync.RWMutex
}

func (o *OriginWithRWLock) Get(key int) (struct{}, bool) {
	o.l.RLock()
	v, ok := o.m[key]
	o.l.RUnlock()
	return v, ok
}
func (o *OriginWithRWLock) Set(key int, value struct{}) {
	o.l.Lock()
	o.m[key] = value
	o.l.Unlock()
}

func (o *OriginWithRWLock) Del(key int) {
	o.l.Lock()
	delete(o.m, key)
	o.l.Unlock()
}

func NewOriginWithRWLock() *OriginWithRWLock {
	return &OriginWithRWLock{
		m: make(map[int]struct{}),
		l: sync.RWMutex{},
	}
}

type SyncMap struct {
	m sync.Map
}

func (o *SyncMap) Get(key int) (struct{}, bool) {
	_, ok := o.m.Load(key)
	return struct{}{}, ok
}
func (o *SyncMap) Set(key int, value struct{}) {
	o.m.Store(key, value)
}

func (o *SyncMap) Del(key int) {
	o.m.Delete(key)
}

func NewSyncMap() *SyncMap {
	return &SyncMap{}
}

type SkipMap struct {
	m *skipmap.IntMap
}

func (o *SkipMap) Get(key int) (struct{}, bool) {
	_, ok := o.m.Load(key)
	return struct{}{}, ok
}
func (o *SkipMap) Set(key int, value struct{}) {
	o.m.Store(key, value)
}

func (o *SkipMap) Del(key int) {
	o.m.Delete(key)
}

func NewSkipMap() *SkipMap {
	return &SkipMap{
		m: skipmap.NewInt(),
	}
}

func benchmarkMaps(b *testing.B, reads, writes uint32) {
	b.Logf("Writer: %d,Reader: %d", writes, reads)

	// read write goroutine
	b.Run("SyncMap", func(b *testing.B) {
		hm := NewSyncMap()
		benchmarkMap(b, hm, writes, reads)
	})

	b.Run("SkipMap", func(b *testing.B) {
		hm := NewSkipMap()
		benchmarkMap(b, hm, writes, reads)
	})

	b.Run("OriginWithRWLock", func(b *testing.B) {
		hm := NewOriginWithRWLock()
		benchmarkMap(b, hm, writes, reads)
	})

}

func BenchmarkMaps_W100_R100(b *testing.B) {
	benchmarkMaps(b, 100, 100)

}

func BenchmarkMaps_W100_R1000(b *testing.B) {
	benchmarkMaps(b, 1000, 100)
}
