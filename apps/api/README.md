# apps/api

Go backend — stdlib `net/http` (Go 1.22 method-based routing), `database/sql` +
`lib/pq` for Postgres, and a small hand-rolled SQL migration runner (no external
migration tool needed).

## Requires
- Go 1.22+
- A local Postgres reachable at `DATABASE_URL` (see `.env.example`)

## Run
```bash
cp .env.example .env   # then export the vars, or use direnv/dotenv of your choice
go run ./cmd/api
```
Server starts on `127.0.0.1:8080` and applies any new files in `migrations/`
automatically on startup.

## Test
```bash
go test ./... -v
```
Tests run against a real local Postgres (`TEST_DATABASE_URL` env var to override,
defaults to the same DB as `DATABASE_URL`). They reset and re-migrate the schema
before each run, so don't point this at a database with data you care about.

## Endpoints
| Method | Path      | Purpose                                      |
|--------|-----------|-----------------------------------------------|
| GET    | /health   | Liveness + DB reachability check              |
| GET    | /items    | List items (empty array, not null, when none) |
| POST   | /items    | Create an item — `{"name": "..."}`            |

## Structure
```
cmd/api/main.go          entrypoint: wires config, db, migrations, routes
internal/config/         env-based config with local defaults
internal/db/             connection + migration runner
internal/handlers/       HTTP handlers (one file per resource)
migrations/              plain .sql files, applied in filename order
```

## Why no chi / golang-migrate
This sandbox's network allowlist only reaches `github.com` directly (no
`proxy.golang.org`, no `go.uber.org`, no `gopkg.in`). `lib/pq` has zero
problematic transitive deps and resolves fine. `golang-migrate` pulls
`go.uber.org/atomic`, which doesn't. Go 1.22's stdlib router also makes chi
unnecessary for this template's scope. If your environment isn't restricted
this way, feel free to swap in `pgx` and `golang-migrate` — the code is
structured so that's a contained change (`internal/db`).
