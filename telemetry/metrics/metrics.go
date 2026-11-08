package metrics

import (
	"github.com/louie-jones-strong/go-shared/telemetry/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type MetricScope struct {
	name string
	m    metric.Meter
	err  error
}

func NewMetricScope(name string) (*MetricScope, func() error) {
	ms := &MetricScope{
		name: name,
		m:    otel.Meter(name),
		err:  nil,
	}
	return ms, func() error {
		return ms.err
	}
}

func (ms *MetricScope) Histogram(name string, description string, unit string) metric.Float64Histogram {
	h, err := ms.m.Float64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		logger.Error("error creating histogram %q in scope %q %v", name, ms.name, err)
		ms.err = err
	}
	return h
}

func (ms *MetricScope) Counter(name string, description string, unit string) metric.Float64Counter {
	h, err := ms.m.Float64Counter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		logger.Error("error creating counter %q in scope %q %v", name, ms.name, err)
		ms.err = err
	}
	return h
}

func (ms *MetricScope) UpDownCounter(name string, description string, unit string) metric.Float64UpDownCounter {
	h, err := ms.m.Float64UpDownCounter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		logger.Error("error creating updowncounter %q in scope %q %v", name, ms.name, err)
		ms.err = err
	}
	return h
}

func (ms *MetricScope) Gauge(name string, description string, unit string) metric.Float64Gauge {
	h, err := ms.m.Float64Gauge(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		logger.Error("error creating gauge %q in scope %q %v", name, ms.name, err)
		ms.err = err
	}
	return h
}
