/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"testing"
)

// A Syncer is a spy for the Sync portion of zapcore.WriteSyncer.
type Syncer struct {
	err    error
	called bool
}

// SetError sets the error that the Sync method will return.
func (s *Syncer) SetError(err error) {
	s.err = err
}

// Sync records that it was called, then returns the user-supplied error (if
// any).
func (s *Syncer) Sync() error {
	s.called = true
	return s.err
}

// Called reports whether the Sync method was called.
func (s *Syncer) Called() bool {
	return s.called
}

// A Discarder sends all writes to ioutil.Discard.
type Discarder struct{ Syncer }

// Write implements io.Writer.
func (d *Discarder) Write(b []byte) (int, error) {
	return io.Discard.Write(b)
}

func Benchmark_zapLogger(b *testing.B) {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
			//zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
			&Discarder{},
			zap.DebugLevel,
		))
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("No context.")
		}
	})

}

func Benchmark_zapLoggerSugar(b *testing.B) {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
			//zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
			&Discarder{},
			zap.DebugLevel,
		)).Sugar()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("No context.")
		}
	})

}

func BenchmarkSysLog(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.")
		}
	})
}

func Benchmark_fmtPrint(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Fprintf(io.Discard, "No context.")
		}
	})
}
