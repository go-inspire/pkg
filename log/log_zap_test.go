/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestZapLogger_Log tests the Log method of zapLogger
func TestZapLogger_Log(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		msg      string
		keyvals  []interface{}
		expected struct {
			level   Level
			msg     string
			keyvals []interface{}
		}
	}{
		{
			name:    "Debug with message",
			level:   DebugLevel,
			msg:     "debug message",
			keyvals: nil,
			expected: struct {
				level   Level
				msg     string
				keyvals []interface{}
			}{
				level:   DebugLevel,
				msg:     "debug message",
				keyvals: nil,
			},
		},
		{
			name:    "Info with key-value pairs",
			level:   InfoLevel,
			msg:     "info message",
			keyvals: []interface{}{"key", "value"},
			expected: struct {
				level   Level
				msg     string
				keyvals []interface{}
			}{
				level:   InfoLevel,
				msg:     "info message",
				keyvals: []interface{}{"key", "value"},
			},
		},
		{
			name:    "Warn with multiple key-value pairs",
			level:   WarnLevel,
			msg:     "warn message",
			keyvals: []interface{}{"key1", "value1", "key2", "value2"},
			expected: struct {
				level   Level
				msg     string
				keyvals []interface{}
			}{
				level:   WarnLevel,
				msg:     "warn message",
				keyvals: []interface{}{"key1", "value1", "key2", "value2"},
			},
		},
		{
			name:    "Error with zap.Field",
			level:   ErrorLevel,
			msg:     "error message",
			keyvals: []interface{}{zap.String("field", "value")},
			expected: struct {
				level   Level
				msg     string
				keyvals []interface{}
			}{
				level:   ErrorLevel,
				msg:     "error message",
				keyvals: []interface{}{zap.String("field", "value")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for log output
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.log")

			// Create zap logger configuration
			config := ZapConfig{
				Config: zap.NewProductionConfig(),
			}
			config.OutputPaths = []string{tmpFile}
			config.ErrorOutputPaths = []string{tmpFile}
			config.Config.Level = zap.NewAtomicLevelAt(zapcore.Level(tt.level))
			config.Config.Encoding = "json"
			config.Config.EncoderConfig = zap.NewProductionEncoderConfig()

			// Create logger
			logger := newZapLogger(config)
			if logger == nil {
				t.Fatal("Failed to create logger")
			}
			defer logger.Close()

			// Log message
			logger.Log(context.Background(), tt.level, tt.msg, tt.keyvals...)

			// Read log file
			data, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}

			// Get the last line (most recent log entry)
			lines := strings.Split(string(data), "\n")
			var lastLine string
			for i := len(lines) - 1; i >= 0; i-- {
				if lines[i] != "" {
					lastLine = lines[i]
					break
				}
			}

			// Parse log entry
			var entry struct {
				Level string `json:"level"`
				Msg   string `json:"msg"`
			}
			if err := json.Unmarshal([]byte(lastLine), &entry); err != nil {
				t.Fatalf("Failed to parse log entry: %v", err)
			}

			// Verify level
			if entry.Level != tt.level.String() {
				t.Errorf("Log level = %v, want %v", entry.Level, tt.level)
			}

			// Verify message
			if entry.Msg != tt.msg {
				t.Errorf("Log message = %v, want %v", entry.Msg, tt.msg)
			}
		})
	}
}

// TestZapLogger_Close tests the Close method of zapLogger
func TestZapLogger_Close(t *testing.T) {
	// Create a temporary file for log output
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.log")

	// Create zap logger configuration
	config := ZapConfig{
		Config: zap.NewProductionConfig(),
	}
	config.OutputPaths = []string{tmpFile}
	config.ErrorOutputPaths = []string{tmpFile}

	// Create logger
	logger := newZapLogger(config)
	if logger == nil {
		t.Fatal("Failed to create logger")
	}

	// Test Close method
	err := logger.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

// TestZapLogger_Config tests the configuration of zapLogger
func TestZapLogger_Config(t *testing.T) {
	tests := []struct {
		name     string
		config   ZapConfig
		expected struct {
			level Level
			named map[string]Level
		}
	}{
		{
			name: "Default configuration",
			config: ZapConfig{
				Config: zap.NewProductionConfig(),
			},
			expected: struct {
				level Level
				named map[string]Level
			}{
				level: InfoLevel,
				named: nil,
			},
		},
		{
			name: "Custom level configuration",
			config: ZapConfig{
				Config: zap.NewProductionConfig(),
				Named: map[string]Level{
					"test": DebugLevel,
				},
			},
			expected: struct {
				level Level
				named map[string]Level
			}{
				level: DebugLevel,
				named: map[string]Level{
					"test": DebugLevel,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the expected level in the config
			tt.config.Config.Level = zap.NewAtomicLevelAt(zapcore.Level(tt.expected.level))
			tt.config.Config.Encoding = "json"
			tt.config.Config.EncoderConfig = zap.NewProductionEncoderConfig()

			// Create logger
			logger := newZapLogger(tt.config)
			if logger == nil {
				t.Fatal("Failed to create logger")
			}
			defer logger.Close()

			// Verify configuration
			if tt.config.Level.Level() != tt.expected.level {
				t.Errorf("Log level = %v, want %v", tt.config.Level.Level(), tt.expected.level)
			}

			if tt.config.Named != nil {
				for name, level := range tt.expected.named {
					if tt.config.Named[name] != level {
						t.Errorf("Named level for %s = %v, want %v", name, tt.config.Named[name], level)
					}
				}
			}

			// Test logging with configured level
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.log")
			tt.config.OutputPaths = []string{tmpFile}
			tt.config.ErrorOutputPaths = []string{tmpFile}

			logger = newZapLogger(tt.config)
			if logger == nil {
				t.Fatal("Failed to create logger with updated config")
			}
			defer logger.Close()

			// Log a test message
			logger.Log(context.Background(), tt.expected.level, "test message")

			// Read and verify log file
			data, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}

			// Get the last line (most recent log entry)
			lines := strings.Split(string(data), "\n")
			var lastLine string
			for i := len(lines) - 1; i >= 0; i-- {
				if lines[i] != "" {
					lastLine = lines[i]
					break
				}
			}

			// Parse log entry
			var entry struct {
				Level string `json:"level"`
			}
			if err := json.Unmarshal([]byte(lastLine), &entry); err != nil {
				t.Fatalf("Failed to parse log entry: %v", err)
			}

			if entry.Level != tt.expected.level.String() {
				t.Errorf("Log entry level = %v, want %v", entry.Level, tt.expected.level)
			}
		})
	}
}

// TestZapLogger_KeyValuesToField tests the keyValuesToField function
func TestZapLogger_KeyValuesToField(t *testing.T) {
	tests := []struct {
		name     string
		args     []interface{}
		expected struct {
			field zap.Field
			rest  []interface{}
		}
	}{
		{
			name: "String key-value pair",
			args: []interface{}{"key", "value"},
			expected: struct {
				field zap.Field
				rest  []interface{}
			}{
				field: zap.Any("key", "value"),
				rest:  nil,
			},
		},
		{
			name: "Single string",
			args: []interface{}{"value"},
			expected: struct {
				field zap.Field
				rest  []interface{}
			}{
				field: zap.String(badKey, "value"),
				rest:  nil,
			},
		},
		{
			name: "Zap field",
			args: []interface{}{zap.String("field", "value")},
			expected: struct {
				field zap.Field
				rest  []interface{}
			}{
				field: zap.String("field", "value"),
				rest:  nil,
			},
		},
		{
			name: "Non-string key",
			args: []interface{}{123, "value"},
			expected: struct {
				field zap.Field
				rest  []interface{}
			}{
				field: zap.Any(badKey, 123),
				rest:  []interface{}{"value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, rest := keyValuesToField(tt.args)
			if field.Key != tt.expected.field.Key {
				t.Errorf("Field key = %v, want %v", field.Key, tt.expected.field.Key)
			}
			if field.Type != tt.expected.field.Type {
				t.Errorf("Field type = %v, want %v", field.Type, tt.expected.field.Type)
			}
			if field.Interface != tt.expected.field.Interface {
				t.Errorf("Field value = %v, want %v", field.Interface, tt.expected.field.Interface)
			}
			if len(rest) != len(tt.expected.rest) {
				t.Errorf("Rest length = %v, want %v", len(rest), len(tt.expected.rest))
			}
		})
	}
}
