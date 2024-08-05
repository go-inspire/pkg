/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"fmt"
	"github.com/bytedance/gopkg/lang/fastrand"
	"math"
	"strconv"
	"testing"
)

const initSize = 1 << 10 // for `load` `1Delete9Store90Load` `1Range9Delete90Store900Load`
const randM = math.MaxUint32

type bench[T any] struct {
	setup func(*testing.B, T)
	perG  func(b *testing.B, pb *testing.PB, i int, m T)
}

func benchMap(b *testing.B, bench bench[mapInterface[struct{}]]) {
	ms := []mapInterface[struct{}]{
		//NewDeepCopyMap[struct{}](),
		NewSyncMap[struct{}](),
		NewSkipMap[struct{}](),
		NewSafeMap[struct{}](),
		NewSyncMapShared[struct{}](),
		NewSharedSafeMap[struct{}](),
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

func BenchmarkMapStore(b *testing.B) {
	benchMap(b, bench[mapInterface[struct{}]]{
		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface[struct{}]) {
			for pb.Next() {
				m.Store(strconv.Itoa(int(fastrand.Uint32())), struct{}{})
			}
		},
	})

}

func benchmarkMapLoadHits(b *testing.B, hits, misses int) {
	benchMap(b, bench[mapInterface[struct{}]]{
		setup: func(_ *testing.B, m mapInterface[struct{}]) {
			for i := 0; i < initSize*(hits+misses); i++ {
				if fastrand.Uint32n(uint32(hits+misses)) < uint32(hits) {
					m.Store(strconv.Itoa(i), struct{}{})
				}
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface[struct{}]) {
			for pb.Next() {
				m.Load(strconv.Itoa(int(fastrand.Uint32n(uint32(initSize * (hits + misses))))))
			}
		},
	})
}

func BenchmarkMapLoad90Hits(b *testing.B) {
	benchmarkMapLoadHits(b, 9, 1)
}

func BenchmarkMapLoad50Hits(b *testing.B) {
	benchmarkMapLoadHits(b, 5, 5)
}

func benchmarkMapStoreLoad(b *testing.B, writes, reads uint32) {
	benchMap(b, bench[mapInterface[struct{}]]{
		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface[struct{}]) {
			for pb.Next() {
				u := fastrand.Uint32n(writes + reads)
				if u < writes {
					m.Store(strconv.Itoa(int(fastrand.Uint32n(randM))), struct{}{})
				} else {
					m.Load(strconv.Itoa(int(fastrand.Uint32n(randM))))
				}
			}
		},
	})
}

func BenchmarkMap50Store50Load(b *testing.B) {
	benchmarkMapStoreLoad(b, 50, 50)
}

func BenchmarkMap30Store70Load(b *testing.B) {
	benchmarkMapStoreLoad(b, 30, 70)
}

func BenchmarkMap10Store90Load(b *testing.B) {
	benchmarkMapStoreLoad(b, 10, 90)
}
