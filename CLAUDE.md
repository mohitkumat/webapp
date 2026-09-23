# Web APP Project Template

## Reference
- Role: @reference/role.md
- Learning path: @reference/learning-path.md

This repo is the capstone for the learning path: a synthetic **AI engineering operations workspace**. The user is the accountable owner; Claude is a bounded agent. Scoring, day tests, and automatic failures in `reference/learning-path.md` are binding.

## Response style

- **Lead with the result.** What changed or what was found first, then evidence, then limitations and the next step (teach-back shape: decision → limitation → next step).
- **Label epistemic status** on any material claim: *observed* (I ran it / read it), *inferred*, *assumption*, *unknown*. Never present an inference as an observation.
- **Cite primary evidence.** `path:line` for code, the exact command + relevant output for runtime claims. No citation → say it's unverified.
- **Report failures verbatim.** Failing test/command output is quoted, not paraphrased. Skipped steps are stated as skipped.
- **Never claim done from green tests or merged code.** "Done" requires the observed live path (see Definition of done). If I couldn't observe it, say "implemented, not verified live" and give the command to verify.
- **Material decisions:** give a recommendation plus at least one viable alternative and the evidence for choosing — never a single unexamined option.
- **Discovery / ambiguity:** ask one question at a time; don't propose a feature until the problem and measure are accepted.
- Concise. No filler, no restating the request, no marketing tone.

## Working rules

### Task contract
Before non-trivial work, confirm (or ask for) the seven fields: **Objective, Authority, Scope, Constraints, Non-goals, Completion evidence, Escalation**. Work outside the stated scope is not absorbed silently — list it as a follow-up.

### Stop and escalate (do not proceed) when
- the change would alter architecture, a contract, or data authority that wasn't approved;
- evidence contradicts the plan or the stated starting state;
- an action is irreversible, touches secrets/credentials, or crosses a trust boundary;
- a control (CI, review, security check) fails and the "fix" would be to bypass it.

### Phases
Research → Decide → Plan → Implement → Review → Verify. Implementation never invents design; a design gap found mid-implementation goes back to Decide, visibly.

### Hard rules (automatic failures — never do these)
- Accept an AI (including my own) conclusion without checking its source.
- Expose a datastore driver, connection string, or credential to the browser.
- Take tenant/account scope from a client-supplied field — scope is server-derived only.
- Put secrets or real production data in prompts, code, logs, fixtures, or commits (synthetic only).
- Bypass or weaken a failing CI, review, security, or deployment control.
- Change architecture silently during implementation.
- Hide an ambiguous contract with retries or longer timeouts.
- Declare completion from tests while the live path is unobserved.
- Leave known broken behavior, orphan/superseded code, or unexplained data loss.

### Local-only runtime
Everything runs on the laptop: no hosted DB, object store, identity, model runtime, webhook relay, or deploy target. GitHub holds source + CI only. Published ports bind `127.0.0.1`, never `0.0.0.0` on the host side. No `latest` image tags. CI never deploys, never calls a hosted model, never needs real credentials.

## Architecture

### Target (capstone — `reference/learning-path.md` → Required capstone system)
```
Browser (React+TS, generated client only)
   │  HTTP, OpenAPI contract, synthetic signed token
   ▼
Go API (stateless) ──► PostgreSQL  (single authority; RLS default-deny per account)
   │                ──► Redis       (bounded read model: max rows, freshness, invalidation, partial-coverage flag)
   │                ──► evidence/   (local dir; retention + expiring-link simulation)
   │  outbox
   ▼
Go worker (leases, heartbeat, idempotency keys) ──► provider stub :8090 (signed webhooks back via inbox/dedup)
   │
   └─ AI workflow: deterministic local stub, structured output validated, tool authz outside model, approval before consequential action
Observability: X-Request-Id propagated browser → API → DB → worker → stub → UI; slog JSON; /metrics counters
```

### Invariants
- **One authority, one writer per data class.** Redis is never authoritative; reads declare coverage and freshness.
- **Contract first.** OpenAPI is the browser↔server contract; the web app uses only the generated client — no handwritten request/response types. A breaking change must fail CI before deploy.
- **Identity:** locally signed synthetic tokens; server derives principal + account; every handler authorizes; Postgres RLS enforces isolation as a second line.
- **Durable work:** at-most-once consequential external effect per idempotency key; every failure (duplicate, timeout, rate limit, process death, partial effect) reaches a deterministic terminal state without manual repair.
- **UI states:** every view defines loading, empty, populated, validation, permission-denied, failure, unavailable, recovery. Keyboard-operable, visible focus, works at 360 CSS px.
- **Forward-only migrations** in `apps/api/migrations/`, applied in filename order.
- **Artifact identity:** source revision → image digest → running container must match before acceptance.
- Observability stays stdlib (`log/slog`, hand-rolled `/metrics`) per `docs/ADR-002-observability.md`; seams live in `apps/api/internal/observability/`.

### Current state (observed 2026-09-23)
- `apps/web` — Vite + React 18 + TS + react-router; `/health` page with loading/ok/error states; `lib/analytics.ts` console-only.
- `apps/api` — Go stdlib `net/http`, `database/sql` + `lib/pq`, hand-rolled migrator, `items` resource, CORS → metrics → request-logger middleware chain (`cmd/api/main.go`).
- Postgres 16 via root `docker-compose.yaml`; CI runs web (lint/test/build) + api (gofmt/vet/build/test) + `ci-ok` aggregator.

### Known gaps vs target (fix when the relevant day comes up; don't paper over)
- No OpenAPI spec or generated client; `App.tsx` uses raw `fetch`.
- No auth, account scope, or RLS; `items` is global.
- No Redis, worker, outbox/inbox, provider stub, evidence store, or AI workflow.
- CI lacks OpenAPI validation and JSONL submission validation.
- `apps/web` uses npm + `package-lock.json` inside a pnpm workspace; `tsconfig.tsbuildinfo` is committed.
- Compose file is at repo root, not `deploy/compose.yaml`; no `tests/` or `docs/day-NN/` dirs yet.

## Delivery

- Branch per day: `day-N-topic`; one PR into `main`, linking the five daily files (`day-N-session.jsonl`, `-inputs.md`, `-artifacts.md`, `-review.md`, `-score.md` under `docs/day-NN/`), the measured score, and the CI run. Merge only after `CI OK` passes and review is done.
- Decisions go in `docs/ADR-NNN-*.md` (context, decision, alternatives table, consequences, revisit trigger). Product framing uses `docs/PRD-template.md`.
- Commit or push only when asked.

## Definition of done
1. Tests pass locally and in CI (quote the result).
2. Running artifact matches the reviewed source/config.
3. Real user path observed end to end (browser or curl) with the request ID visible in logs.
4. Failure/denied path exercised and behaves as designed.
5. No open high-severity review findings; superseded code removed.
6. Limitation and next step stated.

## Commands
```bash
docker compose up --build                 # full stack: web :3000, api :8080, postgres :5432
cd apps/web && npm run lint && npm test && npm run build
cd apps/api && gofmt -l . && go vet ./... && go test -p 1 ./...   # needs local Postgres; -p 1: packages share one DB
bash scripts/record-versions.sh > workstation-versions.txt
```
