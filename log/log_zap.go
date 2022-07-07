package log

import (
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

func (l *zapLogger) Log(level Level, msg string, keyvals ...interface{}) {
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case DebugLevel:
		l.log.Debug(msg, data...)
	case InfoLevel:
		l.log.Info(msg, data...)
	case WarnLevel:
		l.log.Warn(msg, data...)
	case PanicLevel:
		l.log.Panic(msg, data...)
	case ErrorLevel:
		l.log.Error(msg, data...)
	case FatalLevel:
		l.log.Fatal(msg, data...)
	}
}

func (l *zapLogger) Close() error {
	return l.log.Sync()
}

// 1. 读 zap.config.json 配置文件
// 2. 读配置失败后读环境变量配置 zapLogger
func init() {
	var logger *zapLogger

	defer func() {
		if logger != nil {
			SetDefaultLogger(logger)
		}
	}()

	zapConfig := "zap.config.json"
	if val, ok := os.LookupEnv("LOG_ZAP_CONFIG"); ok {
		zapConfig = val
	}

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

	} //允许配置不存在或者配置错误

	if logger != nil { //说明从配置初始化日志成功
		return
	}

	config := zap.NewProductionConfig()
	if val := os.Getenv("LOG_LEVEL"); len(val) > 0 {
		if err := config.Level.UnmarshalText([]byte(val)); err != nil {
			fmt.Printf("parse %s to zapcore.Level fail\n", val)
		}
	}
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	debug := os.Getenv("DEBUG")
	if debug != "" {
		config.Level.SetLevel(zapcore.DebugLevel)
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	if val := os.Getenv("LOG_FILE"); val != "" {
		logfile, _ := filepath.Abs(val)
		config.OutputPaths = append(config.OutputPaths, logfile)
	}
	logger = newZapLogger(ZapConfig{Config: config})
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
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Remove == fsnotify.Remove {
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
	SetConfig(LogConfig{
		DefaultLevel: config.Level.Level(),
		Named:        config.Named,
	})
	return logger, nil
}
