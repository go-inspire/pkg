/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

// Package log 实现了一个线程安全的日志系统，支持多级别日志记录和组件隔离配置
package log

import (
	"strings"
	"sync"
)

var (
	// 全局默认日志适配器实例
	defaultAdapter *Adapter
	// 全局日志配置，包含默认级别和各组件级别
	cfg = Config{DefaultLevel: InfoLevel, Named: make(map[string]Level)}
	// 组件名称到日志适配器的映射表
	adapters = make(map[string]*Adapter)
	// 保护全局变量的互斥锁
	mu sync.Mutex
)

// Config 定义日志系统的配置结构
type Config struct {
	DefaultLevel Level            // 系统默认日志级别
	Named        map[string]Level // 各组件特定的日志级别配置
}

// findConfigLevel 实现组件日志级别的层级查找逻辑
// 例如对于组件"a.b.c"，依次查找"a.b.c"、"a.b"、"a"的配置
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

// Named 获取或创建指定组件的日志适配器
// 组件名称不区分大小写，会自动转换为小写
func Named(s string) *Adapter {
	mu.Lock()
	defer mu.Unlock()

	s = strings.ToLower(s)
	if a, ok := adapters[s]; ok {
		return a
	}

	level := findConfigLevel(&cfg, s)
	a := NewAdapter(defaultAdapter.logger, WithLevel(level))
	adapters[s] = a
	return a
}

// SetConfig 更新全局日志配置并刷新所有适配器
func SetConfig(config Config) {
	mu.Lock()
	defer mu.Unlock()

	cfg = config
	for k, a := range adapters {
		a.SetLevel(findConfigLevel(&cfg, k))
	}
}

// SetDefaultAdapter 替换默认日志适配器
func SetDefaultAdapter(adapter *Adapter) {
	mu.Lock()
	defer mu.Unlock()
	defaultAdapter = adapter
}

// SetDefaultLogger 通过Logger接口创建并设置默认适配器
func SetDefaultLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	defaultAdapter = NewAdapter(logger, WithLevel(cfg.DefaultLevel))
}

// 以下是各级别日志方法的快捷方式，均委托给defaultAdapter处理

func Debug(args ...interface{}) {
	defaultAdapter.Debug(args...)
}

func Debugf(msg string, args ...interface{}) {
	defaultAdapter.Debugf(msg, args...)
}

func Debugw(keyvals ...interface{}) {
	defaultAdapter.Debugw(keyvals...)
}

// Info 记录信息级别日志，参数通过fmt.Sprint格式化
func Info(args ...interface{}) {
	defaultAdapter.Info(args...)
}

// Infof 记录信息级别日志，使用格式化字符串和参数
func Infof(msg string, args ...interface{}) {
	defaultAdapter.Infof(msg, args...)
}

// Infow 记录带键值对的信息级别日志
func Infow(keyvals ...interface{}) {
	defaultAdapter.Infow(keyvals...)
}

// Warn 记录警告级别日志，参数通过fmt.Sprint格式化
func Warn(args ...interface{}) {
	defaultAdapter.Warn(args...)
}

// Warnf 记录警告级别日志，使用格式化字符串和参数
func Warnf(msg string, args ...interface{}) {
	defaultAdapter.Warnf(msg, args...)
}

// Warnw 记录带键值对的警告级别日志
func Warnw(keyvals ...interface{}) {
	defaultAdapter.Warnw(keyvals...)
}

// Error 记录错误级别日志，参数通过fmt.Sprint格式化
func Error(args ...interface{}) {
	defaultAdapter.Error(args...)
}

// Errorf 记录错误级别日志，使用格式化字符串和参数
func Errorf(msg string, args ...interface{}) {
	defaultAdapter.Errorf(msg, args...)
}

// Errorw 记录带键值对的错误级别日志
func Errorw(keyvals ...interface{}) {
	defaultAdapter.Errorw(keyvals...)
}

// Print 兼容标准日志接口，等同于Info级别
func Print(args ...interface{}) {
	defaultAdapter.Info(args...)
}

// Printf 兼容标准日志接口，等同于Infof级别
func Printf(msg string, args ...interface{}) {
	defaultAdapter.Infof(msg, args...)
}

// Printw 兼容标准日志接口，等同于Infow级别
func Printw(keyvals ...interface{}) {
	defaultAdapter.Printw(keyvals...)
}

// Panic 记录日志后触发panic，使用fmt.Sprint格式化消息
func Panic(args ...interface{}) {
	defaultAdapter.Panic(args...)
}

// Panicf 记录日志后触发panic，使用格式化字符串和参数
func Panicf(msg string, args ...interface{}) {
	defaultAdapter.Panicf(msg, args...)
}

// Panicw 记录带键值对的日志后触发panic
func Panicw(keyvals ...interface{}) {
	defaultAdapter.Panicw(keyvals...)
}

// Fatal 记录日志后调用os.Exit(1)，使用fmt.Sprint格式化消息
func Fatal(args ...interface{}) {
	defaultAdapter.Fatal(args...)
}

// Fatalf 记录日志后调用os.Exit(1)，使用格式化字符串和参数
func Fatalf(msg string, args ...interface{}) {
	defaultAdapter.Fatalf(msg, args...)
}

// Fatalw 记录带键值对的日志后调用os.Exit(1)
func Fatalw(keyvals ...interface{}) {
	defaultAdapter.Fatalw(keyvals...)
}

// Flush 确保所有缓冲日志条目被写入
func Flush() error {
	return defaultAdapter.Flush()
}

// SetLevel 动态设置指定组件的日志级别
// name: 组件名称
// lvl: 日志级别字符串(debug/info/warn/error/panic/fatal)
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
