/*
 * Copyright 2022 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func TestCrit(t *testing.T) {
	Debugf("test %d", 1)
	Infof("test")
	Print("test", "a")
	Printw("test", "a")
}

func TestLevel(t *testing.T) {
	atom := zap.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()

	logger.Info("", zap.Any("msg", "info logging enabled"))

	atom.SetLevel(zap.ErrorLevel)
	logger.Info("info logging disabled")
	// Output:
	// {"level":"info","msg":"info logging enabled"}
}

func TestSetLevel(t *testing.T) {

	log1 := defaultAdapter
	defer log1.Flush()
	dlog := Named("debug").SetLevel(DebugLevel)
	defer dlog.Flush()
	ilog := Named("info").SetLevel(InfoLevel)
	defer ilog.Flush()
	wlog := Named("warn").SetLevel(WarnLevel)
	defer wlog.Flush()

	log1.Debug("Debug logging enabled")
	dlog.Debug("Debug logging enabled")
	ilog.Debug("Debug logging enabled")
	wlog.Debug("Debug logging enabled")

	log1.Info("Info logging enabled")
	dlog.Info("Info logging enabled")
	ilog.Info("Info logging enabled")
	wlog.Info("Info logging enabled")

	log1.Warn("Warn logging enabled")
	dlog.Warn("Warn logging enabled")
	ilog.Warn("Warn logging enabled")
	wlog.Warn("Warn logging enabled")

	log1.SetLevel(zap.ErrorLevel)
	dlog.SetLevel(zap.ErrorLevel)
	ilog.SetLevel(zap.ErrorLevel)
	wlog.SetLevel(zap.ErrorLevel)

	log1.Debug("Debug logging disabled")
	dlog.Debug("Debug logging disabled")
	ilog.Debug("Debug logging disabled")
	wlog.Debug("Debug logging disabled")

	log1.Info("Info logging disabled")
	dlog.Info("Info logging disabled")
	ilog.Info("Info logging disabled")
	wlog.Info("Info logging disabled")

	log1.Warn("Warn logging disabled")
	dlog.Warn("Warn logging disabled")
	ilog.Warn("Warn logging disabled")
	wlog.Warn("Warn logging disabled")
}
