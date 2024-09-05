/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"strings"
	"sync"
)

var (
	defaultAdapter *Adapter
	cfg            = Config{DefaultLevel: InfoLevel, Named: make(map[string]Level)}
	adapters       = make(map[string]*Adapter)
	mu             sync.Mutex
)

type Config struct {
	DefaultLevel Level            `json:"level" yaml:"level"`
	Named        map[string]Level `json:"named" yaml:"named"`
}

// findAdapterLevel 按照 . 分隔符，逐级向下查找
// 例如：a.b.c -> a.b -> a
// 如果找不到，则使用默认的 DefaultLevel
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

func SetDefaultAdapter(adapter *Adapter) {
	mu.Lock()
	defer mu.Unlock()

	defaultAdapter = adapter
}

func SetDefaultLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()

	defaultAdapter = NewAdapter(logger, WithLevel(cfg.DefaultLevel))
}

// Debug uses fmt.Sprint to construct and logger a message.
func Debug(args ...interface{}) {
	defaultAdapter.Debug(args...)
}

// Debugf uses fmt.Sprintf to logger a templated message.
func Debugf(msg string, args ...interface{}) {
	defaultAdapter.Debugf(msg, args...)
}

// Debugw uses key,value to construct and logger a message.
func Debugw(keyvals ...interface{}) {
	defaultAdapter.Debugw(keyvals...)

}

// Info uses fmt.Sprint to construct and logger a message.
func Info(args ...interface{}) {
	defaultAdapter.Info(args...)
}

// Infof uses fmt.Sprintf to logger a templated message.
func Infof(msg string, args ...interface{}) {
	defaultAdapter.Infof(msg, args...)
}

// Infow uses key,value to construct and logger a message.
func Infow(keyvals ...interface{}) {
	defaultAdapter.Infow(keyvals...)

}

// Warn uses fmt.Sprint to construct and logger a message.
func Warn(args ...interface{}) {
	defaultAdapter.Warn(args...)
}

// Warnf uses fmt.Sprintf to logger a templated message.
func Warnf(msg string, args ...interface{}) {
	defaultAdapter.Warnf(msg, args...)
}

// Warnw uses key,value to construct and logger a message.
func Warnw(keyvals ...interface{}) {
	defaultAdapter.Warnw(keyvals...)

}

// Error uses fmt.Sprint to construct and logger a message.
func Error(args ...interface{}) {
	defaultAdapter.Error(args...)
}

// Errorf uses fmt.Sprintf to logger a templated message.
func Errorf(msg string, args ...interface{}) {
	defaultAdapter.Errorf(msg, args...)
}

// Errorw uses key,value to construct and logger a message.
func Errorw(keyvals ...interface{}) {
	defaultAdapter.Errorw(keyvals...)

}

// Print uses fmt.Sprint to construct and logger a message.
func Print(args ...interface{}) {
	defaultAdapter.Info(args...)
}

// Printf uses fmt.Sprintf to logger a templated message.
func Printf(msg string, args ...interface{}) {
	defaultAdapter.Infof(msg, args...)
}

// Printw uses key,value to construct and logger a message.
func Printw(keyvals ...interface{}) {
	defaultAdapter.Printw(keyvals...)

}

// Panic uses fmt.Sprint to construct and logger a message, then panics.
func Panic(args ...interface{}) {
	defaultAdapter.Panic(args...)
}

// Panicf uses fmt.Sprintf to logger a templated message, then panics.
func Panicf(msg string, args ...interface{}) {
	defaultAdapter.Panicf(msg, args...)
}

// Panicw uses key,value to construct and logger a message.
func Panicw(keyvals ...interface{}) {
	defaultAdapter.Panicw(keyvals...)

}

// Fatal uses fmt.Sprint to construct and logger a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	defaultAdapter.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to logger a templated message, then calls os.Exit.
func Fatalf(msg string, args ...interface{}) {
	defaultAdapter.Fatalf(msg, args...)
}

// Fatalw uses key,value to construct and logger a message.
func Fatalw(keyvals ...interface{}) {
	defaultAdapter.Fatalw(keyvals...)

}

// Flush flushing any buffered log entries. Applications should take care to call Sync before exiting.
func Flush() error {
	return defaultAdapter.Flush()
}

// SetLevel alters the logging level.
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
