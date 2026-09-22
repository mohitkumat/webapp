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
Phase B complete — frontend is scaffolded, tested, and builds clean.

- [x] A. Repo foundation
- [x] B. Frontend (Vite + React + TS)
- [ ] C. Backend (Go API + Postgres)
- [ ] D. Containers (Docker + Compose)
- [ ] E. CI (GitHub Actions)
- [ ] F. Observability (logs, traces, metrics)
- [ ] G. Docs (PRD/ADR templates)

## Setup

### Frontend (`apps/web`) — available now
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

### Backend (`apps/api`) — lands in Phase C
Not scaffolded yet. Placeholder folder only.

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
