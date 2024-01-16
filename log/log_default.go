/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"fmt"
	"os"
)

// SetDefaultLogger set default logger.
func init() {
	var logger Logger
	defer func() {
		if logger != nil {
			SetDefaultLogger(logger)
		}
	}()

	//读取 zap 配置文件
	zapConfig := "zap.config.json"
	if val, ok := os.LookupEnv("LOG_ZAP_CONFIG"); ok {
		zapConfig = val
	}

	var lvl Level
	if val := os.Getenv("LOG_LEVEL"); len(val) > 0 {
		if err := lvl.UnmarshalText([]byte(val)); err != nil {
			fmt.Printf("parse %s to zapcore.Level fail\n", val)
		}
	}

	//如果 zap 配置文件不存在，则使用 slog
	if _, err := os.Stat(zapConfig); os.IsNotExist(err) {
		logger = initSlogLogger(lvl)
		fmt.Println("use slog logger as default logger")
	} else {
		logger = initZapLogger(zapConfig, lvl)
		fmt.Println("use zap logger as default logger")
	}
}
