/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import "sync/atomic"

// SharedStats 用于记录分片状态
type SharedStats struct {
	Hits map[int]*atomic.Uint64
}

// NewSharedStats 创建一个新的分片状态
func NewSharedStats(n int) *SharedStats {
	hits := make(map[int]*atomic.Uint64, n)
	for i := 0; i < n; i++ {
		hits[i] = &atomic.Uint64{}
	}
	return &SharedStats{Hits: hits}
}

// AddHit 增加分片命中次数
func (s *SharedStats) AddHit(i int) {
	s.Hits[i].Add(1)
}

// GetTotalHits 统计总命中次数
func (s *SharedStats) GetTotalHits() uint64 {
	var total uint64
	for _, hit := range s.Hits {
		total += hit.Load()
	}
	return total
}
