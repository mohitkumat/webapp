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
This is Phase A of the build plan — workspace skeleton only.

- [x] A. Repo foundation
- [ ] B. Frontend (Vite + React + TS)
- [ ] C. Backend (Go API + Postgres)
- [ ] D. Containers (Docker + Compose)
- [ ] E. CI (GitHub Actions)
- [ ] F. Observability (logs, traces, metrics)
- [ ] G. Docs (PRD/ADR templates)

## Setup (once later phases land)
```bash
pnpm install        # installs web workspace deps
docker compose up   # starts web + api + postgres locally
```

## Tool versions
Pinned in `.tool-versions`. After installing tools, run
`bash scripts/record-versions.sh > workstation-versions.txt` and commit the result.

## Contributing
- One feature per branch: `feature/short-description`
- PRs must pass lint + test + build (enforced by CI once Phase E lands)
