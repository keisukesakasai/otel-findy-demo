package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	opentelemetry "app/internal"
)

var (
	tracer = otel.Tracer("otel-findy-demo")
	meter  = otel.Meter("otel-findy-demo")

	histogram metric.Float64Histogram
)

func main() {
	// OpenTelemetry Traces
	tracerProvider, err := opentelemetry.InitTracer()
	if err != nil {
		log.Fatalf("Error setting up trace provider: %v", err)
	}
	defer func() { _ = tracerProvider.Shutdown(context.Background()) }()

	// OpenTelemetry Metrics
	meterProvider, err := opentelemetry.InitMeter()
	if err != nil {
		log.Fatalf("Error setting up metrics provider: %v", err)
	}
	defer func() { _ = meterProvider.Shutdown(context.Background()) }()

	histogram, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("A histogram of the HTTP request durations in seconds."),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Fatalf("failed to initialize histogram: %v", err)
	}

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(mainHandler), "/")
	http.Handle("/", otelHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := tracer.Start(ctx, "main handler")
	defer span.End()

	startTime := time.Now()
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	prosessing(ctx)

	duration := time.Since(startTime)
	histogram.Record(
		ctx,
		float64(duration.Seconds()),
	)
}

func prosessing(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "processing...")
	defer span.End()

	if rand.Float64() < 1.0/100.0 {
		funcAbnormal(ctx)
	} else {
		funcNormal(ctx)
	}
}

func funcNormal(ctx context.Context) {
	_, span := tracer.Start(ctx, "funcNormal")
	defer span.End()
	time.Sleep(10 * time.Millisecond)
}

func funcAbnormal(ctx context.Context) {
	_, span := tracer.Start(ctx, "funcAbNormal(Oh...taking a lot of time...)")
	defer span.End()
	time.Sleep(3 * time.Second)
}
