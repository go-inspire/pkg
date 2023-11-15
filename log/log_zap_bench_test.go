/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"go.uber.org/zap"
	"io"
	"log"
	"testing"
)

func Benchmark_zapLogger(b *testing.B) {
	logger := zap.NewNop()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("No context.")
		}
	})

}

func Benchmark_zapLoggerSugar(b *testing.B) {
	logger := zap.NewNop().Sugar()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("No context.")
		}
	})

}

func Benchmark_log(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.")
		}
	})
}
