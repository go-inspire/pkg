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

// init 初始化默认日志记录器
// 根据配置自动选择使用 zap 或 slog 作为日志实现
func init() {
	var logger Logger
	defer func() {
		if logger != nil {
			SetDefaultLogger(logger)
		}
	}()

	// 读取 zap 配置文件路径
	// 优先从环境变量 LOG_ZAP_CONFIG 获取，默认使用 zap.config.json
	zapConfig := "zap.config.json"
	if val, ok := os.LookupEnv("LOG_ZAP_CONFIG"); ok {
		zapConfig = val
	}

	// 从环境变量 LOG_LEVEL 获取日志级别
	var lvl Level
	if val := os.Getenv("LOG_LEVEL"); len(val) > 0 {
		if err := lvl.UnmarshalText([]byte(val)); err != nil {
			fmt.Printf("解析日志级别 %s 失败\n", val)
		}
	}

	// 检查 zap 配置文件是否存在
	// 如果不存在则使用 slog 作为默认日志实现
	if _, err := os.Stat(zapConfig); os.IsNotExist(err) {
		logger = initSlogLogger(lvl)
		fmt.Println("使用 slog 作为默认日志记录器")
	} else {
		logger = initZapLogger(zapConfig, lvl)
		fmt.Println("使用 zap 作为默认日志记录器")
	}
}
