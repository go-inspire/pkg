/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"runtime"
	"sync"
	"testing"
)

type OriginWithRWLockShared struct {
	shares []OriginWithRWLock
}

func (o *OriginWithRWLockShared) Get(key int) (struct{}, bool) {
	i := share(key, len(o.shares))
	return o.shares[i].Get(key)
}
func (o *OriginWithRWLockShared) Set(key int, value struct{}) {
	i := share(key, len(o.shares))
	o.shares[i].Set(key, value)
}

func (o *OriginWithRWLockShared) Del(key int) {
	i := share(key, len(o.shares))
	o.shares[i].Del(key)
}

func NewOriginWithRWLockShared() *OriginWithRWLockShared {
	num := runtime.GOMAXPROCS(0)
	shares := make([]OriginWithRWLock, 0, num)
	for i := 0; i < num; i++ {
		shares = append(shares, OriginWithRWLock{
			m: make(map[int]struct{}),
			l: sync.RWMutex{},
		})
	}
	return &OriginWithRWLockShared{
		shares: shares,
	}
}

func share(key int, buckets int) int {
	return key % buckets
}

type SyncMapShared struct {
	shares []SyncMap
}

func NewSyncMapShared() *SyncMapShared {
	return &SyncMapShared{
		shares: make([]SyncMap, runtime.GOMAXPROCS(0)),
	}
}

func (o *SyncMapShared) Get(key int) (struct{}, bool) {
	i := share(key, len(o.shares))
	return o.shares[i].Get(key)
}
func (o *SyncMapShared) Set(key int, value struct{}) {
	i := share(key, len(o.shares))
	o.shares[i].Set(key, value)
}

func (o *SyncMapShared) Del(key int) {
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

	b.Run("OriginWithRWLock", func(b *testing.B) {
		hm := NewOriginWithRWLock()
		benchmarkMap(b, hm, writes, reads)
	})

	b.Run("OriginWithRWLockShare", func(b *testing.B) {
		hm := NewOriginWithRWLockShared()
		benchmarkMap(b, hm, writes, reads)
	})

}

func BenchmarkShareMaps_W100_R100(b *testing.B) {
	benchmarkShareMaps(b, 100, 100)

}

func BenchmarkShareMaps_W100_R1000(b *testing.B) {
	benchmarkShareMaps(b, 1000, 100)
}
