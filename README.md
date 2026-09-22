# Project Template

Quick-start monorepo: React + TypeScript frontend, Go backend, Postgres, containerized.

## Structure
```
apps/web/       React + TypeScript UI (Vite)
apps/api/       Go API + background workers
packages/       Shared code (types, config) — added as needed
docs/           PRD, ADRs, runbooks
```

## Status
Phase C complete — Go API scaffolded, tested against real Postgres.

- [x] A. Repo foundation
- [x] B. Frontend (Vite + React + TS)
- [x] C. Backend (Go API + Postgres)
- [ ] D. Containers (Docker + Compose)
- [ ] E. CI (GitHub Actions)
- [ ] F. Observability (logs, traces, metrics)
- [ ] G. Docs (PRD/ADR templates)

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

### Full stack — lands in Phase D
```bash
docker compose up   # starts web + api + postgres locally
```

## Tool versions
Pinned in `.tool-versions`. After installing tools, run
`bash scripts/record-versions.sh > workstation-versions.txt` and commit the result.

## Contributing
- One feature per branch: `feature/short-description`
- PRs must pass lint + test + build (enforced by CI once Phase E lands)
