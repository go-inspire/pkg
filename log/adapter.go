/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
)

// Option is Helper option.
type Option func(*Adapter)

type Adapter struct {
	logger Logger
	lvl    atomic.Int32
}

// WithLevel with log level.
func WithLevel(level Level) Option {
	return func(la *Adapter) {
		la.SetLevel(level)
	}
}

// NewAdapter new a logger adapter.
func NewAdapter(logger Logger, opts ...Option) *Adapter {
	la := &Adapter{
		logger: logger,
	}
	la.SetLevel(InfoLevel)

	if opts != nil {
		for _, o := range opts {
			o(la)
		}
	}
	return la
}

// Level returns the minimum enabled log level.
func (la *Adapter) Level() Level {
	return Level(int(la.lvl.Load()))
}

// SetLevel alters the logging level.
func (la *Adapter) SetLevel(l Level) *Adapter {
	la.lvl.Store(int32(l))
	return la
}

// Enabled implements the zapcore.LevelEnabler interface, which allows the
// AtomicLevel to be used in place of traditional static levels.
func (la *Adapter) Enabled(l Level) bool {
	return la.Level().Enabled(l)
}

var ctx = context.Background()

// Print uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Print(args ...interface{}) {
	la.Info(args...)
}

// Printf uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Printf(msg string, args ...interface{}) {
	la.Infof(msg, args...)
}

// Printw uses key,value to construct and logger a message.
func (la *Adapter) Printw(keyvals ...interface{}) {
	la.Infow(keyvals...)
}

// Debug uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Debug(args ...interface{}) {
	la.output(ctx, DebugLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, sprint(args...))
	})

}

// Debugf uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Debugf(msg string, args ...interface{}) {
	la.output(ctx, DebugLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, sprintf(msg, args...))
	})
}

// Debugw uses key,value to construct and logger a message.
func (la *Adapter) Debugw(keyvals ...interface{}) {
	la.output(ctx, DebugLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", keyvals...)
	})

}

// Info uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Info(args ...interface{}) {
	la.output(ctx, InfoLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, sprint(args...))
	})
}

// Infof uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Infof(msg string, args ...interface{}) {
	la.output(ctx, InfoLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, sprintf(msg, args...))
	})
}

// Infow uses key,value to construct and logger a message.
func (la *Adapter) Infow(keyvals ...interface{}) {
	la.output(ctx, InfoLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", keyvals...)
	})

}

// Warn uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Warn(args ...interface{}) {
	la.output(ctx, WarnLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprint(args...))
	})
}

// Warnf uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Warnf(msg string, args ...interface{}) {
	la.output(ctx, WarnLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprintf(msg, args...))
	})
}

// Warnw uses key,value to construct and logger a message.
func (la *Adapter) Warnw(keyvals ...interface{}) {
	la.output(ctx, WarnLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", keyvals...)
	})

}

// Error uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Error(args ...interface{}) {
	la.output(ctx, ErrorLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprint(args...))
	})

}

// Errorf uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Errorf(msg string, args ...interface{}) {
	la.output(ctx, ErrorLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprintf(msg, args...))
	})
}

// Errorw uses key,value to construct and logger a message.
func (la *Adapter) Errorw(keyvals ...interface{}) {
	la.output(ctx, ErrorLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", keyvals...)
	})
}

// Panic uses fmt.Sprint to construct and logger a message, then panics.
func (la *Adapter) Panic(args ...interface{}) {
	la.output(ctx, PanicLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprint(args...))
	})

	panic(fmt.Sprint(args...))
}

// Panicf uses fmt.Sprintf to logger a templated message, then panics.
func (la *Adapter) Panicf(msg string, args ...interface{}) {
	la.output(ctx, PanicLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprintf(msg, args...))
	})

	panic(fmt.Sprintf(msg, args...))
}

// Panicw uses key,value to construct and logger a message.
func (la *Adapter) Panicw(keyvals ...interface{}) {
	la.output(ctx, PanicLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", keyvals...)
	})
	panic(fmt.Sprint(keyvals...))
}

// Fatal uses fmt.Sprint to construct and logger a message, then calls os.Exit.
func (la *Adapter) Fatal(args ...interface{}) {
	la.output(ctx, FatalLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprint(args...))
	})

	os.Exit(1)
}

// Fatalf uses fmt.Sprintf to logger a templated message, then calls os.Exit.
func (la *Adapter) Fatalf(msg string, args ...interface{}) {
	la.output(ctx, FatalLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", sprintf(msg, args...))
	})
	os.Exit(1)
}

// Fatalw uses key,value to construct and logger a message.
func (la *Adapter) Fatalw(keyvals ...interface{}) {
	la.output(ctx, FatalLevel, func(ctx context.Context, level Level) {
		la.logger.Log(ctx, level, "", keyvals...)
	})
	os.Exit(1)

}

// Flush flushing any buffered log entries. Applications should take care to call Sync before exiting.
func (la *Adapter) Flush() error {
	return la.Close()
}

func (la *Adapter) Close() error {
	return la.logger.Close()
}

func (la *Adapter) Log(_ context.Context, level Level, msg string, keyValues ...interface{}) {
	la.logger.Log(ctx, level, msg, keyValues...)
}

func (la *Adapter) output(ctx context.Context, l Level, log func(ctx context.Context, level Level)) {
	if la.Enabled(l) {
		log(ctx, l)
	}

}
