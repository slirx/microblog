package logger

import (
	"fmt"

	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	StackTrackField = "stack_trace"
)

var _ Logger = (*zapLogger)(nil)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(err error, fields ...zap.Field)
	Fatal(err error, fields ...zap.Field)
}

type zapLogger struct {
	Logger *zap.Logger
}

func (z zapLogger) Debug(msg string, fields ...zap.Field) {
	z.Logger.Debug(msg, fields...)
}

func (z zapLogger) Info(msg string, fields ...zap.Field) {
	z.Logger.Info(msg, fields...)
}

func (z zapLogger) Error(err error, fields ...zap.Field) {
	if err == nil {
		return
	}

	fields = append(fields, StackTraceField(err))
	z.Logger.Error(err.Error(), fields...)
}

func (z zapLogger) Fatal(err error, fields ...zap.Field) {
	if err == nil {
		return
	}

	fields = append(fields, StackTraceField(err))
	z.Logger.Fatal(err.Error(), fields...)
}

func NewZapLogger() (Logger, error) {
	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "", // caller
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zl, err := zapConfig.Build(zap.WrapCore((&apmzap.Core{}).WrapCore))
	if err != nil {
		return nil, fmt.Errorf("can not initialize zap logger: %w", err)
	}

	l := zapLogger{Logger: zl}

	return l, nil
}

func StackTraceField(err error) zap.Field {
	t := getStackTrace(0, err)
	if t == "" { // in case stack trace is empty - write it
		err = errors.WithStack(err)
		t = getStackTrace(2, err)
	}

	return zap.String(StackTrackField, t)
}
