/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"testing"
)

func getLogger() *Adapter {
	//a := NewAdapter(NopLogger, WithLevel(WarnLevel))
	a := NewAdapter(NopLogger, WithLevel(InfoLevel))
	return a
}

func Benchmark_Logger_Print(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.")
		}
	})
}

func Benchmark_Logger_Print_With3StringgArgs(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.", "1", "2")
		}
	})
}

func Benchmark_Logger_Print_With5StringArgs(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.", "1", "2", "1", "2")
		}
	})
}

func Benchmark_Logger_Print_With3StringArgsAnd2IntArgs(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.", "1", "2", 1, 2)
		}
	})
}

func Benchmark_Logger_Printf(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.")
		}
	})
}

func Benchmark_Logger_Printf_With2StringArgs(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.%s%s", "1", "2")
		}
	})
}

func Benchmark_Logger_Printf_With4StringArgs(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.%s%s%s%s", "1", "2", "1", "2")
		}
	})
}

func Benchmark_Logger_Printf_With2StringArgsAnd2IntArgs(b *testing.B) {
	logger := getLogger()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.%s%s%d%d", "1", "2", 1, 2)
		}
	})
}

var NopWriter = discard{}

type discard struct{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

func Benchmark_SysLog_Print(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.")
		}
	})
}

func Benchmark_SysLog_Print_With3StringArgs(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.", "1", "2")
		}
	})
}

func Benchmark_SysLog_Print_With3StringArgsAnd2IntArgs(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.", "1", "2", 1, 2)
		}
	})
}

func Benchmark_SysLog_Printf(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Printf("No context.")
		}
	})
}

func Benchmark_SysLog_Printf_With2StringArgs(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Printf("No context.%s%s", "1", "2")
		}
	})
}

func Benchmark_SysLog_Printf_With2StringArgsAnd2IntArgs(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Printf("No context.%s%s%d%d", "1", "2", 1, 2)
		}
	})
}

func Benchmark_fmt_Print(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Fprint(io.Discard, "No context.")
		}
	})

}

func Benchmark_fmt_Printf(b *testing.B) {
	log.SetOutput(NopWriter)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Fprintf(io.Discard, "No context.")
		}
	})
}

func Benchmark_slog(b *testing.B) {
	log := slog.New(nopHandler{})
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("No context.")
		}
	})
}

func Benchmark_slog_With3StringArgs(b *testing.B) {
	log := slog.New(nopHandler{})
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("No context.", "1", "2")
		}
	})
}

func Benchmark_slog_With3StringArgsAnd2IntArgs(b *testing.B) {
	log := slog.New(nopHandler{})
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("No context.", "1", "2", "1", 2)
		}
	})
}

type nopHandler struct{}

func (nopHandler) Enabled(context.Context, slog.Level) bool  { return true }
func (nopHandler) Handle(context.Context, slog.Record) error { return nil }
func (h nopHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}
func (h nopHandler) WithGroup(string) slog.Handler {
	return h
}
