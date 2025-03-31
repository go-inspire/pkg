/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level is a logging priority. Higher levels are more important.
type Level = zapcore.Level

// AtomicLevel is an atomically changeable, dynamic logging level. It lets you
type AtomicLevel = zap.AtomicLevel

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

// Logger is a logger interface.
type Logger interface {
	io.Closer

	// Log logs a message at a given level.
	Log(ctx context.Context, level Level, msg string, keyValues ...interface{})
}

// NopLogger is a logger that does nothing.
var NopLogger Logger = nopLogger{}

type nopLogger struct{}

// Log is a no-op implementation of the Logger interface.
func (nopLogger) Log(ctx context.Context, level Level, msg string, keyValues ...interface{}) {}

// Close is a no-op implementation of the io.Closer interface.
func (nopLogger) Close() error {
	return nil
}
