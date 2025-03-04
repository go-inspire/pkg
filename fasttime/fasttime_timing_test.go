/*
 * Copyright 2025 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package fasttime

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkUnixTimestamp(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var ts int64
		for pb.Next() {
			ts += UnixTimestamp()
		}
		atomic.StoreInt64(&Sink, ts)
	})
}

func BenchmarkTimeNowUnix(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var ts int64
		for pb.Next() {
			ts += time.Now().Unix()
		}
		atomic.StoreInt64(&Sink, ts)
	})
}

// Sink should prevent from code elimination by optimizing compiler
var Sink int64

func BenchmarkTimeNow(b *testing.B) {
	b.Run("TimeNow", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			var t time.Time
			for pb.Next() {
				t = time.Now()
			}
			runtime.KeepAlive(t)
		})
	})

	b.Run("FasttimeNow", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			var t time.Time
			for pb.Next() {
				t = Now()
			}
			runtime.KeepAlive(t)
		})
	})
}
