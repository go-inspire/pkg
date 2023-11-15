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
	"slices"
	"testing"
)

func Test_Slog(t *testing.T) {
	slog.Info("debug")
	//slog.Debug("debug", "k1")
	slog.Error("debug", "k1", "v1")
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	l.Info("debug", "k1", "v1")

}

func TestSlogLogger(t *testing.T) {
	logger := newSlogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		ReplaceAttr: customLevel,
	}))

	ctx := context.Background()

	logger.Log(ctx, DebugLevel, "debug message", "key1", "value1")
	logger.Log(ctx, InfoLevel, "info message", "key2", "value2")
	logger.Log(ctx, WarnLevel, "warn message", "key3", "value3")
	logger.Log(ctx, ErrorLevel, "error message", "key4", "value4")
	logger.Log(ctx, PanicLevel, "panic message", "key5", "value5")
	logger.Log(ctx, FatalLevel, "fatal message", "key6", "value6")
}

func TestSlogLogger_Close(t *testing.T) {
	logger := newSlogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		ReplaceAttr: customLevel,
	}))

	err := logger.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestToSlogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		want  slog.Level
	}{
		{"DebugLevel", DebugLevel, levelDebug},
		{"InfoLevel", InfoLevel, levelInfo},
		{"WarnLevel", WarnLevel, levelWarn},
		{"ErrorLevel", ErrorLevel, levelError},
		{"PanicLevel", PanicLevel, levelPanic},
		{"FatalLevel", FatalLevel, levelFatal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toSlogLevel(tt.level); got != tt.want {
				t.Errorf("toSlogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitSlogLogger(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		want  slog.Level
	}{
		{"DebugLevel", DebugLevel, levelDebug},
		{"InfoLevel", InfoLevel, levelInfo},
		{"WarnLevel", WarnLevel, levelWarn},
		{"ErrorLevel", ErrorLevel, levelError},
		{"PanicLevel", PanicLevel, levelPanic},
		{"FatalLevel", FatalLevel, levelFatal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := initSlogLogger(tt.level)
			slogLogger, ok := logger.(*slogLogger)
			if !ok {
				t.Errorf("initSlogLogger() did not return a *slogLogger")
				return
			}
			handler, ok := slogLogger.log.Handler().(*slog.TextHandler)
			if !ok {
				t.Errorf("initSlogLogger() did not use a *slog.TextHandler")
				return
			}
			if !handler.Enabled(nil, tt.want) {
				t.Errorf("initSlogLogger() level > want %v", tt.want)
			}
		})
	}
}

func TestSlogLogger_Log(t *testing.T) {
	h := &mockHandler{}
	logger := newSlogLogger(h)

	ctx := context.Background()
	msg := "test message"
	keyValues := []interface{}{"key1", "value1"}

	tests := []struct {
		name  string
		level Level
	}{
		{"DebugLevel", DebugLevel},
		{"InfoLevel", InfoLevel},
		{"WarnLevel", WarnLevel},
		{"ErrorLevel", ErrorLevel},
		{"PanicLevel", PanicLevel},
		{"FatalLevel", FatalLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Log(ctx, tt.level, msg, keyValues...)
			if h.r.Message != msg {
				t.Errorf("Log() Message = %v, want %v", h.r.Message, msg)
			}
			attrs := make([]interface{}, 0, len(keyValues))
			h.r.Attrs(func(a slog.Attr) bool {
				attrs = append(attrs, a.Key, a.Value.String())
				return true
			})
			if !slices.Equal(attrs, keyValues) {
				t.Errorf("Log() Attrs = %v, want %v", attrs, keyValues)
			}

			if h.r.Level != toSlogLevel(tt.level) {
				t.Errorf("Log() level = %v, want %v", h.r.Level, tt.level)
			}

		})
	}
}

type mockHandler struct {
	r slog.Record
}

func (h *mockHandler) Enabled(ctx2 context.Context, level slog.Level) bool {
	return true
}

func (h *mockHandler) Handle(ctx2 context.Context, record slog.Record) error {
	h.r = record
	return nil
}

func (h *mockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *mockHandler) WithGroup(name string) slog.Handler {
	return h
}
