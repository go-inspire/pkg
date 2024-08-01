/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"github.com/zhangyunhao116/wyhash"
	"unsafe"
)

func hashString(key string) uint64 {
	return wyhash.Sum64String(key)
}

func share(key string, shares int) int {
	i := hashString(key) % uint64(shares)
	return int(i)
}

// String2Bytes zero-copy string convert to slice
func String2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
