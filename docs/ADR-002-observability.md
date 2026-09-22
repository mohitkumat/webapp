# ADR-002: Observability — stdlib-only, not OpenTelemetry-Go / prometheus/client_golang

- **Status:** Accepted
- **Date:**

## Context
Phase F needs structured logs, request tracing, and metrics. The obvious
choices are `go.opentelemetry.io/otel` for tracing and
`github.com/prometheus/client_golang` for metrics.

## Decision
Use Go's standard library only: `log/slog` for structured logs, a small
custom request-ID + timing middleware for tracing, and a hand-rolled
`/metrics` endpoint that emits the Prometheus text exposition format.

## Why
Both real libraries were tested directly in this environment and fail to
resolve:
- `go.opentelemetry.io/otel` — the module's own vanity import domain
  (`go.opentelemetry.io`) isn't reachable.
- `github.com/prometheus/client_golang` — resolves at the top level, but its
  transitive deps require `golang.org/x/sys`, `google.golang.org/protobuf`,
  neither reachable here.

This isn't just a sandbox workaround, though — it's a reasonable default for
a *template*:
- `log/slog` is in the standard library since Go 1.21. Zero dependencies,
  zero version drift risk.
- A hand-rolled `/metrics` endpoint that speaks the Prometheus text format is
  still scrapeable by a real Prometheus server or Grafana Agent — nothing
  downstream needs to know it wasn't built with the official client.
- Request-ID-based log correlation (matching Day 9 of the learning path) is
  enough to trace a request end-to-end locally without needing a full
  distributed tracing backend (Jaeger/Tempo) that this template doesn't run.

## Consequences
- No trace spans exportable to Jaeger/Tempo/Honeycomb out of the box — you
  get correlated structured logs, not a trace UI.
- No histogram/percentile metrics (p50/p95/p99) — only counters and simple
  gauges. Good enough for local dev and small deployments; not a production
  SLO dashboard.
- If your real environment isn't network-restricted like this sandbox, swap
  in `otel-go` and `client_golang` — the seams are contained:
  `internal/observability/logger.go` (→ OTel logs/traces) and
  `internal/observability/metrics.go` (→ prometheus.Registry + promhttp).

## Revisit Trigger
Once the project needs real percentile latency metrics or exportable traces
across services, replace the hand-rolled metrics/tracing with the real
libraries — the swap only touches `internal/observability/`.
