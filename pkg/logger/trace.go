package logger

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func getStackTrace(skip int, err error) string {
	t, ok := err.(stackTracer)
	if !ok {
		return ""
	}

	st := t.StackTrace()

	return fmt.Sprintf("%+v", st[skip:int(math.Min(5, float64(len(st))))])
}
