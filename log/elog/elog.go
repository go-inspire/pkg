package elog

import (
	elog "github.com/ethereum/go-ethereum/log"
	"go.uber.org/zap"
)

// InitLogger init elog by zap logger.
func Init(zapLogger *zap.Logger) {
	elog.Root().SetHandler(zapHandler(zapLogger))
}

func zapHandler(log *zap.Logger) elog.Handler {
	logger := log.WithOptions(zap.AddCallerSkip(4)).Sugar()
	return elog.FuncHandler(func(r *elog.Record) error {
		args := make([]interface{}, 0, len(r.Ctx)+1)
		args = append(args, r.Msg)
		args = append(args, r.Ctx)
		switch r.Lvl {
		case elog.LvlCrit:
			logger.Fatal(args...)
		case elog.LvlError:
			logger.Error(args...)
		case elog.LvlWarn:
			logger.Warn(args...)
		case elog.LvlDebug:
		case elog.LvlTrace:
			logger.Debug(args...)
		default:
			logger.Info(args...)
		}
		return nil
	})
}
