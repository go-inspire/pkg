/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"hash/fnv"
	"hash/maphash"
	"testing"
)

var buf = make([]byte, 8192+1)

func Benchmark_Hash_fnv(b *testing.B) {
	b.ReportAllocs()
	buf := buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash := fnv.New64a()
		_, _ = hash.Write(buf[:153])
		_ = hash.Sum64()
	}
}

func Benchmark_Hash_maphash(b *testing.B) {
	b.ReportAllocs()
	buf := buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := maphash.Hash{}
		_, _ = h.Write(buf[:153])
		_ = h.Sum64()
	}
}
