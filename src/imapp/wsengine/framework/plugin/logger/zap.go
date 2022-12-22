package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"sync"
	"ws/framework/env"
)

var log *zap.SugaredLogger
var core zapcore.Core
var once sync.Once

// DefaultLogger .
func DefaultLogger() *zap.SugaredLogger {
	once.Do(func() {
		logLevel := env.NacosConfig.LogLevel

		if logLevel == "" {
			logLevel = "debug"
		}

		mapping := map[string]zapcore.Level{
			"debug": zapcore.DebugLevel,
			"info":  zapcore.InfoLevel,
			"warn":  zapcore.WarnLevel,
			"error": zapcore.ErrorLevel,
		}

		logLevelCode := mapping[logLevel]

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("01-02 15:04:05.00")
		encoderConfig.ConsoleSeparator = " | "

		if runtime.GOOS == "windows" {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), logLevelCode)

		log = zap.New(core).Sugar()
	})

	return log
}

// New .
func New(name string) *zap.SugaredLogger {
	return DefaultLogger().Named(name)
}

// EnabledDebug .
func EnabledDebug() bool {
	return core.Enabled(zap.DebugLevel)
}
