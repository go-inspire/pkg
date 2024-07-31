/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"hash/maphash"
	"unsafe"
)

var h maphash.Hash

func hashString(key string) uint64 {
	h.Reset()
	_, _ = h.WriteString(key)
	return h.Sum64()
}

func share(key string, buckets int) int {
	i := hashString(key) % uint64(buckets)
	return int(i)
}

// String2Bytes zero-copy string convert to slice
func String2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
