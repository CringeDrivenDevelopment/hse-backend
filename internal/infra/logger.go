package infra

import (
	"fmt"
	"time"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Zap *zap.Logger
	*zap.SugaredLogger
	Name string
}

func NewLogger(cfg *Config) (*Logger, error) {
	var l Logger
	l.Name = "main"

	var config zap.Config
	if cfg.Debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(time.FixedZone("GMT+0", 3*60*60)).Format("2006-01-02 15:04:05"))
		}
	} else {
		config = zap.NewProductionConfig()
	}

	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	l.SugaredLogger = log.Sugar()
	l.Zap = log
	return &l, nil
}

type ZapFxLogger struct{ *zap.Logger }

func (z *ZapFxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		z.Info("fx on start executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName))

	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			z.Error("fx on start failed",
				zap.String("callee", e.FunctionName),
				zap.Error(e.Err))
		} else {
			z.Info("fx on start succeeded",
				zap.String("callee", e.FunctionName),
				zap.Duration("runtime", e.Runtime))
		}

	default:
		z.Debug("fx event", zap.String("type", fmt.Sprintf("%T", e)))
	}
}

type ZapGooseAdapter struct{ zap *zap.Logger }

func (a *ZapGooseAdapter) Printf(format string, args ...any) {
	a.zap.Sugar().Infof(format, args...)
}

func (a *ZapGooseAdapter) Fatalf(format string, args ...interface{}) {
	a.zap.Sugar().Fatalf(format, args...)
}
