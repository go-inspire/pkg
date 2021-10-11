package log

import (
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
)

type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

type Slogger interface {
	// Print uses fmt.Sprint to construct and logger a message.
	Print(args ...interface{})
	// Printf uses fmt.Sprintf to logger a templated message.
	Printf(msg string, args ...interface{})
}

type Logger interface {
	Slogger
	// Debug uses fmt.Sprint to construct and logger a message.
	Debug(args ...interface{})
	// Debugf uses fmt.Sprintf to logger a templated message.
	Debugf(msg string, args ...interface{})

	// Info uses fmt.Sprint to construct and logger a message.
	Info(args ...interface{})
	// Infof uses fmt.Sprintf to logger a templated message.
	Infof(msg string, args ...interface{})

	// Warn uses fmt.Sprint to construct and logger a message.
	Warn(args ...interface{})
	// Warnf uses fmt.Sprintf to logger a templated message.
	Warnf(msg string, args ...interface{})

	// Error uses fmt.Sprint to construct and logger a message.
	Error(args ...interface{})
	// Errorf uses fmt.Sprintf to logger a templated message.
	Errorf(msg string, args ...interface{})

	// Panic uses fmt.Sprint to construct and logger a message, then panics.
	Panic(args ...interface{})
	// Panicf uses fmt.Sprintf to logger a templated message, then panics.
	Panicf(msg string, args ...interface{})

	// Fatal uses fmt.Sprint to construct and logger a message, then calls os.Exit.
	Fatal(args ...interface{})
	// Fatalf uses fmt.Sprintf to logger a templated message, then calls os.Exit.
	Fatalf(msg string, args ...interface{})

	// Flush Flushing any buffered logger entries. Applications should take care to call Sync before exiting.
	Flush() error

	// Named adds a new path segment to the logger's name. Segments are joined by
	// periods. By default, Loggers are unnamed.
	Named(s string) Logger

	//SetLevel alters the logging level.
	SetLevel(lvl Level) Logger
}

var (
	loggers = make(map[string]Logger)
	mu      sync.Mutex
	std     Logger
)

func Named(s string) Logger {
	mu.Lock()
	defer mu.Unlock()

	s = strings.ToLower(s)
	l, ok := loggers[s]
	if !ok {
		l = std.Named(s)
		loggers[s] = l
	}
	return l
}

// Debug uses fmt.Sprint to construct and logger a message.
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Debugf uses fmt.Sprintf to logger a templated message.
func Debugf(msg string, args ...interface{}) {
	std.Debugf(msg, args...)
}

// Info uses fmt.Sprint to construct and logger a message.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Infof uses fmt.Sprintf to logger a templated message.
func Infof(msg string, args ...interface{}) {
	std.Infof(msg, args...)
}

// Warn uses fmt.Sprint to construct and logger a message.
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Warnf uses fmt.Sprintf to logger a templated message.
func Warnf(msg string, args ...interface{}) {
	std.Warnf(msg, args...)
}

// Error uses fmt.Sprint to construct and logger a message.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Errorf uses fmt.Sprintf to logger a templated message.
func Errorf(msg string, args ...interface{}) {
	std.Errorf(msg, args...)
}

// Print uses fmt.Sprint to construct and logger a message.
func Print(args ...interface{}) {
	std.Info(args...)
}

// Printf uses fmt.Sprintf to logger a templated message.
func Printf(msg string, args ...interface{}) {
	std.Infof(msg, args...)
}

// Panic uses fmt.Sprint to construct and logger a message, then panics.
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// Panicf uses fmt.Sprintf to logger a templated message, then panics.
func Panicf(msg string, args ...interface{}) {
	std.Panicf(msg, args...)
}

// Fatal uses fmt.Sprint to construct and logger a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to logger a templated message, then calls os.Exit.
func Fatalf(msg string, args ...interface{}) {
	std.Fatalf(msg, args...)
}

// Flush flushing any buffered log entries. Applications should take care to call Sync before exiting.
func Flush() error {
	for _, l := range loggers {
		_ = l.Flush()
	}
	return std.Flush()
}

// SetLevel alters the logging level.
func SetLevel(name, lvl string) Logger {
	name = strings.ToLower(name)
	logger, ok := loggers[name]
	if !ok {
		return nil
	}
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(lvl)); err == nil {
		return logger.SetLevel(l)
	}
	return nil
}
