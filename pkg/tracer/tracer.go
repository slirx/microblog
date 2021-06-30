// tracer package contains interfaces for application monitoring. for example for integration with APM.
package tracer

import (
	"context"

	"go.elastic.co/apm"
)

// Tracer defines methods for performance monitoring, metrics.
type Tracer interface {
	RequestID(ctx context.Context) string
}

type apmTracer struct {
}

func (t apmTracer) RequestID(ctx context.Context) string {
	tx := apm.TransactionFromContext(ctx)
	traceContext := tx.TraceContext()

	return traceContext.Trace.String()
}

func NewAPMTracer() Tracer {
	return apmTracer{}
}
