# Monorepo template — build plan

## Phase A — Repo foundation (~2hrs)
1. Init pnpm workspace monorepo (`apps/web`, `apps/api`, `packages/*`) — 1hr
2. Pin tool versions (`.tool-versions` / `engines` field) — 30min
3. Root README with setup/start/stop/test commands — 30min

## Phase B — Frontend (~3hrs)
4. Scaffold React + TypeScript + Vite in `apps/web` — 1hr
5. Add ESLint, Prettier, Vitest — 1hr
6. Sample routed page + one component test — 1hr

## Phase C — Backend (~3-4hrs)
7. Scaffold Go API (`net/http` or chi) with `/health` — 1hr
8. Connect Postgres via pgx + a migration tool (golang-migrate) — 1-2hr
9. Add tests with Go's built-in `testing` package — 1hr

## Phase D — Containers (~3hrs)
10. Dockerfile for web, Dockerfile for api — 1hr
11. `docker-compose.yml`: web + api + postgres, ports bound to `127.0.0.1` — 1-2hr
12. Verify `docker compose up` boots the full stack — 1hr

## Phase E — CI (~2hrs)
13. GitHub Actions: lint → test → build on every PR — 1hr
14. Branch protection + required checks — 1hr

## Phase F — Observability (~3-4hrs)
15. Structured logging (pino) with request IDs — 1hr
16. Basic OpenTelemetry tracing (API → DB) — 1-2hr
17. `/metrics` endpoint (prom-client) — 1hr
18. Analytics hook placeholder (PostHog/Plausible) — 30min

## Phase G — Docs
19. Drop `PRD-template.md` into `/docs`
20. Add a one-page ADR (architecture decision record) template

---

## As built

The plan above is the original intent. The table records where the build diverged from it; steps not listed were built as planned. Status reflects the repo on 2026-09-23.

| Step | Planned | As built | Why / record |
|---|---|---|---|
| 5 | ESLint, Prettier, Vitest | As planned, but `apps/web` uses npm + `package-lock.json` inside the pnpm workspace | Not recorded |
| 8 | pgx + golang-migrate | `database/sql` + `lib/pq`; hand-rolled migration runner (`apps/api/internal/db/migrate.go`) | Build sandbox couldn't resolve their deps — `apps/api/README.md` |
| 11 | `docker-compose.yml` | `docker-compose.yaml` at repo root | — |
| 12 | Verify `docker compose up` boots the stack | **Not verified** — Dockerfiles/compose hand-checked only | No Docker daemon in build sandbox — root `README.md` |
| 13 | lint → test → build on every PR | Jobs `web` + `api` + `ci-ok` aggregator; workflow commands run by hand, not on a real runner at build time | Root `README.md` |
| 14 | Branch protection + required checks | Documented in `docs/branch-protection.md`; **not verified** as applied on GitHub | — |
| 15 | pino + request IDs | `log/slog` JSON + `X-Request-Id` middleware (`apps/api/internal/observability/`) | pino is Node; backend is Go — `docs/ADR-001-backend-language.md` |
| 16 | OpenTelemetry tracing (API → DB) | Not built; request-ID log correlation instead | `docs/ADR-002-observability.md` |
| 17 | `/metrics` via prom-client | Hand-rolled Prometheus text endpoint, counters only | `docs/ADR-002-observability.md` |
| 18 | PostHog/Plausible placeholder | Console-only `track`/`identify` hook (`apps/web/src/lib/analytics.ts`) | As intended (placeholder) |
| 20 | One-page ADR template | No template file; two filled ADRs (`ADR-001`, `ADR-002`) | — |
