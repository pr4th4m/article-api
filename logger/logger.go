package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

func New(level zapcore.Level) Logger {

	var consoleEncoder zapcore.Encoder
	var consoleLogPriority zapcore.LevelEnabler

	devConfig := zap.NewDevelopmentEncoderConfig()
	devConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	devConfig.TimeKey = ""
	devConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	devConfig.StacktraceKey = ""
	if level <= zapcore.DebugLevel {
		devConfig.StacktraceKey = "S"
		devConfig.TimeKey = "ts"
	}
	consoleEncoder = zapcore.NewConsoleEncoder(devConfig)
	consoleLogPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(io.Writer(os.Stdout)), consoleLogPriority),
	)
	return zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}

// NewNoop No operation logger
func NewNoop() Logger {
	return zap.NewNop().Sugar()
}
