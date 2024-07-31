/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"runtime"
	"testing"
)

// Load returns the value stored under the given key.
func (sm *SharedSafeMap[T]) Get(key string) (T, bool) {
	i := share(key, len(sm.buckets))
	return sm.buckets[i].Get(key)
}

// Store stores the given value under the given key.
func (sm *SharedSafeMap[T]) Set(key string, value T) {
	i := share(key, len(sm.buckets))
	sm.buckets[i].Set(key, value)
}

type SyncMapShared struct {
	shares []SyncMap
}

func NewSyncMapShared() *SyncMapShared {
	return &SyncMapShared{
		shares: make([]SyncMap, runtime.GOMAXPROCS(0)),
	}
}

func (o *SyncMapShared) Get(key string) (struct{}, bool) {
	i := share(key, len(o.shares))
	return o.shares[i].Get(key)
}
func (o *SyncMapShared) Set(key string, value struct{}) {
	i := share(key, len(o.shares))
	o.shares[i].Set(key, value)
}

func (o *SyncMapShared) Del(key string) {
	i := share(key, len(o.shares))
	o.shares[i].Del(key)
}

func benchmarkShareMaps(b *testing.B, reads, writes uint32) {
	b.Logf("Writer: %d,Reader: %d", writes, reads)

	// read write goroutine
	b.Run("SyncMap", func(b *testing.B) {
		hm := NewSyncMap()
		benchmarkMap(b, hm, writes, reads)
	})

	b.Run("SyncMapShared", func(b *testing.B) {
		hm := NewSyncMapShared()
		benchmarkMap(b, hm, writes, reads)
	})

	b.Run("SafeMap", func(b *testing.B) {
		hm := NewSafeMap[struct{}]()
		benchmarkMap(b, hm, writes, reads)
	})

	b.Run("SharedSafeMap", func(b *testing.B) {
		hm := NewSharedSafeMap[struct{}]()
		benchmarkMap(b, hm, writes, reads)
	})

}

func BenchmarkShareMaps_W100_R100(b *testing.B) {
	benchmarkShareMaps(b, 100, 100)

}

func BenchmarkShareMaps_W100_R1000(b *testing.B) {
	benchmarkShareMaps(b, 1000, 100)
}
