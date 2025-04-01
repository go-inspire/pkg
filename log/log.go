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

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level 定义了日志级别。
// 它使用 zapcore.Level 作为底层实现，提供了标准的日志级别支持。
type Level = zapcore.Level

// AtomicLevel 是一个原子的日志级别。
// 它允许在运行时动态地调整日志级别，而无需重启程序。
type AtomicLevel = zap.AtomicLevel

// 预定义的日志级别常量。
// 这些常量定义了从低到高的不同日志级别，用于控制日志输出的详细程度。
const (
	// DebugLevel 用于调试信息，通常在开发环境中启用
	DebugLevel = zapcore.DebugLevel
	// InfoLevel 用于一般信息，是默认的日志级别
	InfoLevel = zapcore.InfoLevel
	// WarnLevel 用于警告信息，表示可能的问题
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel 用于错误信息，表示程序遇到了错误但可以继续运行
	ErrorLevel = zapcore.ErrorLevel
	// PanicLevel 用于严重错误，记录日志后会触发 panic
	PanicLevel = zapcore.PanicLevel
	// FatalLevel 用于致命错误，记录日志后会导致程序退出
	FatalLevel = zapcore.FatalLevel
)

// Logger 定义了日志记录器的接口。
// 它提供了基本的日志记录功能，包括上下文支持和资源清理。
type Logger interface {
	// Log 记录一条日志消息。
	// ctx: 上下文信息，可用于传递请求级别的数据
	// level: 日志级别，决定消息的重要程度
	// msg: 日志消息内容
	// keyValues: 键值对形式的结构化数据
	Log(ctx context.Context, level Level, msg string, keyValues ...interface{})

	// Close 关闭日志记录器并释放相关资源。
	// 在程序退出前应该调用此方法以确保所有日志都被正确写入。
	Close() error
}
