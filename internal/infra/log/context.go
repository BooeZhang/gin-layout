package log

import (
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type requestIDKey struct{}

const (
	FieldRequestID = "request_id"
	FieldTraceID   = "trace_id"
	FieldSpanID    = "span_id"
)

type contextHook struct{}

func (contextHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	ctx := e.GetCtx()
	if ctx == nil {
		return
	}
	if requestID := RequestID(ctx); requestID != "" {
		e.Str(FieldRequestID, requestID)
	}
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		e.Str(FieldTraceID, sc.TraceID().String()).
			Str(FieldSpanID, sc.SpanID().String())
	}
}

func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}

func Ctx(ctx context.Context) *zerolog.Logger {
	if ctx == nil {
		ctx = context.Background()
	}
	l := zerolog.Ctx(ctx).With().Ctx(ctx).Logger()
	return &l
}
