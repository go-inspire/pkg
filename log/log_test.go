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
}

func TestLevel(t *testing.T) {
	atom := zap.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()

	logger.Info("info logging enabled")

	atom.SetLevel(zap.ErrorLevel)
	logger.Info("info logging disabled")
	// Output:
	// {"level":"info","msg":"info logging enabled"}
}

func TestSetLevel(t *testing.T) {

	log1 := std
	defer log1.Flush()
	dlog := std.Named("debug").SetLevel(DebugLevel)
	defer dlog.Flush()
	ilog := std.Named("info").SetLevel(InfoLevel)
	defer ilog.Flush()
	wlog := std.Named("warn").SetLevel(WarnLevel)
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
