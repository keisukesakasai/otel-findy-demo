FROM ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:0.97.0
COPY otelcol-config.yaml /etc/otel/otelcol-config.yaml
EXPOSE 4317
CMD ["--config", "/etc/otel/otelcol-config.yaml"]