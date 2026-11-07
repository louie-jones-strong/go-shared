// Package telemetry provides logging, development metrics and traces for applications.
package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/louie-jones-strong/go-shared/telemetry/logger"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Setup initializes the telemetry systems
func Setup(
	ctx context.Context,
	consoleLogLevel logger.Level,
	fileLogLevel logger.Level,
	logFilePath string,
	serviceName string,
	serviceVersion string,
	environment string,
	otlPEndpoint string,
	metricInterval time.Duration,
) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			semconv.DeploymentEnvironmentName(environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	conn, err := grpc.NewClient(otlPEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OTLP collector at %s: %w", otlPEndpoint, err)
	}

	err = setupLogger(ctx, serviceName, consoleLogLevel, fileLogLevel, logFilePath, conn, res)
	if err != nil {
		return nil, err
	}
	logger.Debug("Logger setup complete")

	tp, err := setupTracer(ctx, conn, res)
	if err != nil {
		return nil, fmt.Errorf("failed to setup tracer: %w", err)
	}
	logger.Debug("Tracer setup complete")

	mp, err := setupMetrics(ctx, conn, res, metricInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to setup metrics: %w", err)
	}
	logger.Debug("Metrics setup complete")

	// Define graceful shutdown handle
	shutdown := func(shutdownCtx context.Context) error {
		var errs []error
		if err := tp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("tracer shutdown error: %w", err))
		}
		if err := mp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("meter shutdown error: %w", err))
		}
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("gRPC connection close error: %w", err))
		}
		if len(errs) > 0 {
			return fmt.Errorf("telemetry shutdown errors: %v", errs)
		}
		return nil
	}

	return shutdown, nil
}

func setupLogger(
	ctx context.Context,
	serviceName string,
	consoleLogLevel logger.Level,
	fileLogLevel logger.Level,
	logFilePath string,
	conn *grpc.ClientConn,
	res *resource.Resource,
) error {
	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return fmt.Errorf("failed to create log exporter: %w", err)
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)
	otel.SetLoggerProvider(loggerProvider)

	slogHandler := otelslog.NewHandler(serviceName, otelslog.WithLoggerProvider(loggerProvider))

	err = logger.SetupLogging(consoleLogLevel, fileLogLevel, logFilePath, slogHandler)
	if err != nil {
		return err
	}
	return nil
}

func setupTracer(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*sdktrace.TracerProvider, error) {

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func setupMetrics(
	ctx context.Context,
	conn *grpc.ClientConn,
	res *resource.Resource,
	interval time.Duration,
) (*sdkmetric.MeterProvider, error) {

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(interval))),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}
