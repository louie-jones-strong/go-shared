package telemetry

import (
	"context"
	"reflect"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	End()
	RecordError(err error)
}

type spanImpl struct {
	span trace.Span
}

func (s *spanImpl) End() {
	s.span.End()
}

func (s *spanImpl) RecordError(err error) {
	if err == nil {
		return
	}

	s.span.RecordError(err)
	s.span.SetStatus(codes.Error, err.Error())
}

func StartTrace(ctx context.Context, name string) (context.Context, Span) {

	tracer := otel.Tracer("")

	// Start a root trace span for this HTTP request
	subCtx, span := tracer.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindInternal),
		// trace.WithAttributes(attribute.String("http.method", r.Method)),
	)

	return subCtx, &spanImpl{span: span}
}

// GetFuncName extracts the clean function name using reflection and runtime metadata.
func GetFuncName(fn any) string {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		return "unknown.func"
	}

	pc := val.Pointer()
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown.func"
	}

	// f.Name() returns something like "main.ProcessOrder" or "your-module/service.Checkout"
	fullName := f.Name()

	// Optional: strip package path and keep just the function name
	// parts := strings.Split(fullName, "/")
	// shortName := parts[len(parts)-1]

	return fullName
}

type SpanFunc0 func(ctx context.Context) error
type SpanFunc1[T any] func(ctx context.Context) T

func WrapWithSpan(ctx context.Context, fn SpanFunc0) error {

	funcName := GetFuncName(fn)
	ctx, span := StartTrace(ctx, funcName)
	defer span.End()

	err := fn(ctx)

	if err != nil {
		span.RecordError(err)
	}

	return err
}
