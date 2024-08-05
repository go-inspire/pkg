/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"fmt"
	"github.com/bytedance/gopkg/collection/skipset"
	"github.com/bytedance/gopkg/lang/fastrand"
	"slices"
	"strconv"
	"sync"
	"testing"
)

type Set interface {
	Add(key string) bool
	Contains(key string) bool
	Remove(key string) bool
}

func benchSet(b *testing.B, bench bench[Set]) {
	ms := []Set{
		//&Array{
		//	dirty: make([]string, 0, initSize),
		//},
		&SyncSet[string]{},
		NewSafeHashSet[string](),
		skipset.NewString(),
	}
	for _, m := range ms {
		b.Run(fmt.Sprintf("%T", m), func(b *testing.B) {
			if bench.setup != nil {
				bench.setup(b, m)
			}

			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				bench.perG(b, pb, b.N, m)
			})
		})
	}
}

type Array struct {
	m []string
}

func (a *Array) Add(key string) bool {
	a.m = append(a.m, key)
	return true
}

func (a *Array) Contains(key string) bool {
	return slices.Contains(a.m, key)
}

func (a *Array) Remove(key string) bool {
	for i, v := range a.m {
		if v == key {
			a.m = append(a.m[:i], a.m[i+1:]...)
			return true
		}
	}
	return false
}

// SyncMap
type SyncSet[T comparable] struct {
	m sync.Map
}

func (a *SyncSet[T]) Add(key string) bool {
	_, loaded := a.m.LoadOrStore(key, struct{}{})
	return !loaded
}

func (a *SyncSet[T]) Contains(key string) bool {
	_, ok := a.m.Load(key)
	return ok
}

func (a *SyncSet[T]) Remove(key string) bool {
	a.m.Delete(key)
	return true
}

func benchmarkSet(b *testing.B, adds, contains uint32) {
	benchSet(b, bench[Set]{
		perG: func(b *testing.B, pb *testing.PB, i int, m Set) {
			for pb.Next() {
				u := fastrand.Uint32n(adds + contains)
				if u < adds {
					m.Add(strconv.Itoa(int(fastrand.Uint32n(randM))))
				} else {
					m.Contains(strconv.Itoa(int(fastrand.Uint32n(randM))))
				}
			}
		},
	})
}

func BenchmarkSet10Add90Contains(b *testing.B) {
	benchmarkSet(b, 10, 90)
}

func BenchmarkSet30Add70Contains(b *testing.B) {
	benchmarkSet(b, 30, 70)
}

func BenchmarkSet50Add50Contains(b *testing.B) {
	benchmarkSet(b, 50, 50)
}

func benchmarkSetHits(b *testing.B, hits, misses int) {
	benchSet(b, bench[Set]{
		setup: func(_ *testing.B, m Set) {
			for i := 0; i < initSize*(hits+misses); i++ {
				if fastrand.Uint32n(uint32(hits+misses)) < uint32(hits) {
					m.Add(strconv.Itoa(i))
				}
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m Set) {
			for pb.Next() {
				m.Contains(strconv.Itoa(int(fastrand.Uint32n(uint32(initSize * (hits + misses))))))
			}
		},
	})
}

func BenchmarkSetContains90Hits(b *testing.B) {
	benchmarkSetHits(b, 9, 1)
}

func BenchmarkSetContains50Hits(b *testing.B) {
	benchmarkSetHits(b, 5, 5)
}

func benchmarkSetAddHits(b *testing.B, hits, misses int) {
	benchSet(b, bench[Set]{
		setup: func(_ *testing.B, m Set) {
			for i := 0; i < initSize*(hits+misses); i++ {
				if fastrand.Uint32n(uint32(hits+misses)) < uint32(hits) {
					m.Add(strconv.Itoa(i))
				}
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m Set) {
			for pb.Next() {
				n := int(fastrand.Uint32n(uint32(initSize * (hits + misses))))
				if n < hits {
					m.Add(strconv.Itoa(n))
				} else {
					m.Contains(strconv.Itoa(n))
				}
			}
		},
	})
}

func BenchmarkSetAddContains90Hits(b *testing.B) {
	benchmarkSetAddHits(b, 9, 1)
}

func BenchmarkSetAddContains50Hits(b *testing.B) {
	benchmarkSetAddHits(b, 5, 5)
}
