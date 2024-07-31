/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"github.com/bytedance/gopkg/lang/fastrand"
	"math"
	"strconv"
	"testing"
)

const initSize = 1 << 10 // for `load` `1Delete9Store90Load` `1Range9Delete90Store900Load`
const randM = math.MaxInt

type Map interface {
	Set(key string, val struct{})
	Get(key string) (struct{}, bool)
	Del(key string)
}

// benchmarkMap 大量并发读写的场景
func benchmarkMap(b *testing.B, hm Map, reads, writes uint32) {
	for i := 0; i < initSize; i++ {
		hm.Set(strconv.Itoa(fastrand.Intn(randM)), struct{}{})
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrand.Uint32n(reads + writes)
			if u < writes {
				hm.Set(strconv.Itoa(fastrand.Intn(randM)), struct{}{})
			} else {
				_, _ = hm.Get(strconv.Itoa(fastrand.Intn(randM)))
			}
		}
	})
}
