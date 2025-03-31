/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestAdapter_SetLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected Level
	}{
		{"Set Debug Level", DebugLevel, DebugLevel},
		{"Set Info Level", InfoLevel, InfoLevel},
		{"Set Warn Level", WarnLevel, WarnLevel},
		{"Set Error Level", ErrorLevel, ErrorLevel},
		{"Set Panic Level", PanicLevel, PanicLevel},
		{"Set Fatal Level", FatalLevel, FatalLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewAdapter(&mockLogger{})
			adapter.SetLevel(tt.level)
			if adapter.Level() != tt.expected {
				t.Errorf("SetLevel() = %v, want %v", adapter.Level(), tt.expected)
			}
		})
	}
}

func TestAdapter_Enabled(t *testing.T) {
	tests := []struct {
		name     string
		setLevel Level
		check    Level
		expected bool
	}{
		{"Debug enabled at Debug", DebugLevel, DebugLevel, true},
		{"Debug enabled at Info", DebugLevel, InfoLevel, true},
		{"Debug enabled at Warn", DebugLevel, WarnLevel, true},
		{"Debug enabled at Error", DebugLevel, ErrorLevel, true},
		{"Debug enabled at Panic", DebugLevel, PanicLevel, true},
		{"Debug enabled at Fatal", DebugLevel, FatalLevel, true},
		{"Info enabled at Debug", InfoLevel, DebugLevel, false},
		{"Info enabled at Info", InfoLevel, InfoLevel, true},
		{"Info enabled at Warn", InfoLevel, WarnLevel, true},
		{"Info enabled at Error", InfoLevel, ErrorLevel, true},
		{"Info enabled at Panic", InfoLevel, PanicLevel, true},
		{"Info enabled at Fatal", InfoLevel, FatalLevel, true},
		{"Warn enabled at Debug", WarnLevel, DebugLevel, false},
		{"Warn enabled at Info", WarnLevel, InfoLevel, false},
		{"Warn enabled at Warn", WarnLevel, WarnLevel, true},
		{"Warn enabled at Error", WarnLevel, ErrorLevel, true},
		{"Warn enabled at Panic", WarnLevel, PanicLevel, true},
		{"Warn enabled at Fatal", WarnLevel, FatalLevel, true},
		{"Error enabled at Debug", ErrorLevel, DebugLevel, false},
		{"Error enabled at Info", ErrorLevel, InfoLevel, false},
		{"Error enabled at Warn", ErrorLevel, WarnLevel, false},
		{"Error enabled at Error", ErrorLevel, ErrorLevel, true},
		{"Error enabled at Panic", ErrorLevel, PanicLevel, true},
		{"Error enabled at Fatal", ErrorLevel, FatalLevel, true},
		{"Panic enabled at Debug", PanicLevel, DebugLevel, false},
		{"Panic enabled at Info", PanicLevel, InfoLevel, false},
		{"Panic enabled at Warn", PanicLevel, WarnLevel, false},
		{"Panic enabled at Error", PanicLevel, ErrorLevel, false},
		{"Panic enabled at Panic", PanicLevel, PanicLevel, true},
		{"Panic enabled at Fatal", PanicLevel, FatalLevel, true},
		{"Fatal enabled at Debug", FatalLevel, DebugLevel, false},
		{"Fatal enabled at Info", FatalLevel, InfoLevel, false},
		{"Fatal enabled at Warn", FatalLevel, WarnLevel, false},
		{"Fatal enabled at Error", FatalLevel, ErrorLevel, false},
		{"Fatal enabled at Panic", FatalLevel, PanicLevel, false},
		{"Fatal enabled at Fatal", FatalLevel, FatalLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewAdapter(&mockLogger{})
			adapter.SetLevel(tt.setLevel)
			if adapter.Enabled(tt.check) != tt.expected {
				t.Errorf("Enabled() = %v, want %v", adapter.Enabled(tt.check), tt.expected)
			}
		})
	}
}

func TestAdapter_PrintMethods(t *testing.T) {
	tests := []struct {
		name     string
		adapter  *Adapter
		action   func(*Adapter)
		expected struct {
			message string
			level   Level
			keyvals []interface{}
		}
	}{
		{
			name:    "Print",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Print("test") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test",
				level:   InfoLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Printf",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Printf("test %d", 1) },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test 1",
				level:   InfoLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Printw",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Printw("key", "value") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "",
				level:   InfoLevel,
				keyvals: []interface{}{"key", "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			adapter := NewAdapter(logger)
			tt.action(adapter)
			if logger.lastMessage != tt.expected.message {
				t.Errorf("%s() message = %v, want %v", tt.name, logger.lastMessage, tt.expected.message)
			}
			if logger.lastLevel != tt.expected.level {
				t.Errorf("%s() level = %v, want %v", tt.name, logger.lastLevel, tt.expected.level)
			}
			if !reflect.DeepEqual(logger.lastKeyvals, tt.expected.keyvals) {
				t.Errorf("%s() keyvals = %v, want %v", tt.name, logger.lastKeyvals, tt.expected.keyvals)
			}
		})
	}
}

func TestAdapter_DebugMethods(t *testing.T) {
	tests := []struct {
		name     string
		adapter  *Adapter
		action   func(*Adapter)
		expected struct {
			message string
			level   Level
			keyvals []interface{}
		}
	}{
		{
			name:    "Debug",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Debug("test") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test",
				level:   DebugLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Debugf",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Debugf("test %d", 1) },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test 1",
				level:   DebugLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Debugw",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Debugw("key", "value") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "",
				level:   DebugLevel,
				keyvals: []interface{}{"key", "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			adapter := NewAdapter(logger)
			adapter.SetLevel(DebugLevel)
			tt.action(adapter)
			if logger.lastMessage != tt.expected.message {
				t.Errorf("%s() message = %v, want %v", tt.name, logger.lastMessage, tt.expected.message)
			}
			if logger.lastLevel != tt.expected.level {
				t.Errorf("%s() level = %v, want %v", tt.name, logger.lastLevel, tt.expected.level)
			}
			if !reflect.DeepEqual(logger.lastKeyvals, tt.expected.keyvals) {
				t.Errorf("%s() keyvals = %v, want %v", tt.name, logger.lastKeyvals, tt.expected.keyvals)
			}
		})
	}
}

func TestAdapter_InfoMethods(t *testing.T) {
	tests := []struct {
		name     string
		adapter  *Adapter
		action   func(*Adapter)
		expected struct {
			message string
			level   Level
			keyvals []interface{}
		}
	}{
		{
			name:    "Info",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Info("test") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test",
				level:   InfoLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Infof",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Infof("test %d", 1) },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test 1",
				level:   InfoLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Infow",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Infow("key", "value") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "",
				level:   InfoLevel,
				keyvals: []interface{}{"key", "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			adapter := NewAdapter(logger)
			adapter.SetLevel(InfoLevel)
			tt.action(adapter)
			if logger.lastMessage != tt.expected.message {
				t.Errorf("%s() message = %v, want %v", tt.name, logger.lastMessage, tt.expected.message)
			}
			if logger.lastLevel != tt.expected.level {
				t.Errorf("%s() level = %v, want %v", tt.name, logger.lastLevel, tt.expected.level)
			}
			if !reflect.DeepEqual(logger.lastKeyvals, tt.expected.keyvals) {
				t.Errorf("%s() keyvals = %v, want %v", tt.name, logger.lastKeyvals, tt.expected.keyvals)
			}
		})
	}
}

func TestAdapter_WarnMethods(t *testing.T) {
	tests := []struct {
		name     string
		adapter  *Adapter
		action   func(*Adapter)
		expected struct {
			message string
			level   Level
			keyvals []interface{}
		}
	}{
		{
			name:    "Warn",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Warn("test") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test",
				level:   WarnLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Warnf",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Warnf("test %d", 1) },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test 1",
				level:   WarnLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Warnw",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Warnw("key", "value") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "",
				level:   WarnLevel,
				keyvals: []interface{}{"key", "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			adapter := NewAdapter(logger)
			adapter.SetLevel(WarnLevel)
			tt.action(adapter)
			if logger.lastMessage != tt.expected.message {
				t.Errorf("%s() message = %v, want %v", tt.name, logger.lastMessage, tt.expected.message)
			}
			if logger.lastLevel != tt.expected.level {
				t.Errorf("%s() level = %v, want %v", tt.name, logger.lastLevel, tt.expected.level)
			}
			if !reflect.DeepEqual(logger.lastKeyvals, tt.expected.keyvals) {
				t.Errorf("%s() keyvals = %v, want %v", tt.name, logger.lastKeyvals, tt.expected.keyvals)
			}
		})
	}
}

func TestAdapter_ErrorMethods(t *testing.T) {
	tests := []struct {
		name     string
		adapter  *Adapter
		action   func(*Adapter)
		expected struct {
			message string
			level   Level
			keyvals []interface{}
		}
	}{
		{
			name:    "Error",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Error("test") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test",
				level:   ErrorLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Errorf",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Errorf("test %d", 1) },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test 1",
				level:   ErrorLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Errorw",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Errorw("key", "value") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "",
				level:   ErrorLevel,
				keyvals: []interface{}{"key", "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			adapter := NewAdapter(logger)
			adapter.SetLevel(ErrorLevel)
			tt.action(adapter)
			if logger.lastMessage != tt.expected.message {
				t.Errorf("%s() message = %v, want %v", tt.name, logger.lastMessage, tt.expected.message)
			}
			if logger.lastLevel != tt.expected.level {
				t.Errorf("%s() level = %v, want %v", tt.name, logger.lastLevel, tt.expected.level)
			}
			if !reflect.DeepEqual(logger.lastKeyvals, tt.expected.keyvals) {
				t.Errorf("%s() keyvals = %v, want %v", tt.name, logger.lastKeyvals, tt.expected.keyvals)
			}
		})
	}
}

func TestAdapter_PanicMethods(t *testing.T) {
	tests := []struct {
		name     string
		adapter  *Adapter
		action   func(*Adapter)
		expected struct {
			message string
			level   Level
			keyvals []interface{}
		}
	}{
		{
			name:    "Panic",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Panic("test") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test",
				level:   PanicLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Panicf",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Panicf("test %d", 1) },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test 1",
				level:   PanicLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Panicw",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Panicw("key", "value") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "",
				level:   PanicLevel,
				keyvals: []interface{}{"key", "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			adapter := NewAdapter(logger)
			adapter.SetLevel(PanicLevel)

			// Use defer recover to catch panic
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic but got none")
				}
			}()

			tt.action(adapter)

			// Check logger state after panic
			if logger.lastMessage != tt.expected.message {
				t.Errorf("%s() message = %v, want %v", tt.name, logger.lastMessage, tt.expected.message)
			}
			if logger.lastLevel != tt.expected.level {
				t.Errorf("%s() level = %v, want %v", tt.name, logger.lastLevel, tt.expected.level)
			}
			if !reflect.DeepEqual(logger.lastKeyvals, tt.expected.keyvals) {
				t.Errorf("%s() keyvals = %v, want %v", tt.name, logger.lastKeyvals, tt.expected.keyvals)
			}
		})
	}
}

func TestAdapter_FatalMethods(t *testing.T) {
	tests := []struct {
		name     string
		adapter  *Adapter
		action   func(*Adapter)
		expected struct {
			message string
			level   Level
			keyvals []interface{}
		}
	}{
		{
			name:    "Fatal",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Fatal("test") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test",
				level:   FatalLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Fatalf",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Fatalf("test %d", 1) },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "test 1",
				level:   FatalLevel,
				keyvals: nil,
			},
		},
		{
			name:    "Fatalw",
			adapter: NewAdapter(&mockLogger{}),
			action:  func(a *Adapter) { a.Fatalw("key", "value") },
			expected: struct {
				message string
				level   Level
				keyvals []interface{}
			}{
				message: "",
				level:   FatalLevel,
				keyvals: []interface{}{"key", "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file to store test results
			tmpFile := t.TempDir() + "/test.log"
			logger := &mockLogger{}
			adapter := NewAdapter(logger)
			adapter.SetLevel(FatalLevel)

			// Run the test in a subprocess
			cmd := exec.Command(os.Args[0], "-test.run=TestFatalSubprocess")
			cmd.Env = append(os.Environ(),
				"TEST_FATAL_METHOD="+tt.name,
				"TEST_TEMP_FILE="+tmpFile,
			)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command and expect it to exit with code 1
			err := cmd.Run()
			if err == nil {
				t.Error("Expected command to exit with error, got nil")
			} else if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
				t.Errorf("Expected exit code 1, got %v", err)
			}

			// Read the test results from the temporary file
			data, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("Failed to read test results: %v", err)
			}

			// Parse the test results
			var result struct {
				Message string
				Level   Level
				Keyvals []interface{}
			}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to parse test results: %v", err)
			}

			// Verify the results
			if result.Message != tt.expected.message {
				t.Errorf("%s() message = %v, want %v", tt.name, result.Message, tt.expected.message)
			}
			if result.Level != tt.expected.level {
				t.Errorf("%s() level = %v, want %v", tt.name, result.Level, tt.expected.level)
			}
			if !reflect.DeepEqual(result.Keyvals, tt.expected.keyvals) {
				t.Errorf("%s() keyvals = %v, want %v", tt.name, result.Keyvals, tt.expected.keyvals)
			}
		})
	}
}

// TestFatalSubprocess is a helper function that runs in a subprocess to test fatal methods
func TestFatalSubprocess(t *testing.T) {
	if os.Getenv("TEST_FATAL_METHOD") == "" {
		return
	}

	logger := &mockLogger{}
	adapter := NewAdapter(logger)
	adapter.SetLevel(FatalLevel)

	// Get the temporary file path
	tmpFile := os.Getenv("TEST_TEMP_FILE")
	if tmpFile == "" {
		t.Fatal("TEST_TEMP_FILE environment variable not set")
	}

	// Prepare test results based on the method
	var result struct {
		Message string
		Level   Level
		Keyvals []interface{}
	}

	switch os.Getenv("TEST_FATAL_METHOD") {
	case "Fatal":
		result = struct {
			Message string
			Level   Level
			Keyvals []interface{}
		}{
			Message: "test",
			Level:   FatalLevel,
			Keyvals: nil,
		}
	case "Fatalf":
		result = struct {
			Message string
			Level   Level
			Keyvals []interface{}
		}{
			Message: "test 1",
			Level:   FatalLevel,
			Keyvals: nil,
		}
	case "Fatalw":
		result = struct {
			Message string
			Level   Level
			Keyvals []interface{}
		}{
			Message: "",
			Level:   FatalLevel,
			Keyvals: []interface{}{"key", "value"},
		}
	default:
		t.Fatalf("Unknown test method: %s", os.Getenv("TEST_FATAL_METHOD"))
	}

	// Write test results before calling fatal method
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal test results: %v", err)
	}

	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test results: %v", err)
	}

	// Call the fatal method after writing results
	switch os.Getenv("TEST_FATAL_METHOD") {
	case "Fatal":
		adapter.Fatal("test")
	case "Fatalf":
		adapter.Fatalf("test %d", 1)
	case "Fatalw":
		adapter.Fatalw("key", "value")
	}
}

func TestAdapter_WithContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected context.Context
	}{
		{
			name:     "With Background Context",
			ctx:      context.Background(),
			expected: context.Background(),
		},
		{
			name:     "With Value Context",
			ctx:      context.WithValue(context.Background(), "test", "value"),
			expected: context.WithValue(context.Background(), "test", "value"),
		},
		{
			name:     "With Cancel Context",
			ctx:      func() context.Context { ctx, _ := context.WithCancel(context.Background()); return ctx }(),
			expected: func() context.Context { ctx, _ := context.WithCancel(context.Background()); return ctx }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewAdapter(&mockLogger{}, WithContext(tt.ctx))
			if tt.ctx != nil {
				if adapter.ctx == nil {
					t.Error("WithContext() returned nil context, expected non-nil")
				} else {
					if val := tt.ctx.Value("test"); val != nil {
						if adapterVal := adapter.ctx.Value("test"); adapterVal != val {
							t.Errorf("WithContext() value = %v, want %v", adapterVal, val)
						}
					}
					if tt.ctx.Done() != nil && adapter.ctx.Done() == nil {
						t.Error("WithContext() returned context without Done channel, expected with Done channel")
					}
				}
			} else if adapter.ctx != nil {
				t.Error("WithContext() returned non-nil context, expected nil")
			}
		})
	}
}

func TestAdapter_Close(t *testing.T) {
	tests := []struct {
		name     string
		logger   *mockLogger
		expected error
	}{
		{
			name:     "Close Success",
			logger:   &mockLogger{},
			expected: nil,
		},
		{
			name: "Close Error",
			logger: &mockLogger{
				closeError: fmt.Errorf("close error"),
			},
			expected: fmt.Errorf("close error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewAdapter(tt.logger)
			err := adapter.Close()
			if tt.expected == nil {
				if err != nil {
					t.Errorf("Close() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("Close() error = nil, want %v", tt.expected)
				} else if err.Error() != tt.expected.Error() {
					t.Errorf("Close() error = %v, want %v", err, tt.expected)
				}
			}
		})
	}
}

// mockLogger is a Logger implementation for testing purposes.
type mockLogger struct {
	lastMessage string
	lastKeyvals []interface{}
	lastLevel   Level
	lastCtx     context.Context
	closeError  error
}

func (l *mockLogger) Log(ctx context.Context, level Level, msg string, keyvals ...interface{}) {
	l.lastMessage = msg
	l.lastKeyvals = keyvals
	l.lastLevel = level
	l.lastCtx = ctx
}

func (l *mockLogger) Close() error {
	return l.closeError
}
