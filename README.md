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
Phase D complete — full stack runs with one command via Docker Compose.

- [x] A. Repo foundation
- [x] B. Frontend (Vite + React + TS)
- [x] C. Backend (Go API + Postgres)
- [x] D. Containers (Docker + Compose)
- [ ] E. CI (GitHub Actions)
- [ ] F. Observability (logs, traces, metrics)
- [ ] G. Docs (PRD/ADR templates)

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

## Tool versions
Pinned in `.tool-versions`. After installing tools, run
`bash scripts/record-versions.sh > workstation-versions.txt` and commit the result.

## Contributing
- One feature per branch: `feature/short-description`
- PRs must pass lint + test + build (enforced by CI once Phase E lands)
