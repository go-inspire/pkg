/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"github.com/go-inspire/pkg/internal/wyhash"
)

func hashString(key string) uint64 {
	return wyhash.Sum64String(key)
}

func share(key string, shares int) int {
	i := hashString(key) % uint64(shares)
	return int(i)
}
