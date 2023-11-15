/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"strings"
	"testing"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name string
		cfg  ZapConfig
		want *zapLogger
	}{
		{
			name: "ValidConfig",
			cfg: ZapConfig{
				Config: zap.Config{
					Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
					Development: true,
					Encoding:    "console",
					EncoderConfig: zapcore.EncoderConfig{
						MessageKey: "message",
					},
				},
			},
			want: &zapLogger{},
		},
		{
			name: "InvalidConfig",
			cfg: ZapConfig{
				Config: zap.Config{
					Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
					Development: true,
					Encoding:    "invalid",
					EncoderConfig: zapcore.EncoderConfig{
						MessageKey: "message",
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newZapLogger(tt.cfg)
			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("newZapLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZapLogger_Log(t *testing.T) {
	message := make([]string, 5)
	ts := newTestLogSpy(t)
	defer ts.AssertPassed()

	zl := zaptest.NewLogger(ts, zaptest.WrapOptions(zap.AddCallerSkip(2)), zaptest.Level(zap.DebugLevel))
	logger := &zapLogger{
		log: zl,
	}
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
		{"FatalLevel", FatalLevel},
		//{"PanicLevel", PanicLevel},
	}

	for _, tt := range tests {
		logger.Log(ctx, tt.level, msg, keyValues...)
		message = append(message, fmt.Sprintf(`%s	test message	{"key1": "value1"}`, strings.ToUpper(tt.level.String())))
	}
	ts.AssertMessages(message...)
}

// testLogSpy is a testing.TB that captures logged messages.
type testLogSpy struct {
	testing.TB

	failed   bool
	Messages []string
}

func newTestLogSpy(t testing.TB) *testLogSpy {
	return &testLogSpy{TB: t}
}

func (t *testLogSpy) Fail() {
	t.failed = true
}

func (t *testLogSpy) Failed() bool {
	return t.failed
}

func (t *testLogSpy) FailNow() {
	t.Fail()
	t.TB.FailNow()
}

func (t *testLogSpy) Logf(format string, args ...interface{}) {
	// Log messages are in the format,
	//
	//   2017-10-27T13:03:01.000-0700	DEBUG	your message here	{data here}
	//
	// We strip the first part of these messages because we can't really test
	// for the timestamp from these tests.
	m := fmt.Sprintf(format, args...)
	m = m[strings.IndexByte(m, '\t')+1:]
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *testLogSpy) AssertMessages(msgs ...string) {
	assert.Equal(t.TB, msgs, t.Messages, "logged messages did not match")
}

func (t *testLogSpy) AssertPassed() {
	t.assertFailed(false, "expected test to pass")
}

func (t *testLogSpy) AssertFailed() {
	t.assertFailed(true, "expected test to fail")
}

func (t *testLogSpy) assertFailed(v bool, msg string) {
	assert.Equal(t.TB, v, t.failed, msg)
}
