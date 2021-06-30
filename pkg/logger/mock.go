package logger

import (
	"go.uber.org/zap"
)

var _ Logger = (*Mock)(nil)

type Mock struct {
	DebugFn func(msg string, fields ...zap.Field)
	InfoFn  func(msg string, fields ...zap.Field)
	ErrorFn func(err error, fields ...zap.Field)
	FatalFn func(err error, fields ...zap.Field)
}

func (m Mock) Debug(msg string, fields ...zap.Field) {
	m.DebugFn(msg, fields...)
}

func (m Mock) Info(msg string, fields ...zap.Field) {
	m.InfoFn(msg, fields...)
}

func (m Mock) Error(err error, fields ...zap.Field) {
	m.ErrorFn(err, fields...)
}

func (m Mock) Fatal(err error, fields ...zap.Field) {
	m.FatalFn(err, fields...)
}
