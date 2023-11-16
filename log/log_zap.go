/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type ZapConfig struct {
	zap.Config

	Named map[string]Level `json:"named" yaml:"named"`
}

var _ Logger = (*zapLogger)(nil)

// zapLogger zap.Logger 的实现
type zapLogger struct {
	log *zap.Logger
}

func newZapLogger(cfg ZapConfig) *zapLogger {
	logger, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		fmt.Printf("zap.Config.Build[%v] fail\n", cfg)
		return nil
	}

	return &zapLogger{
		log: logger,
	}
}

func (l *zapLogger) Log(_ context.Context, level Level, msg string, keyValues ...interface{}) {
	var fields []zap.Field
	var f zap.Field
	for len(keyValues) > 0 {
		f, keyValues = keyValuesToField(keyValues)
		fields = append(fields, f)
	}

	switch level {
	case DebugLevel:
		l.log.Debug(msg, fields...)
	case InfoLevel:
		l.log.Info(msg, fields...)
	case WarnLevel:
		l.log.Warn(msg, fields...)
	case PanicLevel:
		l.log.Panic(msg, fields...)
	case ErrorLevel:
		l.log.Error(msg, fields...)
	case FatalLevel:
		l.log.Fatal(msg, fields...)
	}
}

func (l *zapLogger) Close() error {
	return l.log.Sync()
}

const badKey = "!BADKEY"

func keyValuesToField(args []interface{}) (zap.Field, []interface{}) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return zap.String(badKey, x), nil
		}
		return zap.Any(x, args[1]), args[2:]

	case zap.Field:
		return x, args[1:]

	default:
		return zap.Any(badKey, x), args[1:]
	}
}

// 1. 读 zap.config.json 配置文件
// 2. 读配置失败后读环境变量配置 zapLogger
func initZapLogger(zapConfig string, lvl Level) Logger {
	var logger *zapLogger

	defer func() {
		if logger != nil {
			zap.RedirectStdLog(logger.log)
		}
	}()

	file, err := filepath.Abs(zapConfig)
	if err == nil {
		logger, err = buildFrom(file)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			go func(path string) {
				err := watch(path)
				if err != nil {
					fmt.Println("error:", err)
				}
			}(file)
		}
		return logger

	} //允许配置不存在或者配置错误

	config := zap.NewProductionConfig()
	config.Level.SetLevel(lvl)
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	debug := os.Getenv("DEBUG")
	if debug != "" {
		config.Level.SetLevel(zapcore.DebugLevel)
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logger = newZapLogger(ZapConfig{Config: config})
	return logger
}

func watch(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			fmt.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
				fmt.Print("modified file:", event.Name)
				logger, err := buildFrom(path)
				if err != nil {
					fmt.Println("error:", err)
				} else {
					SetDefaultLogger(logger)
				}

			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
			fmt.Println("error:", err)
		}
	}

}

func buildFrom(file string) (*zapLogger, error) {
	config := ZapConfig{
		Config: zap.NewProductionConfig(),
		Named:  make(map[string]Level),
	}
	fmt.Printf("Loading config: %v\n", file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}

	logger := newZapLogger(config)
	SetConfig(Config{
		DefaultLevel: config.Level.Level(),
		Named:        config.Named,
	})
	return logger, nil
}
