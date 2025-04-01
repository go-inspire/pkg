/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

// Package log 提供了一个灵活的日志系统，支持多种日志后端和分层的日志级别。
// 它同时支持结构化和非结构化的日志记录功能。
package log

import (
	"context"
	"os"
	"sync/atomic"
)

// Option 是一个函数类型，用于配置 Adapter。
// 它提供了一种灵活的方式来设置 Adapter 的各种选项。
type Option func(*Adapter)

// Adapter 提供了一个高级的日志记录接口，支持不同级别和格式的日志。
// 它封装了一个 Logger 实现，并提供了便捷的方法来记录不同级别的日志。
type Adapter struct {
	logger Logger          // 底层的日志记录器
	lvl    atomic.Int32    // 原子操作的日志级别
	ctx    context.Context // 默认上下文
}

// WithLevel 返回一个 Option，用于设置 Adapter 的最低启用日志级别。
// 只有等于或高于此级别的消息才会被记录。
func WithLevel(level Level) Option {
	return func(la *Adapter) {
		la.SetLevel(level)
	}
}

// WithContext 返回一个 Option，用于设置日志适配器的上下文。
// 除非被覆盖，否则此上下文将用于所有日志操作。
func WithContext(ctx context.Context) Option {
	return func(la *Adapter) {
		la.ctx = ctx
	}
}

// NewAdapter 使用给定的日志记录器和选项创建一个新的日志适配器。
// 它使用默认设置初始化适配器，并应用所提供的选项。
func NewAdapter(logger Logger, opts ...Option) *Adapter {
	la := &Adapter{
		logger: logger,
		ctx:    context.Background(),
	}
	la.SetLevel(InfoLevel)

	if opts != nil {
		for _, o := range opts {
			o(la)
		}
	}
	return la
}

// Level 返回适配器的最低启用日志级别。
// 只有等于或高于此级别的消息才会被记录。
func (la *Adapter) Level() Level {
	return Level(int(la.lvl.Load()))
}

// SetLevel 更改适配器的日志级别。
// 它返回适配器本身，以支持方法链式调用。
func (la *Adapter) SetLevel(l Level) *Adapter {
	la.lvl.Store(int32(l))
	return la
}

// Enabled 实现了 zapcore.LevelEnabler 接口。
// 如果给定的日志级别已启用，则返回 true。
func (la *Adapter) Enabled(l Level) bool {
	return la.Level().Enabled(l)
}

// logWithLevel 是一个通用的日志记录方法，处理所有日志级别。
// 它格式化消息并在级别启用时调用底层日志记录器。
func (la *Adapter) logWithLevel(level Level, msg string, args ...interface{}) {
	la.output(la.ctx, level, func(ctx context.Context, l Level) {
		la.logger.Log(ctx, l, msg, args...)
	})
}

// logWithFormat 是一个用于格式化日志的通用方法。
// 它在记录之前使用 sprintf 格式化消息。
func (la *Adapter) logWithFormat(level Level, format string, args ...interface{}) {
	la.logWithLevel(level, sprintf(format, args...))
}

// logWithKeyValues 是一个用于键值对日志记录的通用方法。
// 它记录键值对而不带消息。
func (la *Adapter) logWithKeyValues(level Level, keyvals ...interface{}) {
	la.logWithLevel(level, "", keyvals...)
}

// Print 系列方法在 Info 级别记录日志。
// Print 记录一个简单的消息。
func (la *Adapter) Print(args ...interface{}) {
	la.logWithLevel(InfoLevel, sprint(args...))
}

// Printf 在 Info 级别记录格式化的消息。
func (la *Adapter) Printf(msg string, args ...interface{}) {
	la.logWithFormat(InfoLevel, msg, args...)
}

// Printw 在 Info 级别记录键值对。
func (la *Adapter) Printw(keyvals ...interface{}) {
	la.logWithKeyValues(InfoLevel, keyvals...)
}

// Debug 系列方法在 Debug 级别记录日志。
// Debug 记录一个简单的消息。
func (la *Adapter) Debug(args ...interface{}) {
	la.logWithLevel(DebugLevel, sprint(args...))
}

// Debugf 在 Debug 级别记录格式化的消息。
func (la *Adapter) Debugf(msg string, args ...interface{}) {
	la.logWithFormat(DebugLevel, msg, args...)
}

// Debugw 在 Debug 级别记录键值对。
func (la *Adapter) Debugw(keyvals ...interface{}) {
	la.logWithKeyValues(DebugLevel, keyvals...)
}

// Info 系列方法在 Info 级别记录日志。
// Info 记录一个简单的消息。
func (la *Adapter) Info(args ...interface{}) {
	la.logWithLevel(InfoLevel, sprint(args...))
}

// Infof 在 Info 级别记录格式化的消息。
func (la *Adapter) Infof(msg string, args ...interface{}) {
	la.logWithFormat(InfoLevel, msg, args...)
}

// Infow 在 Info 级别记录键值对。
func (la *Adapter) Infow(keyvals ...interface{}) {
	la.logWithKeyValues(InfoLevel, keyvals...)
}

// Warn 系列方法在 Warn 级别记录日志。
// Warn 记录一个简单的消息。
func (la *Adapter) Warn(args ...interface{}) {
	la.logWithLevel(WarnLevel, sprint(args...))
}

// Warnf 在 Warn 级别记录格式化的消息。
func (la *Adapter) Warnf(msg string, args ...interface{}) {
	la.logWithFormat(WarnLevel, msg, args...)
}

// Warnw 在 Warn 级别记录键值对。
func (la *Adapter) Warnw(keyvals ...interface{}) {
	la.logWithKeyValues(WarnLevel, keyvals...)
}

// Error 系列方法在 Error 级别记录日志。
// Error 记录一个简单的消息。
func (la *Adapter) Error(args ...interface{}) {
	la.logWithLevel(ErrorLevel, sprint(args...))
}

// Errorf 在 Error 级别记录格式化的消息。
func (la *Adapter) Errorf(msg string, args ...interface{}) {
	la.logWithFormat(ErrorLevel, msg, args...)
}

// Errorw 在 Error 级别记录键值对。
func (la *Adapter) Errorw(keyvals ...interface{}) {
	la.logWithKeyValues(ErrorLevel, keyvals...)
}

// Panic 系列方法在 Panic 级别记录日志，然后触发 panic。
// Panic 记录一个简单的消息，然后使用相同的消息触发 panic。
func (la *Adapter) Panic(args ...interface{}) {
	msg := sprint(args...)
	la.logWithLevel(PanicLevel, msg)
	panic(msg)
}

// Panicf 在 Panic 级别记录格式化的消息，然后触发 panic。
func (la *Adapter) Panicf(msg string, args ...interface{}) {
	formatted := sprintf(msg, args...)
	la.logWithLevel(PanicLevel, formatted)
	panic(formatted)
}

// Panicw 在 Panic 级别记录键值对，然后触发 panic。
func (la *Adapter) Panicw(keyvals ...interface{}) {
	msg := sprint(keyvals...)
	la.logWithLevel(PanicLevel, "", keyvals...)
	panic(msg)
}

// Fatal 系列方法在 Fatal 级别记录日志，然后退出程序。
// Fatal 记录一个简单的消息，然后以状态码 1 退出程序。
func (la *Adapter) Fatal(args ...interface{}) {
	msg := sprint(args...)
	la.logWithLevel(FatalLevel, msg)
	la.Flush()
	os.Exit(1)
}

// Fatalf 在 Fatal 级别记录格式化的消息，然后退出程序。
func (la *Adapter) Fatalf(msg string, args ...interface{}) {
	formatted := sprintf(msg, args...)
	la.logWithLevel(FatalLevel, formatted)
	la.Flush()
	os.Exit(1)
}

// Fatalw 在 Fatal 级别记录键值对，然后退出程序。
func (la *Adapter) Fatalw(keyvals ...interface{}) {
	la.logWithLevel(FatalLevel, "", keyvals...)
	la.Flush()
	os.Exit(1)
}

// Flush 刷新所有缓冲的日志条目。
// 它调用底层日志记录器的 Close() 方法。
func (la *Adapter) Flush() error {
	return la.Close()
}

// Close 关闭底层的日志记录器。
// 当不再需要日志记录器时应该调用此方法。
func (la *Adapter) Close() error {
	return la.logger.Close()
}

// Log 使用指定的上下文在指定的级别记录消息。
// 如果没有提供上下文，则使用适配器的默认上下文。
func (la *Adapter) Log(ctx context.Context, level Level, msg string, keyValues ...interface{}) {
	if ctx == nil {
		ctx = la.ctx
	}
	la.logger.Log(ctx, level, msg, keyValues...)
}

// output 处理实际的日志输出，并进行级别检查。
// 只有在级别启用时才调用提供的日志函数。
func (la *Adapter) output(ctx context.Context, l Level, log func(ctx context.Context, level Level)) {
	if la.Enabled(l) {
		log(ctx, l)
	}
}
