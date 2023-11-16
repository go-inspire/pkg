/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"log/slog"
	"os"
)

var _ Logger = (*slogLogger)(nil)

// slogLogger zap.Logger 的实现
type slogLogger struct {
	log *slog.Logger
}

func newSlogLogger(h slog.Handler) *slogLogger {
	logger := slog.New(h)
	return &slogLogger{
		log: logger,
	}
}

const (
	levelDebug = slog.LevelDebug
	levelInfo  = slog.LevelInfo
	levelWarn  = slog.LevelWarn
	levelError = slog.LevelError
	levelPanic = slog.Level(10)
	levelFatal = slog.Level(12)
)

func customLevel(groups []string, a slog.Attr) slog.Attr {
	// Customize the name of the level key and the output string, including
	// custom level values.
	if a.Key == slog.LevelKey {
		// Handle custom level values.
		level := a.Value.Any().(slog.Level)

		// This could also look up the name from a map or other structure, but
		// this demonstrates using a switch statement to rename levels. For
		// maximum performance, the string values should be constants, but this
		// example uses the raw strings for readability.
		switch {
		case level < levelDebug:
			a.Value = slog.StringValue("TRACE")
		case level < levelInfo:
			a.Value = slog.StringValue("INFO")
		case level < levelWarn:
			a.Value = slog.StringValue("WARN")
		case level < levelError:
			a.Value = slog.StringValue("ERROR")
		case level < levelPanic:
			a.Value = slog.StringValue("ERROR")
		case level < levelFatal:
			a.Value = slog.StringValue("FATAL")
		default:
			a.Value = slog.StringValue("INFO")
		}
	}

	return a
}

func (l *slogLogger) Log(ctx context.Context, level Level, msg string, keyValues ...interface{}) {
	switch level {
	case DebugLevel:
		l.log.Log(ctx, levelDebug, msg, keyValues...)
	case InfoLevel:
		l.log.Log(ctx, levelInfo, msg, keyValues...)
	case WarnLevel:
		l.log.Log(ctx, levelWarn, msg, keyValues...)
	case ErrorLevel:
		l.log.Log(ctx, levelError, msg, keyValues...)
	case PanicLevel:
		l.log.Log(ctx, levelPanic, msg, keyValues...)
	case FatalLevel:
		l.log.Log(ctx, levelFatal, msg, keyValues...)
	}
}

func (l *slogLogger) Close() error {
	return nil
}

func initSlogLogger(lvl Level) Logger {
	logger := newSlogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		ReplaceAttr: customLevel,
		Level:       toSlogLevel(lvl),
	}))

	return logger
}

func toSlogLevel(l Level) slog.Level {
	switch l {
	case DebugLevel:
		return levelDebug
	case InfoLevel:
		return levelInfo
	case WarnLevel:
		return levelWarn
	case ErrorLevel:
		return levelError
	case PanicLevel:
		return levelPanic
	case FatalLevel:
		return levelFatal
	default:
		return levelInfo
	}
}
