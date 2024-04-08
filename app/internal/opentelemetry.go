package opentelemetry

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	ServiceNameEnv           = "SERVICE_NAME"
	OtelCollectorEndpointEnv = "OTEL_COLLECTOR_ENDPOINT"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	otelCollectorEndpoint := os.Getenv(OtelCollectorEndpointEnv)

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(otelCollectorEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	serviceName := os.Getenv(ServiceNameEnv)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Println("Initialize Trace Provider ...: ", otelCollectorEndpoint)

	return tp, nil
}

func InitMeter() (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()
	otelCollectorEndpoint := os.Getenv(OtelCollectorEndpointEnv)

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(otelCollectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	serviceName := os.Getenv(ServiceNameEnv)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var boundaries []float64
	for i := 0; i < 11; i++ {
		boundary := 0.01 * math.Pow(2, float64(i))
		boundaries = append(boundaries, boundary)
	}
	view := metric.NewView(
		metric.Instrument{Kind: metric.InstrumentKindHistogram},
		metric.Stream{Aggregation: metric.AggregationExplicitBucketHistogram{
			Boundaries: boundaries,
			NoMinMax:   false,
		}},
	)
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			exporter,
			sdkmetric.WithInterval(1*time.Minute),
		)),
		sdkmetric.WithResource(res),
		metric.WithView(view),
	)
	otel.SetMeterProvider(mp)

	log.Println("Initialize Metrics Provider ...: ", otelCollectorEndpoint)

	return mp, nil
}
