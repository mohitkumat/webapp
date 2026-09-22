# ADR-001: Backend Language — Go over Node.js

- **Status:** Accepted
- **Date:**
- **Decider(s):**

## Context
The monorepo template needs one backend language for all new projects. Both Go and
Node.js/TypeScript were viable. This decision only applies to the API/worker layer —
the frontend stays React + TypeScript either way.

## Decision
Use **Go** for the backend (API + background workers).

## Alternatives Considered

| Option | Pros | Cons |
|---|---|---|
| **Go** (chosen) | Compiled, single static binary; strong concurrency (goroutines) for workers/webhooks; typed end-to-end; low memory footprint; fast startup in containers | Second language in the stack (JS + Go); smaller pool of Go-familiar hires; more verbose than TS for CRUD-heavy code |
| **Node.js + TypeScript** | One language across frontend/backend; huge package ecosystem; faster to prototype | Single-threaded event loop needs extra care under CPU-heavy or concurrent I/O load; runtime type safety weaker without discipline; dependency surface area (supply-chain risk) larger |

## Why Go Wins Here
1. **Durable background work.** Templates need workers, leases, retries, idempotent
   webhook processing — Go's goroutines + standard library make this straightforward
   without extra runtime libraries.
2. **Predictable performance under load.** Compiled binary, low GC pressure, consistent
   latency — matters once services carry real traffic, not just prototypes.
3. **Deployment simplicity.** One static binary per service → smaller container images,
   fewer runtime dependencies to patch/secure.
4. **Type safety at the API boundary.** Combined with OpenAPI-generated clients, Go's
   static typing reduces contract drift between frontend and backend.
5. **Operational maturity.** Native `context`, structured concurrency, and first-class
   tracing/metrics libraries (OpenTelemetry-Go) fit the observability requirements
   directly.

## Consequences
- Two languages in the repo (TypeScript for `apps/web`, Go for `apps/api`) — CI and
  tooling docs must cover both toolchains.
- Onboarding needs a short Go primer for engineers coming from JS-only backgrounds.
- Some rapid-prototyping tasks (quick scripts, one-off scrapers) may still reach for
  Node/TS or Python — that's fine outside the core API/worker services.
- Monorepo build plan's Phase C is updated: Express → Go (net/http or chi), `pg` driver
  → `pgx` or `database/sql`, Jest/Vitest → Go's built-in `testing` package.

## Revisit Trigger
Reconsider if the team consistently struggles to hire/ramp on Go, or if a project's
scope is dominated by rapid CRUD prototyping where Node's speed-to-ship outweighs Go's
operational benefits.
