package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(message string, kv ...interface{})
	Info(message string, kv ...interface{})
	Warn(message string, kv ...interface{})
	Error(message string, kv ...interface{})
	Panic(message string, kv ...interface{})
	With(kv ...interface{}) Logger
}

type loggerBuilder[D any] struct {
	*Injector[D]

	*zap.Logger

	LogLevel Config `config:"log.level,str" default:"debug" usage:"Log level (debug, info, error)"`
	LogJson  Config `config:"log.json,bool" default:"false" usage:"Log in json format"`
	LogDev   Config `config:"log.dev,bool" default:"true" usage:"Log in development mode"`
}

func (l *loggerBuilder[D]) build() Logger {

	zapConfig := zap.NewProductionConfig()

	if l.LogDev.IsSet() && l.LogDev.BoolVal() {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	if l.LogLevel.IsSet() {
		switch l.LogLevel.StringVal() {

		case "error":
			zapConfig.Level.SetLevel(zap.ErrorLevel)

		case "info":
			zapConfig.Level.SetLevel(zap.InfoLevel)

		default:
			zapConfig.Level.SetLevel(zap.DebugLevel)

		}
	}

	zapConfig.Encoding = "console"

	if l.LogJson.IsSet() && l.LogJson.BoolVal() {
		zapConfig.Encoding = "json"
	}

	logger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}

	l.Logger = logger

	return &stdLogger{logger.Sugar()}

}

type stdLogger struct {
	logger *zap.SugaredLogger
}

func (s *stdLogger) Debug(message string, kv ...interface{}) {
	s.logger.Debugw(message, kv...)
}

func (s *stdLogger) Info(message string, kv ...interface{}) {
	s.logger.Infow(message, kv...)
}

func (s *stdLogger) Warn(message string, kv ...interface{}) {
	s.logger.Warnw(message, kv...)
}

func (s *stdLogger) Error(message string, kv ...interface{}) {
	s.logger.Errorw(message, kv...)
}

func (s *stdLogger) Panic(message string, kv ...interface{}) {
	s.logger.Panicw(message, kv...)
}

func (s *stdLogger) With(kv ...interface{}) Logger {
	return &stdLogger{s.logger.With(kv...)}
}
