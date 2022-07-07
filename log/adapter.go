package log

import (
	"fmt"
	"os"
	"sync/atomic"
)

// Option is Helper option.
type Option func(*Adapter)

type Adapter struct {
	logger Logger
	lvl    int32
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
		lvl:    int32(InfoLevel),
	}
	for _, o := range opts {
		o(la)
	}
	return la
}

// Level returns the minimum enabled log level.
func (la *Adapter) Level() Level {
	return Level(int8(atomic.LoadInt32(&la.lvl)))
}

// SetLevel alters the logging level.
func (la *Adapter) SetLevel(l Level) *Adapter {
	atomic.StoreInt32(&la.lvl, int32(l))
	return la
}

// Enabled implements the zapcore.LevelEnabler interface, which allows the
// AtomicLevel to be used in place of traditional static levels.
func (la *Adapter) Enabled(l Level) bool {
	return la.Level().Enabled(l)
}

// Info uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Print(args ...interface{}) {
	if la.Enabled(InfoLevel) {
		la.logger.Log(InfoLevel, sprint(args...))
	}

}

// Infof uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Printf(msg string, args ...interface{}) {
	if la.Enabled(InfoLevel) {
		la.logger.Log(InfoLevel, fmt.Sprintf(msg, args...))
	}
}

// Printw uses key,value to construct and logger a message.
func (la *Adapter) Printw(keyvals ...interface{}) {
	if la.Enabled(InfoLevel) {
		la.logger.Log(InfoLevel, "", keyvals...)
	}
}

// Debug uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Debug(args ...interface{}) {
	if la.Enabled(DebugLevel) {
		la.logger.Log(DebugLevel, sprint(args...))
	}

}

// Debugf uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Debugf(msg string, args ...interface{}) {
	if la.Enabled(DebugLevel) {
		la.logger.Log(DebugLevel, fmt.Sprintf(msg, args...))
	}
}

// Debugw uses key,value to construct and logger a message.
func (la *Adapter) Debugw(keyvals ...interface{}) {
	if la.Enabled(DebugLevel) {
		la.logger.Log(DebugLevel, "", keyvals...)
	}

}

// Info uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Info(args ...interface{}) {
	if la.Enabled(InfoLevel) {
		la.logger.Log(InfoLevel, sprint(args...))
	}
}

// Infof uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Infof(msg string, args ...interface{}) {
	if la.Enabled(InfoLevel) {
		la.logger.Log(InfoLevel, fmt.Sprintf(msg, args...))
	}
}

// Infow uses key,value to construct and logger a message.
func (la *Adapter) Infow(keyvals ...interface{}) {
	if la.Enabled(InfoLevel) {
		la.logger.Log(InfoLevel, "", keyvals...)
	}

}

// Warn uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Warn(args ...interface{}) {
	if la.Enabled(WarnLevel) {
		la.logger.Log(WarnLevel, sprint(args...))
	}
}

// Warnf uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Warnf(msg string, args ...interface{}) {
	if la.Enabled(WarnLevel) {
		la.logger.Log(WarnLevel, fmt.Sprintf(msg, args...))
	}
}

// Warnw uses key,value to construct and logger a message.
func (la *Adapter) Warnw(keyvals ...interface{}) {
	if la.Enabled(WarnLevel) {
		la.logger.Log(WarnLevel, "", keyvals...)
	}

}

// Error uses fmt.Sprint to construct and logger a message.
func (la *Adapter) Error(args ...interface{}) {
	if la.Enabled(ErrorLevel) {
		la.logger.Log(ErrorLevel, sprint(args...))
	}
}

// Errorf uses fmt.Sprintf to logger a templated message.
func (la *Adapter) Errorf(msg string, args ...interface{}) {
	if la.Enabled(ErrorLevel) {
		la.logger.Log(ErrorLevel, fmt.Sprintf(msg, args...))
	}
}

// Errorw uses key,value to construct and logger a message.
func (la *Adapter) Errorw(keyvals ...interface{}) {
	if la.Enabled(ErrorLevel) {
		la.logger.Log(ErrorLevel, "", keyvals...)
	}

}

// Panic uses fmt.Sprint to construct and logger a message, then panics.
func (la *Adapter) Panic(args ...interface{}) {
	if la.Enabled(PanicLevel) {
		la.logger.Log(PanicLevel, sprint(args...))
	}
}

// Panicf uses fmt.Sprintf to logger a templated message, then panics.
func (la *Adapter) Panicf(msg string, args ...interface{}) {
	if la.Enabled(PanicLevel) {
		la.logger.Log(PanicLevel, fmt.Sprintf(msg, args...))
	}
}

// Panicw uses key,value to construct and logger a message.
func (la *Adapter) Panicw(keyvals ...interface{}) {
	if la.Enabled(PanicLevel) {
		la.logger.Log(PanicLevel, "", keyvals...)
	}
}

// Fatal uses fmt.Sprint to construct and logger a message, then calls os.Exit.
func (la *Adapter) Fatal(args ...interface{}) {
	if la.Enabled(FatalLevel) {
		la.logger.Log(FatalLevel, sprint(args...))
	}
	os.Exit(1)
}

// Fatalf uses fmt.Sprintf to logger a templated message, then calls os.Exit.
func (la *Adapter) Fatalf(msg string, args ...interface{}) {
	if la.Enabled(FatalLevel) {
		la.logger.Log(FatalLevel, fmt.Sprintf(msg, args...))
	}
	os.Exit(1)
}

// Fatalw uses key,value to construct and logger a message.
func (la *Adapter) Fatalw(keyvals ...interface{}) {
	if la.Enabled(FatalLevel) {
		la.logger.Log(FatalLevel, "", keyvals...)
	}
	os.Exit(1)

}

// Flush flushing any buffered log entries. Applications should take care to call Sync before exiting.
func (la *Adapter) Flush() error {
	return la.logger.Close()
}
