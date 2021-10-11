package log

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-inspire/pkg/log/elog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	zap.Config

	Named map[string]Level `json:"named" yaml:"named"`
}

func (c Config) clone() Config {
	copied := c
	copied.Level = zap.NewAtomicLevelAt(c.Level.Level())
	return copied
}

// zapLogger zap.Logger 的实现
type zapLogger struct {
	cfg    Config
	level  zap.AtomicLevel
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func newZapLogger(cfg Config) *zapLogger {
	logger, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		fmt.Printf("zap.Config.Build[%v] fail\n", cfg)
		return nil
	}

	return &zapLogger{
		cfg:    cfg,
		level:  cfg.Level,
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

// Debug uses fmt.Sprint to construct and logger a message.
func (l *zapLogger) Debug(args ...interface{}) {
	if l.level.Enabled(DebugLevel) {
		l.logger.Debug(sprint(args...))
	}
}

// Debugf uses fmt.Sprintf to logger a templated message.
func (l *zapLogger) Debugf(msg string, args ...interface{}) {
	l.sugar.Debugf(msg, args...)
}

// Info uses fmt.Sprint to construct and logger a message.
func (l *zapLogger) Info(args ...interface{}) {
	if l.level.Enabled(InfoLevel) {
		l.logger.Info(sprint(args...))
	}
}

// Infof uses fmt.Sprintf to logger a templated message.
func (l *zapLogger) Infof(msg string, args ...interface{}) {
	l.sugar.Infof(msg, args...)
}

// Info uses fmt.Sprint to construct and logger a message.
func (l *zapLogger) Print(args ...interface{}) {
	if l.level.Enabled(InfoLevel) {
		l.logger.Info(sprint(args...))
	}
}

// Infof uses fmt.Sprintf to logger a templated message.
func (l *zapLogger) Printf(msg string, args ...interface{}) {
	l.sugar.Infof(msg, args...)
}

// Warn uses fmt.Sprint to construct and logger a message.
func (l *zapLogger) Warn(args ...interface{}) {
	if l.level.Enabled(WarnLevel) {
		l.logger.Warn(sprint(args...))
	}
}

// Warnf uses fmt.Sprintf to logger a templated message.
func (l *zapLogger) Warnf(msg string, args ...interface{}) {
	l.sugar.Warnf(msg, args...)
}

// Error uses fmt.Sprint to construct and logger a message.
func (l *zapLogger) Error(args ...interface{}) {
	if l.level.Enabled(ErrorLevel) {
		l.logger.Error(sprint(args...))
	}
}

// Errorf uses fmt.Sprintf to logger a templated message.
func (l *zapLogger) Errorf(msg string, args ...interface{}) {
	l.sugar.Errorf(msg, args...)
}

// Panic uses fmt.Sprint to construct and logger a message, then panics.
func (l *zapLogger) Panic(args ...interface{}) {
	if l.level.Enabled(PanicLevel) {
		l.logger.Panic(sprint(args...))
	}
}

// Panicf uses fmt.Sprintf to logger a templated message, then panics.
func (l *zapLogger) Panicf(msg string, args ...interface{}) {
	l.sugar.Panicf(msg, args...)
}

// Fatal uses fmt.Sprint to construct and logger a message, then calls os.Exit.
func (l *zapLogger) Fatal(args ...interface{}) {
	if l.level.Enabled(FatalLevel) {
		l.logger.Fatal(sprint(args...))
	}
}

// Fatalf uses fmt.Sprintf to logger a templated message, then calls os.Exit.
func (l *zapLogger) Fatalf(msg string, args ...interface{}) {
	l.sugar.Fatalf(msg, args...)
}

// Flush flushing any buffered log entries. Applications should take care to call Sync before exiting.
func (l *zapLogger) Flush() error {
	return l.logger.Sync()
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (l *zapLogger) Named(s string) Logger {
	s = strings.ToLower(s)
	cfg := l.cfg.clone()
	lvl, ok := cfg.Named[s]
	if ok {
		cfg.Level.SetLevel(lvl)
	}

	return newZapLogger(cfg)
}

//SetLevel alters the logging level.
func (l *zapLogger) SetLevel(lvl Level) Logger {
	l.level.SetLevel(lvl)
	return l
}

func sprint(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	} else if s, ok := a[0].(string); ok && len(a) == 1 {
		return s
	} else if v := reflect.ValueOf(a[0]); len(a) == 1 && v.Kind() == reflect.String {
		return v.String()
	} else {
		return fmt.Sprint(a...)
	}
}

//func sprintf(format string, a ...interface{}) string {
//	if len(a) == 0 {
//		return format
//	} else {
//		return fmt.Sprintf(format, a...)
//	}
//}

// 1. 读 zap.config.json 配置文件
// 2. 读配置失败后读环境变量配置 zapLogger
func init() {
	zapConfig := "zap.config.json"
	if val, ok := os.LookupEnv("LOG_ZAP_CONFIG"); ok {
		zapConfig = val
	}

	var logger *zapLogger

	defer func() {
		if logger != nil {
			elog.Init(logger.logger)
		}
	}()
	file, err := filepath.Abs(zapConfig)
	if err == nil {
		logger, err = load(file)
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

	} //允许配置不存在或者配置错误

	if logger != nil { //说明从配置初始化日志成功
		std = logger
		return
	}

	cfg := zap.NewProductionConfig()
	if val := os.Getenv("LOG_LEVEL"); len(val) > 0 {
		if err := cfg.Level.UnmarshalText([]byte(val)); err != nil {
			fmt.Printf("parse %s to zapcore.Level fail\n", val)
		}
	}
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	debug := os.Getenv("DEBUG")
	if debug != "" {
		cfg.Level.SetLevel(zapcore.DebugLevel)
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	if val := os.Getenv("LOG_FILE"); val != "" {
		logfile, _ := filepath.Abs(val)
		cfg.OutputPaths = append(cfg.OutputPaths, logfile)
	}

	std = newZapLogger(Config{Config: cfg})
}

func watch(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()
	defer std.Flush()

	err = watcher.Add(path)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			fmt.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Print("modified file:", event.Name)
				logger, err := load(path)
				if err != nil {
					fmt.Println("error:", err)
				} else {
					std = logger
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

func load(file string) (*zapLogger, error) {
	cfg := Config{
		Config: zap.NewProductionConfig(),
		Named:  make(map[string]Level),
	}
	fmt.Printf("Loading config: %v\n", file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if len(cfg.Named) > 0 {
		dlvl := cfg.Level.Level()
		for k, l := range loggers {
			lvl, ok := cfg.Named[k]
			if !ok {
				l.SetLevel(dlvl)
			} else {
				l.SetLevel(lvl)
			}

			loggers[k] = l
		}
	}

	logger := newZapLogger(cfg)
	return logger, nil
}
