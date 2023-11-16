/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"testing"
)

func TestAdapter_SetLevel(t *testing.T) {
	adapter := NewAdapter(&mockLogger{})
	adapter.SetLevel(InfoLevel)
	if adapter.Level() != InfoLevel {
		t.Errorf("SetLevel() = %v, want %v", adapter.Level(), InfoLevel)
	}
}

func TestAdapter_Enabled(t *testing.T) {
	adapter := NewAdapter(&mockLogger{})
	adapter.SetLevel(InfoLevel)
	if !adapter.Enabled(InfoLevel) {
		t.Errorf("Enabled() = %v, want %v", adapter.Enabled(InfoLevel), true)
	}
}

func TestAdapter_Print(t *testing.T) {
	logger := &mockLogger{}
	adapter := NewAdapter(logger)
	adapter.Print("test")
	if logger.lastMessage != "test" {
		t.Errorf("Print() = %v, want %v", logger.lastMessage, "test")
	}
}

func TestAdapter_Printf(t *testing.T) {
	logger := &mockLogger{}
	adapter := NewAdapter(logger)
	adapter.Printf("test %d", 1)
	if logger.lastMessage != "test 1" {
		t.Errorf("Printf() = %v, want %v", logger.lastMessage, "test 1")
	}
}

func TestAdapter_Printw(t *testing.T) {
	logger := &mockLogger{}
	adapter := NewAdapter(logger)
	adapter.Printw("test", "a")
	if logger.lastMessage != "" || logger.lastKeyvals[0] != "test" || logger.lastKeyvals[1] != "a" {
		t.Errorf("Printw() = %v, want %v", logger.lastMessage, "test")
	}
}

// mockLogger is a Logger implementation for testing purposes.
type mockLogger struct {
	lastMessage string
	lastKeyvals []interface{}
}

func (l *mockLogger) Log(ctx context.Context, level Level, msg string, keyvals ...interface{}) {
	l.lastMessage = msg
	l.lastKeyvals = keyvals
}

func (l *mockLogger) Close() error {
	return nil
}
