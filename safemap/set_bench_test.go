/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"github.com/bytedance/gopkg/lang/fastrand"
	"slices"
	"strconv"
	"testing"
)

type Set interface {
	Add(key string) bool
	Contains(key string) bool
	Remove(key string) bool
}

// benchmarkSet 大量并发读写的场景
func benchmarkSet(b *testing.B, s Set, reads, writes uint32) {
	for i := 0; i < initSize; i++ {
		s.Add(strconv.Itoa(fastrand.Intn(randM)))
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrand.Uint32n(reads + writes)
			if u < writes {
				_ = s.Add(strconv.Itoa(fastrand.Intn(randM)))
			} else {
				_ = s.Contains(strconv.Itoa(fastrand.Intn(randM)))
			}
		}
	})
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

func benchmarkSets(b *testing.B, reads, writes uint32) {
	b.Logf("Writer: %d,Reader: %d", writes, reads)

	// read write goroutine
	b.Run("Array", func(b *testing.B) {
		set := &Array{
			m: make([]string, 0, initSize),
		}
		benchmarkSet(b, set, writes, reads)
	})

	//b.Run("HashSet", func(b *testing.B) {
	//	set := NewHashSetWithSize[string](initSize)
	//	benchmarkSet(b, set, writes, reads)
	//})

	b.Run("SafeHashSet", func(b *testing.B) {
		set := NewSafeHashSetWithSize[string](initSize)
		benchmarkSet(b, set, writes, reads)
	})

}

func BenchmarkSets_W100_R100(b *testing.B) {
	benchmarkSets(b, 100, 100)

}

func BenchmarkSets_W100_R1000(b *testing.B) {
	benchmarkSets(b, 1000, 100)
}
