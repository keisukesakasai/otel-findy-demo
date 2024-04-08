OpenTelemetry Observability運用の実例 Lunch LT で使うデモ用のリポジトリです。

https://findy.connpass.com/event/313260/

## Requirements
- Google Kubernetes Engine
- Role
  - roles/cloudtrace.agent
  - roles/monitoring.metricWriter

## /app
Web サーバーアプリケーション。リクエストを受けるとトレースと、サーバーの Duration タイムをヒストグラムメトリクスとして OTel Collector に OTLP で送信する。また、環境変数 `OTEL_GO_X_EXEMPLAR` でトレースエグザンプラーを有効化できる。

### 環境変数
- SERVICE_NAME
- OTEL_COLLECTOR_ENDPOINT
- OTEL_GO_X_EXEMPLAR
- OTEL_METRICS_EXEMPLAR_FILTER

## /deployments

### otel-findy-demo
```sh
# Create Namespace
$ kubectl create ns app

# Deploy Sample Application
$ kubectl apply -f deployments/otel-findy-demo/otel-findy-demo.yaml
```

### otelcol
```sh
# Create Namespace
$ kubectl create ns observability

# Deploy OpenTelemetryCollector
$ kubectl apply -f deployments/otelcol/otelcol.yaml
```

### OSS Observability Stack
```sh

```