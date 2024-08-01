/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"fmt"
	wyhash0 "github.com/go-inspire/pkg/internal/wyhash"
	"github.com/zhangyunhao116/wyhash"
	"hash/fnv"
	"hash/maphash"
	"runtime"
	"testing"
)

type hasher interface {
	Hash(data []byte, seed uint64) uint64
	HashString(data string, seed uint64) uint64
}

var buf = make([]byte, 8192+1)

func benchmarkHash(h hasher, b *testing.B) {
	sizes := []int{
		0, 1, 3, 4, 8, 9, 16, 17, 32,
		//33, 64, 65, 96, 97, 128, 129, 240, 241,
		//512, 1024, 100 * 1024,
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Hash_%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			var acc uint64
			b.ReportAllocs()
			b.ResetTimer()
			d := buf[:size]
			for i := 0; i < b.N; i++ {
				acc = h.Hash(d, 0)
			}
			runtime.KeepAlive(acc)
		})
	}
}

func benchmarkHashString(h hasher, b *testing.B) {
	sizes := []int{
		0, 1, 3, 4, 8, 9, 16, 17, 32,
		//33, 64, 65, 96, 97, 128, 129, 240, 241,
		//512, 1024, 100 * 1024,
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("HashString_%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			var acc uint64

			b.ReportAllocs()
			b.ResetTimer()
			d := string(make([]byte, size))
			for i := 0; i < b.N; i++ {
				acc = h.HashString(d, 0)
			}
			runtime.KeepAlive(acc)
		})
	}
}

type wyhasher struct{}

func (wyhasher) Hash(data []byte, seed uint64) uint64 {
	return wyhash0.Sum64WithSeed(data, seed)
}
func (wyhasher) HashString(data string, seed uint64) uint64 {
	return wyhash0.Sum64StringWithSeed(data, seed)
}

type wyhasher1 struct{}

func (wyhasher1) Hash(data []byte, seed uint64) uint64 {
	return wyhash.Sum64WithSeed(data, seed)
}
func (wyhasher1) HashString(data string, seed uint64) uint64 {
	return wyhash.Sum64StringWithSeed(data, seed)
}

type wyhasher3 struct{}

func (wyhasher3) Hash(data []byte, seed uint64) uint64 {
	return wyhash.Sum64WithSeedV3(data, seed)
}
func (wyhasher3) HashString(data string, seed uint64) uint64 {
	return wyhash.Sum64StringWithSeedV3(data, seed)

}

type Fnv64a struct{}

func (Fnv64a) Hash(data []byte, seed uint64) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write(data)
	return hash.Sum64()
}

func (Fnv64a) HashString(data string, seed uint64) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(data))
	return hash.Sum64()
}

type Maphash struct{}

var mhseed = maphash.MakeSeed()

var h maphash.Hash

func (Maphash) Hash(data []byte, seed uint64) uint64 {
	return maphash.Bytes(mhseed, data)
}

func (m Maphash) HashString(data string, seed uint64) uint64 {
	return maphash.String(mhseed, data)

}

func BenchmarkHash(b *testing.B) {
	b.Run("wyhash0", func(b *testing.B) {
		benchmarkHash(wyhasher{}, b)
		benchmarkHashString(wyhasher{}, b)
	})

	b.Run("wyhash1", func(b *testing.B) {
		benchmarkHash(wyhasher1{}, b)
		benchmarkHashString(wyhasher1{}, b)
	})

	b.Run("wyhash3", func(b *testing.B) {
		benchmarkHash(wyhasher3{}, b)
		benchmarkHashString(wyhasher3{}, b)
	})

	//b.Run("Fnv64a", func(b *testing.B) {
	//	benchmarkHash(Fnv64a{}, b)
	//	benchmarkHashString(Fnv64a{}, b)
	//})
	//
	b.Run("Maphash", func(b *testing.B) {
		benchmarkHash(Maphash{}, b)
		benchmarkHashString(Maphash{}, b)
	})
}
