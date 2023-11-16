/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"go.uber.org/zap"
	"testing"
)

func Test_ZapLogger(t *testing.T) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	l, _ := config.Build()

	l.Info("debug", zap.String("k1", "v1"))
}
