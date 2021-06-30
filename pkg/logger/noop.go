package logger

import (
	"go.uber.org/zap"
)

var _ Logger = (*noop)(nil)

type noop struct {
}

func (n noop) Debug(msg string, fields ...zap.Field) {

}

func (n noop) Info(msg string, fields ...zap.Field) {

}

func (n noop) Error(err error, fields ...zap.Field) {

}

func (n noop) Fatal(err error, fields ...zap.Field) {

}

func NewNoop() Logger {
	return noop{}
}
