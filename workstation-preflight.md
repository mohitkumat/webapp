# Workstation preflight

- **Date:** 2026-09-23
- **Source revision:** `b9f2680` (plus uncommitted preflight changes listed at the end)
- **Result:** Checks 1–8 **pass**. Check 9 is pending the CI result on the `day-0-preflight` PR (recorded below once it finishes).

## Pass criteria

| Criterion | Status |
|---|---|
| Every tool exists | Pass |
| Every local service starts | Pass (Redis only as a throwaway container; it isn't in Compose until Day 5) |
| The AI CLI completes a session and keeps its transcript | Pass |
| Required ports are free | Pass (3000/8080/5432/6379/8090 were free before `compose up`) |
| At least 15 GB free after images | Pass (237 GiB free) |

---

## 1. OS and CPU architecture: pass
```
$ sw_vers; uname -m
ProductName:    macOS
ProductVersion: 26.6.2
BuildVersion:   25G83
arm64
```

## 2. CPU, memory, disk: pass (at the minimum)
```
$ sysctl -n hw.logicalcpu            → 10
$ sysctl -n hw.memsize               → 16 GB
$ df -h /                            → 460Gi size, 237Gi free
$ docker info (Docker Desktop VM)    → 10 CPUs, 7.65 GiB memory
```
16 GB is the minimum tier, so the rule is to run only the current day's services.
Peak memory with the full current stack (`docker stats --no-stream`) is web 18.3 MiB, api 15.8 MiB, postgres 56.0 MiB, about 90 MiB total.
Startup: `docker compose build` took 5 s (cached), `docker compose up -d --wait` took 11 s to reach all healthy.

## 3. Tool versions: pass (after re-pinning to current stable)
The pins had drifted: Node 22.13 was installed against a pin of 20.17 (engines `<21`), pnpm 12.5 against 9.11, and Go 1.22.5 against 1.23.2 (CI and the Dockerfile used 1.22).
Decision (owner): move to the newest stable versions, checked against the official sources on 2026-09-23 (`nodejs.org/dist/index.json`, `go.dev/dl/?mode=json`, `registry.npmjs.org/pnpm/latest`).

| Tool | Installed | Pin | Where pinned |
|---|---|---|---|
| Node.js | 24.21.0 (LTS "Krypton", via nvm) | 24.21.0 | `.tool-versions`, `package.json` engines `>=24 <25`, CI `setup-node`, `node:24.21.0-alpine` |
| pnpm | 12.6.0 | 12.6.0 | `.tool-versions`, engines `>=12` |
| Go | 1.27.1 (go.mod toolchain; system `/usr/local/go` is 1.22.5 and auto-switches) | 1.27.1 | `.tool-versions`, `go.mod`, CI `setup-go`, `golang:1.27.1-alpine` |
| Git / gh | 2.50.1 / 2.101.0 | — | |
| Python / uv | 3.12.1 / 0.12.17 | ≥ 3.12 | |
| Docker / Compose | 29.8.0 / v5.5.1 | — | |
| AI CLI | Claude Code 2.1.280 (Codex CLI 0.156.1 also installed) | — | |
| curl / jq / OpenSSL | 8.7.1 / 1.7.1 / 3.6.4 | — | |

Full output is in `workstation-versions.txt`. Checked again on the new toolchain:
```
apps/web (node v24.21.0): npm ci ✓, lint ✓, test 1/1 ✓, build ✓
apps/api (go1.27.1): gofmt clean, go vet ✓, go build ✓,
  go test -p 1 ./... → db ok, handlers ok, middleware ok, observability ok
```
`go mod tidy` collapsed `go 1.22` + `toolchain go1.27.1` into `go 1.27.1`.
`npm ci` on npm 11 warns that install scripts were not run (`npm install-scripts ls`). Lint, test and build still pass. Revisit when `apps/web` moves to pnpm.

## 4. Container reachable only through loopback: pass
```
$ lsof -nP -iTCP -sTCP:LISTEN | grep -E ':(3000|8080|5432)'
com.docke ... TCP 127.0.0.1:5432 (LISTEN)
com.docke ... TCP 127.0.0.1:8080 (LISTEN)
com.docke ... TCP 127.0.0.1:3000 (LISTEN)

$ curl -m 2 http://192.168.0.103:{3000,8080,5432}   # en0 LAN address
→ exit 7 (connection refused) on all three
```

## 5. AI CLI session and transcript: pass
Claude Code 2.1.280, authenticated, session completed. Its native transcript is
`~/.claude/projects/-Users-mohitkumat-Documents-work-webapp/<session-id>.jsonl`.
```
$ jq -e . 7c00c01c-….jsonl >/dev/null && echo "all lines parse"
all lines parse        # 341 lines: 46 user, 89 assistant, 7 system, …
```

## 6. Local image build and Compose start, no remote app images: pass (after a fix)
The first build **failed**:
```
#30 [web builder 5/6] COPY . .
#30 ERROR: cannot replace to directory .../src/node_modules/@testing-library/jest-dom with file
```
Cause: there was no `.dockerignore`, so the host's pnpm-symlinked `node_modules` was copied over the image's `npm install` output. The same gap also sent `apps/api/.env` into the API builder stage.
Fix: added `apps/web/.dockerignore` and `apps/api/.dockerignore`. Rebuild:
```
$ docker compose build            → Image webapp-web Built / Image webapp-api Built (exit 0)
$ docker compose up -d --wait     → postgres Healthy, api Healthy, web started (exit 0)
```
Only base images (`postgres:16-alpine`, `golang:1.27.1-alpine`, `node:24.21.0-alpine`, `alpine:3.20`, `nginx:1.27-alpine`) come from a registry. App images are built locally:

After the re-pin, `docker compose up -d --build --wait` reached all healthy in 78 s (including base-image pulls). The web builder was confirmed on `node:24.21.0-alpine@sha256:ebfe2f90…` with a `--no-cache` rebuild.

| Image | ID (after re-pin) |
|---|---|
| webapp-api (golang:1.27.1-alpine → alpine:3.20) | `sha256:2570f2d8…` |
| webapp-web (node:24.21.0-alpine → nginx:1.27-alpine) | `sha256:44250a77…` |
| postgres:16-alpine | `sha256:ef738a34…` |

After the re-pin, `GET /health` with `X-Request-Id: preflight-002` returned 200, and the same ID appears in the API log.

API container user: `uid=10001(appuser)` (non-root).

## 7. PostgreSQL + Redis health and persistence, evidence dir, token fixture: pass
```
$ curl -si -H 'X-Request-Id: preflight-001' http://127.0.0.1:8080/health
HTTP/1.1 200 OK
X-Request-Id: preflight-001
{"status":"ok","db":"ok"}

$ curl -s -X POST -d '{"name":"preflight"}' http://127.0.0.1:8080/items   → {"id":1,"name":"preflight"}
$ docker compose restart postgres api
$ curl -s http://127.0.0.1:8080/items                                      → [{"id":1,"name":"preflight"}]

# API log line for the same request ID
{"level":"INFO","msg":"request","request_id":"preflight-001","method":"GET","path":"/health","status":200,"duration_ms":0}
```
Redis (throwaway `redis:7.4-alpine`, `127.0.0.1:6379`, named volume, AOF on), removed afterwards:
```
PING → PONG;  SET preflight ok → OK;  docker restart;  GET preflight → ok
6379/tcp -> 127.0.0.1:6379
```
Evidence dir and token fixture (all under `.local/`, which is git-ignored):
- `.local/evidence/`: empty evidence directory.
- `.local/secrets/token-signing.key`: random 32-byte HMAC key (dir 0700, file 0600), never committed.
- `.local/fixtures/token-acct-a.jwt`: HS256 JWT with `sub=user-synthetic-1`, `account_id=acct-synthetic-a`, 24 h expiry. Its signature re-verifies with OpenSSL.
- *Assumption:* HS256 JWT is a placeholder format. The real token design is decided on Day 4/5.

## 8. Endpoints are loopback or private network only, no external services: pass
```
$ grep -nE 'https?://' docker-compose.yaml apps/api/internal/config/config.go apps/web/nginx.conf apps/web/vite.config.ts \
    | grep -vE '127\.0\.0\.1|localhost'
(no output, exit 1)
```
Inside Compose the API reaches Postgres at `postgres:5432` over the private Compose network.

## 9. GitHub repo, test PR, green CI: pending (see the CI result below)
```
$ gh repo view   → paavik/webapp, PRIVATE, default branch main
                   (git remote still says git@github.com:mohitkumat/webapp.git, so it redirects)
$ gh run list    → 35899612053  main  "Initi a CLAUDE.md file…"  **failure**
                   (the 4 earlier runs on main succeeded)
```
The failing job was `api`, step `Test`:
```
migrate_test.go:35: first run: record migration ../../migrations/0001_init.sql: pq: relation "schema_migrations" does not exist
--- FAIL: TestRunMigrations_IsIdempotent
health_test.go:13: run migrations: create schema_migrations table: pq: duplicate key value violates unique constraint "pg_type_typname_nsp_index"
--- FAIL: TestHealth_WhenDBReachable_ReturnsOK
```
- *Observed:* `internal/db/migrate_test.go:30` and `internal/handlers/testhelpers_test.go:26` both `DROP TABLE … items, schema_migrations` and re-migrate **the same database**. `go test ./...` runs packages in parallel.
- *Inferred:* the two packages race. The errors match exactly: one package drops the table the other just created, and two `CREATE TABLE IF NOT EXISTS` run at the same moment. The commit that failed changed no Go code, so this is a flaky test rather than a regression.
- *Not reproduced locally:* 0 failures in 35 parallel runs on this 10-core machine. The race depends on timing.

Fix (owner decision): run Go tests with packages in series. CI now runs `go test -p 1 ./... -v`, and both READMEs and CLAUDE.md say the same. Trade-off: the tests stay coupled through one shared database and rely on a flag rather than being isolated. The alternative, a database per package, stays open if the suite grows.

Also found:
- Branch protection can't be turned on: `GET …/branches/main/protection` returned `403 Upgrade to GitHub Pro or make this repository public`. So `docs/branch-protection.md` can't be applied on the current plan.
- The only collaborator is `mohitkumat`, so the assessor doesn't have read access yet.

---

## Changes made during preflight
- `apps/web/.dockerignore`, `apps/api/.dockerignore`: new (fix for check 6).
- `.gitignore`: added `.local/`.
- `scripts/record-versions.sh`: now records every tool check 3 requires. `workstation-versions.txt` regenerated.
- Version re-pin: `.tool-versions`, `package.json` engines, `apps/api/go.mod`, both Dockerfiles, `.github/workflows/ci.yml`.
- CI Go tests run with `-p 1`. `README.md`, `apps/api/README.md` and `CLAUDE.md` updated to match.
- `workstation-preflight.md`: this file.

## Still open (don't block preflight)
- The assessor doesn't have read access yet (the only collaborator is `mohitkumat`).
- Branch protection is unavailable on the private free plan (HTTP 403). Options: make the repo public, upgrade the plan, or accept that merges aren't enforced.
- The remote URL `mohitkumat/webapp` redirects to `paavik/webapp`. Consider `git remote set-url`.
