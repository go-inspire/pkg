/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"log/slog"
	"os"
)

// 确保 slogLogger 实现了 Logger 接口
var _ Logger = (*slogLogger)(nil)

// slogLogger 是基于 slog 的日志记录器实现
type slogLogger struct {
	log *slog.Logger // 底层的 slog 日志记录器
}

// newSlogLogger 创建并返回一个新的 slogLogger 实例
func newSlogLogger(h slog.Handler) *slogLogger {
	logger := slog.New(h)
	return &slogLogger{
		log: logger,
	}
}

// 定义自定义日志级别常量
const (
	levelDebug = slog.LevelDebug // 调试级别
	levelInfo  = slog.LevelInfo  // 信息级别
	levelWarn  = slog.LevelWarn  // 警告级别
	levelError = slog.LevelError // 错误级别
	levelPanic = slog.Level(10)  // 恐慌级别(自定义)
	levelFatal = slog.Level(12)  // 致命级别(自定义)
)

// customLevel 自定义日志级别的显示格式
func customLevel(groups []string, a slog.Attr) slog.Attr {
	// 只处理日志级别属性
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)

		// 根据级别设置对应的显示字符串
		switch {
		case level < levelDebug:
			a.Value = slog.StringValue("TRACE")
		case level < levelInfo:
			a.Value = slog.StringValue("DEBUG")
		case level < levelWarn:
			a.Value = slog.StringValue("INFO")
		case level < levelError:
			a.Value = slog.StringValue("WARN")
		case level < levelPanic:
			a.Value = slog.StringValue("ERROR")
		case level < levelFatal:
			a.Value = slog.StringValue("FATAL")
		default:
			a.Value = slog.StringValue("INFO")
		}
	}

	return a
}

// Log 实现 Logger 接口的 Log 方法
func (l *slogLogger) Log(ctx context.Context, level Level, msg string, keyValues ...interface{}) {
	switch level {
	case DebugLevel:
		l.log.Log(ctx, levelDebug, msg, keyValues...)
	case InfoLevel:
		l.log.Log(ctx, levelInfo, msg, keyValues...)
	case WarnLevel:
		l.log.Log(ctx, levelWarn, msg, keyValues...)
	case ErrorLevel:
		l.log.Log(ctx, levelError, msg, keyValues...)
	case PanicLevel:
		l.log.Log(ctx, levelPanic, msg, keyValues...)
	case FatalLevel:
		l.log.Log(ctx, levelFatal, msg, keyValues...)
	}
}

// Close 实现 Logger 接口的 Close 方法
func (l *slogLogger) Close() error {
	return nil // slog 不需要显式关闭
}

// initSlogLogger 初始化并返回一个 slog 日志记录器
func initSlogLogger(lvl Level) Logger {
	logger := newSlogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,    // 不添加源代码位置
		ReplaceAttr: customLevel, // 使用自定义级别格式
		Level:       toSlogLevel(lvl), // 设置日志级别
	}))

	return logger
}

// toSlogLevel 将自定义 Level 转换为 slog.Level
func toSlogLevel(l Level) slog.Level {
	switch l {
	case DebugLevel:
		return levelDebug
	case InfoLevel:
		return levelInfo
	case WarnLevel:
		return levelWarn
	case ErrorLevel:
		return levelError
	case PanicLevel:
		return levelPanic
	case FatalLevel:
		return levelFatal
	default:
		return levelInfo
	}
}
