# Web APP Project Template

Quick-start monorepo: React + TypeScript frontend, Go backend, Postgres, containerized.

## Structure
```
apps/web/       React + TypeScript UI (Vite)
apps/api/       Go API + background workers
packages/       Shared code (types, config) — added as needed
docs/           PRD, ADRs, runbooks
```

## Status
Template complete — Phases A through F built and verified; Phase G folded in
(PRD template + both ADRs live in `docs/`).

- [x] A. Repo foundation
- [x] B. Frontend (Vite + React + TS)
- [x] C. Backend (Go API + Postgres)
- [x] D. Containers (Docker + Compose)
- [x] E. CI (GitHub Actions)
- [x] F. Observability (logs, traces, metrics)
- [x] G. Docs (PRD/ADR templates — see `docs/`)

**Note on Phase D verification:** this sandbox has no Docker daemon and no
access to Docker Hub, so the Dockerfiles/compose file below are hand-verified
(YAML syntax checked, image/base choices double-checked) but not run
end-to-end here. Run `docker compose up --build` on your machine as the real
check, and flag anything that doesn't come up clean.

## Setup

### Frontend (`apps/web`)
```bash
cd apps/web
npm install
npm run dev     # starts at http://127.0.0.1:3000
npm test        # runs the sample test suite
npm run build   # production build
npm run lint    # ESLint check
```
The home page links to `/health`, which calls the API's `/health` endpoint and shows
loading / success / error states — a live example to follow for new pages.

### Backend (`apps/api`)
Requires a local Postgres reachable at `DATABASE_URL` (see `apps/api/.env.example`).
```bash
cd apps/api
go run ./cmd/api    # starts at http://127.0.0.1:8080, auto-applies migrations
go test ./... -v    # runs against a real local Postgres
```
See `apps/api/README.md` for endpoints, structure, and a note on the dependency
substitution made for this sandbox (`lib/pq` instead of `pgx`, a hand-rolled
migration runner instead of `golang-migrate`).

### Full stack (Docker Compose)
```bash
docker compose up --build
```
- Web: http://127.0.0.1:3000
- API: http://127.0.0.1:8080
- Postgres: 127.0.0.1:5432 (user/pass `postgres`, db `appdb`)

All ports are published to `127.0.0.1` only — nothing is reachable from other
machines on your network. The API runs migrations automatically on startup.

**CORS:** the API only allows requests from `ALLOWED_ORIGIN` (defaults to
`http://127.0.0.1:3000`, the web app's origin). Change it via env var if you
run the web app somewhere else.

**Container networking note:** the API binds to `127.0.0.1` by default for
bare-metal `go run` (loopback-only). Inside the container it's overridden to
`HOST=0.0.0.0` via compose — otherwise Docker's published port wouldn't reach
it. The 127.0.0.1-only restriction is instead enforced by compose's
`"127.0.0.1:8080:8080"` port mapping on the host side.

### CI (GitHub Actions)
Every push/PR to `main` runs two jobs in parallel:
- **web**: `npm ci` → lint → test → build
- **api**: gofmt check → `go vet` → build → test (against a real Postgres
  service container spun up by the workflow)

A third job, **ci-ok**, only passes once both do — point branch protection at
that one job (see `docs/branch-protection.md`) so future jobs don't require
re-touching repo settings.

**Verification note:** this sandbox has no Docker, so I couldn't run the
workflow through `act`. Instead I ran every command the workflow runs
(`npm ci`, lint, test, build; `gofmt`, `go vet`, `go build`, `go test`) by
hand, in order — all pass. The one thing not verified here is the GitHub
Actions runner environment itself (service container startup, caching
behavior) — push a PR and confirm the checks actually appear.

### Observability
- **Structured logs** — every API request logs one JSON line (`request_id`,
  `method`, `path`, `status`, `duration_ms`) via `log/slog`.
- **Request correlation ("tracing")** — every request gets an `X-Request-Id`
  (reused if the caller already sent one). Grep any deployment's logs for one
  ID to see that request's full path. See `docs/ADR-002-observability.md` for
  why this is stdlib-based instead of full OpenTelemetry.
- **Metrics** — `GET /metrics` on the API serves Prometheus text format
  (`http_requests_total`, `http_requests_4xx_total`, `http_requests_5xx_total`)
  — scrapeable by a real Prometheus/Grafana Agent even though it's hand-rolled.
- **Frontend analytics** — `apps/web/src/lib/analytics.ts` is a
  provider-agnostic `track`/`identify` hook, console-logging by default.
  Swap the implementation for PostHog/Plausible later without touching call
  sites.

## Tool versions
Pinned in `.tool-versions`. After installing tools, run
`bash scripts/record-versions.sh > workstation-versions.txt` and commit the result.

### Docs
- `docs/PRD-template.md` — fill out at project kickoff, before writing code.
- `docs/ADR-001-backend-language.md` — why Go over Node.js.
- `docs/ADR-002-observability.md` — why stdlib observability over OTel/Prometheus client.
- `docs/branch-protection.md` — how to wire CI to actually block merges.

## Contributing
- One feature per branch: `feature/short-description`
- PRs must pass lint + test + build — enforced by CI (`ci-ok` check)
