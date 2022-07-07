/*
 * Copyright 2022 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"fmt"
	"testing"
)

func Benchmark_fmtSprint(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Sprint("No context.")
		}
	})
}

func Benchmark_fmtSprint_WithArgs(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Sprint("No context.", "arg1", 1)
		}
	})
}

func Benchmark_format(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sprint("No context.")
		}
	})
}

func Benchmark_format_WithArgs(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sprint("No context.", "arg1", 1)
		}
	})
}
