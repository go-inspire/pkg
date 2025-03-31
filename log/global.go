/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

// Package log provides a flexible logging system with support for multiple logging backends
// and hierarchical log levels. It offers both structured and unstructured logging capabilities.
package log

import (
	"strings"
	"sync"
)

var (
	// defaultAdapter is the default logger adapter used by the package-level logging functions.
	defaultAdapter *Adapter
	// cfg holds the global logging configuration including default and named log levels.
	cfg = Config{DefaultLevel: InfoLevel, Named: make(map[string]Level)}
	// adapters stores named logger adapters for different components or modules.
	adapters = make(map[string]*Adapter)
	// mu protects concurrent access to the global variables.
	mu sync.Mutex
)

// Config represents the global logging configuration.
// It allows setting default log levels and named log levels for specific components.
type Config struct {
	// DefaultLevel is the default logging level used when no specific level is set.
	DefaultLevel Level `json:"level" yaml:"level"`
	// Named contains component-specific log levels, indexed by component name.
	Named map[string]Level `json:"named" yaml:"named"`
}

// findConfigLevel searches for the appropriate log level for a given component name.
// It implements a hierarchical lookup where it tries to find the most specific level
// by progressively removing segments from the name (e.g., "a.b.c" -> "a.b" -> "a").
// If no specific level is found, it returns the default level.
func findConfigLevel(cfg *Config, s string) Level {
	name := s
	for {
		if lvl, exists := cfg.Named[name]; exists {
			return lvl
		}
		if index := strings.LastIndex(name, "."); index > 0 {
			name = name[:index]
		} else {
			return cfg.DefaultLevel
		}
	}
}

// Named returns a logger adapter for the specified component name.
// If a logger for the name already exists, it returns the existing one;
// otherwise, it creates a new logger with the appropriate level based on the configuration.
func Named(s string) *Adapter {
	mu.Lock()
	defer mu.Unlock()

	s = strings.ToLower(s)
	a, ok := adapters[s]
	if !ok {
		level := findConfigLevel(&cfg, s)
		a = NewAdapter(defaultAdapter.logger, WithLevel(level))
		adapters[s] = a
	}
	return a
}

// SetConfig updates the global logging configuration.
// It also updates the log levels of all existing named adapters to match the new configuration.
func SetConfig(config Config) {
	mu.Lock()
	defer mu.Unlock()

	cfg = config

	if len(cfg.Named) > 0 {
		for k, a := range adapters {
			lvl := findConfigLevel(&cfg, k)
			a.SetLevel(lvl)
			adapters[k] = a
		}
	}
}

// SetDefaultAdapter sets the default logger adapter used by package-level logging functions.
// This adapter is used when no specific named logger is requested.
func SetDefaultAdapter(adapter *Adapter) {
	mu.Lock()
	defer mu.Unlock()

	defaultAdapter = adapter
}

// SetDefaultLogger creates a new adapter with the specified logger and the current default level,
// then sets it as the default adapter.
func SetDefaultLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()

	defaultAdapter = NewAdapter(logger, WithLevel(cfg.DefaultLevel))
}

// Debug logs a message at debug level using fmt.Sprint to construct the message.
func Debug(args ...interface{}) {
	defaultAdapter.Debug(args...)
}

// Debugf logs a formatted message at debug level using fmt.Sprintf.
func Debugf(msg string, args ...interface{}) {
	defaultAdapter.Debugf(msg, args...)
}

// Debugw logs a message at debug level with key-value pairs.
func Debugw(keyvals ...interface{}) {
	defaultAdapter.Debugw(keyvals...)
}

// Info logs a message at info level using fmt.Sprint to construct the message.
func Info(args ...interface{}) {
	defaultAdapter.Info(args...)
}

// Infof logs a formatted message at info level using fmt.Sprintf.
func Infof(msg string, args ...interface{}) {
	defaultAdapter.Infof(msg, args...)
}

// Infow logs a message at info level with key-value pairs.
func Infow(keyvals ...interface{}) {
	defaultAdapter.Infow(keyvals...)
}

// Warn logs a message at warn level using fmt.Sprint to construct the message.
func Warn(args ...interface{}) {
	defaultAdapter.Warn(args...)
}

// Warnf logs a formatted message at warn level using fmt.Sprintf.
func Warnf(msg string, args ...interface{}) {
	defaultAdapter.Warnf(msg, args...)
}

// Warnw logs a message at warn level with key-value pairs.
func Warnw(keyvals ...interface{}) {
	defaultAdapter.Warnw(keyvals...)
}

// Error logs a message at error level using fmt.Sprint to construct the message.
func Error(args ...interface{}) {
	defaultAdapter.Error(args...)
}

// Errorf logs a formatted message at error level using fmt.Sprintf.
func Errorf(msg string, args ...interface{}) {
	defaultAdapter.Errorf(msg, args...)
}

// Errorw logs a message at error level with key-value pairs.
func Errorw(keyvals ...interface{}) {
	defaultAdapter.Errorw(keyvals...)
}

// Print logs a message at info level using fmt.Sprint to construct the message.
// It is an alias for Info for compatibility with standard logging interfaces.
func Print(args ...interface{}) {
	defaultAdapter.Info(args...)
}

// Printf logs a formatted message at info level using fmt.Sprintf.
// It is an alias for Infof for compatibility with standard logging interfaces.
func Printf(msg string, args ...interface{}) {
	defaultAdapter.Infof(msg, args...)
}

// Printw logs a message at info level with key-value pairs.
// It is an alias for Infow for compatibility with standard logging interfaces.
func Printw(keyvals ...interface{}) {
	defaultAdapter.Printw(keyvals...)
}

// Panic logs a message at panic level using fmt.Sprint to construct the message,
// then panics with the constructed message.
func Panic(args ...interface{}) {
	defaultAdapter.Panic(args...)
}

// Panicf logs a formatted message at panic level using fmt.Sprintf,
// then panics with the formatted message.
func Panicf(msg string, args ...interface{}) {
	defaultAdapter.Panicf(msg, args...)
}

// Panicw logs a message at panic level with key-value pairs,
// then panics with the message.
func Panicw(keyvals ...interface{}) {
	defaultAdapter.Panicw(keyvals...)
}

// Fatal logs a message at fatal level using fmt.Sprint to construct the message,
// then calls os.Exit(1).
func Fatal(args ...interface{}) {
	defaultAdapter.Fatal(args...)
}

// Fatalf logs a formatted message at fatal level using fmt.Sprintf,
// then calls os.Exit(1).
func Fatalf(msg string, args ...interface{}) {
	defaultAdapter.Fatalf(msg, args...)
}

// Fatalw logs a message at fatal level with key-value pairs,
// then calls os.Exit(1).
func Fatalw(keyvals ...interface{}) {
	defaultAdapter.Fatalw(keyvals...)
}

// Flush ensures that any buffered log entries are written.
// Applications should call this before exiting to ensure all logs are written.
func Flush() error {
	return defaultAdapter.Flush()
}

// SetLevel changes the log level for a named logger.
// If the named logger doesn't exist, the call is ignored.
func SetLevel(name, lvl string) {
	name = strings.ToLower(name)
	a, ok := adapters[name]
	if !ok {
		return
	}
	var l Level
	if err := l.UnmarshalText([]byte(lvl)); err == nil {
		a.SetLevel(l)
	}
}
