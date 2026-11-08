package telemetry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/louie-jones-strong/go-shared/telemetry/logger"
	"github.com/louie-jones-strong/go-shared/telemetry/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// responseWriter wraps http.ResponseWriter to capture status code and bytes written.
type responseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int
	wroteHeader  bool
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

// Support optional interfaces like http.Flusher if needed by downstream handlers
func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// TelemetryMiddleware returns a Chi-compatible middleware that handles OTel, metrics, and logging.
func TelemetryMiddleware(next http.Handler) http.Handler {
	propagator := otel.GetTextMapPropagator()

	ms, msErr := metrics.NewMetricScope("http_server")
	httpDurationMeter := ms.Histogram("http.server.request.duration", "how long this request took to process", "{ms}")
	requestCounterMeter := ms.Counter("http.server.request.count", "how many times this request has been requested", "")
	if err := msErr(); err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 1. Extract incoming trace context from headers (for distributed tracing across services)
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// 2. Start root server span
		spanName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		ctx, span := otel.Tracer("").Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.target", r.URL.Path),
				attribute.String("net.peer.ip", r.RemoteAddr),
			),
		)
		defer span.End()

		// Wrap response writer to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Pass down the chain with updated context
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		// 3. Calculate metrics & duration
		durationMs := time.Since(startTime).Milliseconds()
		statusStr := fmt.Sprintf("%d", wrappedWriter.status)

		// Record Dev Metric: HTTP Request Duration Histogram
		commonAttrs := metric.WithAttributes(
			attribute.String("http.request.method", r.Method),
			attribute.String("http.response.status_code", statusStr), // Note: response.status_code
			attribute.String("http.route", r.Pattern),
		)
		httpDurationMeter.Record(ctx, float64(durationMs), commonAttrs)
		requestCounterMeter.Add(ctx, 1, commonAttrs)

		// 4. Record Structured Log with Trace context injection for Loki/Grafana bridging
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		logAttrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrappedWriter.status,
			"duration_ms", durationMs,
			"bytes_written", wrappedWriter.bytesWritten,
			"trace_id", traceID,
			"span_id", spanID,
		}

		span.SetAttributes(attribute.Int("status", wrappedWriter.status))

		logFunc := logger.InfoWithArgs
		if wrappedWriter.status >= 500 {
			logFunc = logger.ErrorWithArgs
		} else if wrappedWriter.status >= 400 {
			logFunc = logger.WarnWithArgs
		}
		logFunc(
			fmt.Sprintf(
				"%v: %s %s in %v",
				wrappedWriter.status,
				r.Method,
				r.URL.Path,
				durationMs,
			),
			logAttrs...,
		)
	})
}
