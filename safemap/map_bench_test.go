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

type SyncMap struct {
	m sync.Map
}

func (o *SyncMap) Get(key string) (struct{}, bool) {
	_, ok := o.m.Load(key)
	return struct{}{}, ok
}
func (o *SyncMap) Set(key string, value struct{}) {
	o.m.Store(key, value)
}

func (o *SyncMap) Del(key string) {
	o.m.Delete(key)
}

func NewSyncMap() *SyncMap {
	return &SyncMap{}
}

type SkipMap struct {
	m *skipmap.StringMap
}

func (o *SkipMap) Get(key string) (struct{}, bool) {
	_, ok := o.m.Load(key)
	return struct{}{}, ok
}
func (o *SkipMap) Set(key string, value struct{}) {
	o.m.Store(key, value)
}

func (o *SkipMap) Del(key string) {
	o.m.Delete(key)
}

func NewSkipMap() *SkipMap {
	return &SkipMap{
		m: skipmap.NewString(),
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

	b.Run("SafeMap", func(b *testing.B) {
		hm := NewSafeMap[struct{}]()
		benchmarkMap(b, hm, writes, reads)
	})

}

func BenchmarkMaps_W100_R100(b *testing.B) {
	benchmarkMaps(b, 100, 100)

}

func BenchmarkMaps_W100_R1000(b *testing.B) {
	benchmarkMaps(b, 1000, 100)
}
