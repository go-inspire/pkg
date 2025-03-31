/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"os"
	"sync/atomic"
)

// Option is a function type that configures an Adapter.
// It provides a flexible way to set various options on the Adapter.
type Option func(*Adapter)

// Adapter provides a high-level interface for logging with various levels and formats.
// It wraps a Logger implementation and provides convenient methods for different log levels.
type Adapter struct {
	logger Logger
	lvl    atomic.Int32
	ctx    context.Context
}

// WithLevel returns an Option that sets the minimum enabled log level for the Adapter.
// Only messages at or above this level will be logged.
func WithLevel(level Level) Option {
	return func(la *Adapter) {
		la.SetLevel(level)
	}
}

// WithContext returns an Option that sets the context for the logger adapter.
// This context will be used for all logging operations unless overridden.
func WithContext(ctx context.Context) Option {
	return func(la *Adapter) {
		la.ctx = ctx
	}
}

// NewAdapter creates a new logger adapter with the given logger and options.
// It initializes the adapter with default settings and applies any provided options.
func NewAdapter(logger Logger, opts ...Option) *Adapter {
	la := &Adapter{
		logger: logger,
		ctx:    context.Background(),
	}
	la.SetLevel(InfoLevel)

	if opts != nil {
		for _, o := range opts {
			o(la)
		}
	}
	return la
}

// Level returns the minimum enabled log level for the adapter.
// Only messages at or above this level will be logged.
func (la *Adapter) Level() Level {
	return Level(int(la.lvl.Load()))
}

// SetLevel alters the logging level of the adapter.
// It returns the adapter for method chaining.
func (la *Adapter) SetLevel(l Level) *Adapter {
	la.lvl.Store(int32(l))
	return la
}

// Enabled implements the zapcore.LevelEnabler interface.
// It returns true if the given level is enabled for logging.
func (la *Adapter) Enabled(l Level) bool {
	return la.Level().Enabled(l)
}

// logWithLevel is a generic logging method that handles all log levels.
// It formats the message and calls the underlying logger if the level is enabled.
func (la *Adapter) logWithLevel(level Level, msg string, args ...interface{}) {
	la.output(la.ctx, level, func(ctx context.Context, l Level) {
		la.logger.Log(ctx, l, msg, args...)
	})
}

// logWithFormat is a generic method for formatted logging.
// It formats the message using sprintf before logging.
func (la *Adapter) logWithFormat(level Level, format string, args ...interface{}) {
	la.logWithLevel(level, sprintf(format, args...))
}

// logWithKeyValues is a generic method for key-value pair logging.
// It logs the key-value pairs without a message.
func (la *Adapter) logWithKeyValues(level Level, keyvals ...interface{}) {
	la.logWithLevel(level, "", keyvals...)
}

// Print methods log messages at the Info level.
// Print logs a simple message.
func (la *Adapter) Print(args ...interface{}) {
	la.logWithLevel(InfoLevel, sprint(args...))
}

// Printf logs a formatted message at the Info level.
func (la *Adapter) Printf(msg string, args ...interface{}) {
	la.logWithFormat(InfoLevel, msg, args...)
}

// Printw logs key-value pairs at the Info level.
func (la *Adapter) Printw(keyvals ...interface{}) {
	la.logWithKeyValues(InfoLevel, keyvals...)
}

// Debug methods log messages at the Debug level.
// Debug logs a simple message.
func (la *Adapter) Debug(args ...interface{}) {
	la.logWithLevel(DebugLevel, sprint(args...))
}

// Debugf logs a formatted message at the Debug level.
func (la *Adapter) Debugf(msg string, args ...interface{}) {
	la.logWithFormat(DebugLevel, msg, args...)
}

// Debugw logs key-value pairs at the Debug level.
func (la *Adapter) Debugw(keyvals ...interface{}) {
	la.logWithKeyValues(DebugLevel, keyvals...)
}

// Info methods log messages at the Info level.
// Info logs a simple message.
func (la *Adapter) Info(args ...interface{}) {
	la.logWithLevel(InfoLevel, sprint(args...))
}

// Infof logs a formatted message at the Info level.
func (la *Adapter) Infof(msg string, args ...interface{}) {
	la.logWithFormat(InfoLevel, msg, args...)
}

// Infow logs key-value pairs at the Info level.
func (la *Adapter) Infow(keyvals ...interface{}) {
	la.logWithKeyValues(InfoLevel, keyvals...)
}

// Warn methods log messages at the Warn level.
// Warn logs a simple message.
func (la *Adapter) Warn(args ...interface{}) {
	la.logWithLevel(WarnLevel, sprint(args...))
}

// Warnf logs a formatted message at the Warn level.
func (la *Adapter) Warnf(msg string, args ...interface{}) {
	la.logWithFormat(WarnLevel, msg, args...)
}

// Warnw logs key-value pairs at the Warn level.
func (la *Adapter) Warnw(keyvals ...interface{}) {
	la.logWithKeyValues(WarnLevel, keyvals...)
}

// Error methods log messages at the Error level.
// Error logs a simple message.
func (la *Adapter) Error(args ...interface{}) {
	la.logWithLevel(ErrorLevel, sprint(args...))
}

// Errorf logs a formatted message at the Error level.
func (la *Adapter) Errorf(msg string, args ...interface{}) {
	la.logWithFormat(ErrorLevel, msg, args...)
}

// Errorw logs key-value pairs at the Error level.
func (la *Adapter) Errorw(keyvals ...interface{}) {
	la.logWithKeyValues(ErrorLevel, keyvals...)
}

// Panic methods log messages at the Panic level and then panic.
// Panic logs a simple message and then panics with the same message.
func (la *Adapter) Panic(args ...interface{}) {
	msg := sprint(args...)
	la.logWithLevel(PanicLevel, msg)
	panic(msg)
}

// Panicf logs a formatted message at the Panic level and then panics.
func (la *Adapter) Panicf(msg string, args ...interface{}) {
	formatted := sprintf(msg, args...)
	la.logWithLevel(PanicLevel, formatted)
	panic(formatted)
}

// Panicw logs key-value pairs at the Panic level and then panics.
func (la *Adapter) Panicw(keyvals ...interface{}) {
	msg := sprint(keyvals...)
	la.logWithLevel(PanicLevel, "", keyvals...)
	panic(msg)
}

// Fatal methods log messages at the Fatal level and then exit.
// Fatal logs a simple message and then exits with status code 1.
func (la *Adapter) Fatal(args ...interface{}) {
	msg := sprint(args...)
	la.logWithLevel(FatalLevel, msg)
	la.Flush()
	os.Exit(1)
}

// Fatalf logs a formatted message at the Fatal level and then exits.
func (la *Adapter) Fatalf(msg string, args ...interface{}) {
	formatted := sprintf(msg, args...)
	la.logWithLevel(FatalLevel, formatted)
	la.Flush()
	os.Exit(1)
}

// Fatalw logs key-value pairs at the Fatal level and then exits.
func (la *Adapter) Fatalw(keyvals ...interface{}) {
	la.logWithLevel(FatalLevel, "", keyvals...)
	la.Flush()
	os.Exit(1)
}

// Flush flushes any buffered log entries.
// It calls Close() on the underlying logger.
func (la *Adapter) Flush() error {
	return la.Close()
}

// Close closes the underlying logger.
// It should be called when the logger is no longer needed.
func (la *Adapter) Close() error {
	return la.logger.Close()
}

// Log logs a message at the specified level with the given context.
// If no context is provided, it uses the adapter's default context.
func (la *Adapter) Log(ctx context.Context, level Level, msg string, keyValues ...interface{}) {
	if ctx == nil {
		ctx = la.ctx
	}
	la.logger.Log(ctx, level, msg, keyValues...)
}

// output handles the actual logging output with level checking.
// It only calls the provided log function if the level is enabled.
func (la *Adapter) output(ctx context.Context, l Level, log func(ctx context.Context, level Level)) {
	if la.Enabled(l) {
		log(ctx, l)
	}
}
