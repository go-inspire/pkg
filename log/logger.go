package log

import (
	"fmt"
	"os"
)

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
)

var logger *zap.Logger

var sugar *zap.SugaredLogger

// 1. 读 zap.config.json 配置文件
// 2. 读配置失败后读环境变量配置 logger
func init() {
	callerSkip := zap.AddCallerSkip(1)
	file, err := filepath.Abs("zap.config.json")
	if err == nil {
		if f, err := os.Open(file); err == nil {
			cfg := zap.NewProductionConfig()
			if err := json.NewDecoder(f).Decode(&cfg); err == nil {
				logger, err = cfg.Build(callerSkip)
				if err != nil {
					fmt.Printf("zap.Config.Build[%v] fail", cfg)
				} else {
					defer logger.Sync()
					sugar = logger.Sugar()
				}
			}
		}
	} //允许配置不存在或者配置错误

	if logger != nil { //说明从配置初始化日志成功
		return
	}

	level := zapcore.InfoLevel
	if val := os.Getenv("LOGLV"); val != "" {
		if err := level.UnmarshalText([]byte(val)); err != nil {
			fmt.Printf("parse %s to zapcore.Level fail", val)
		}
	}
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(os.Stdout), level)

	debug := os.Getenv("DEBUG")
	if debug != "" {
		cfg = zap.NewDevelopmentEncoderConfig()
		level = zapcore.DebugLevel
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.Lock(os.Stdout), level)
	}

	if val := os.Getenv("LOGFILE"); val != "" {
		logfile, _ := filepath.Abs(val)
		file, _ := os.Create(logfile)

		core = zapcore.NewTee(
			core,
			zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(file), level),
		)
	}

	logger = zap.New(core, callerSkip)
	defer logger.Sync()

	sugar = logger.Sugar()

}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(msg string, args ...interface{}) {
	sugar.Debugf(msg, args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	sugar.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(msg string, args ...interface{}) {
	sugar.Infof(msg, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(msg string, args ...interface{}) {
	sugar.Warnf(msg, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	sugar.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(msg string, args ...interface{}) {
	sugar.Errorf(msg, args...)
}

// Print uses fmt.Sprint to construct and log a message.
func Print(args ...interface{}) {
	sugar.Info(args...)
}

// Printf uses fmt.Sprintf to log a templated message.
func Printf(msg string, args ...interface{}) {
	sugar.Infof(msg, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(msg string, args ...interface{}) {
	sugar.Panicf(msg, args...)
}
