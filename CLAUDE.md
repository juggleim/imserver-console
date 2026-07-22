# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Admin console for [JuggleIM](https://github.com/juggleim/im-server): a Go/Gin backend (management APIs + an API gateway to the IM server) with a Vue 3 admin UI whose built output is **embedded into the Go binary**. See README.md for the deployment/feature overview.

## Commands

Backend (run from repo root — config is read from the relative path `conf/config.yml`):

```bash
go run .                 # start console on config.port (default 8091)
go build ./...
go test ./...
go test ./services -run TestQryMsgRealtime      # single test
RUN_DB_TEST=1 go test ./services -run Integration  # DB-backed tests (see below)
```

DB-backed tests are skipped unless `RUN_DB_TEST=1`; `TestMain` in `services/statisticservice_realtime_integration_test.go` then chdirs to the repo root and boots real config + MySQL, so they need a live database matching `conf/config.yml`.

Frontend (`webconsole/web`):

```bash
npm ci
npm run dev     # Vite dev server; proxies /admingateway -> http://127.0.0.1:8090
npm run build   # regenerates web/dist, which is what the Go binary embeds
npm test        # node --test over explicitly listed .mjs test files
```

Two gotchas: the Vite proxy targets port **8090** while `conf/config.yml` defaults to **8091** — align one of them when developing the UI against a local backend. And `npm test` names its test files explicitly in `package.json`; a new `*.test.mjs` file will not run until it is added to that script.

`webconsole/web/dist` is committed on purpose (it is `//go:embed`-ed by `webconsole/webload.go`). Any UI change must be rebuilt and the dist committed, or the served console will be stale. `routers/admin/` holds an older copy of the same assets and is not referenced by any Go code.

## Architecture

Request path: browser → embedded Vue SPA → `/admingateway/...` → Gin.

- `routers/router.go` — the single, flat route table for the whole console. Every endpoint is registered here under the `admingateway` prefix; adding a feature means adding a line here plus a handler in `apis/`.
- `apis/` — thin Gin handlers: bind the request, call a service, respond. `apis/validate.go` is the auth middleware applied to the entire group.
- `services/` — business logic; the only layer that talks to both `dbs/` and the IM server.
- `dbs/` — one GORM DAO struct per table, each with a `TableName()` method and its own query helpers. DAO structs double as the row model.
- `commons/` — `configures` (YAML config singleton), `dbcommons` (connection + migrations), `ctxs` (context keys + HTTP response helpers), `errs`, `imsdk`, `logs`, `tools`, `caches`.
- `apis/models/` and `services/models/` — request/response DTOs, kept separate from DAO structs.

### Responses and errors

Handlers respond via `ctxs.SuccessHttpResp` / `ctxs.FailHttpResp`, which return **HTTP 200 with a business `code` field** from `commons/errs/errs.go` (`AdminErrorCode_*`, 0 = success). Only auth and malformed-body failures use real 4xx statuses. Services return `(errs.AdminErrorCode, data)` rather than Go errors; follow that convention.

### Auth

`apis.Validate` guards every `/admingateway` route except `/login`. Two mutually exclusive modes:

- **Signature** (server-to-server): headers `signature`/`nonce`/`timestamp`, where signature = `SHA1(adminSecret + nonce + timestamp)`. Only active when both the header and `adminSecret` are set.
- **JWT** (browser): `Authorization` header holding a 24h HS256 token signed with `adminSecret` (falls back to a hardcoded default key when unset). The account is then re-checked against the DB and stashed in `ctxs.CtxKey_Account`.

Tenant scoping is not automatic: handlers read the target app from an `app_key` query param or request body field. The frontend also sends an `appkey` header (`webconsole/web/src/services/request.js`), and a 401 there triggers logout.

### Talking to the IM server

- `POST /admingateway/imapiagent` is a pass-through proxy: `services.ApiAgent` re-signs the request with the target app's own `appSecret` and forwards it to `imApiDomain/apigateway`. This is what the console's API-debugging UI uses.
- `commons/imsdk` caches one `JuggleIMSdk` client per appkey. Call `imsdk.Invalidate(appkey)` whenever an app's secret or domain changes, or stale clients keep signing with the old secret.
- `imAdminDomain` (JuggleIM admin API) is separate from `imApiDomain` (IM service API). Console→admin calls are signed with `services.GetImConsoleHeaders()` and hit the `/console/...` prefix.
- `/apps/serverlogs/{userconnect,connect,business}` (the connection inspector) proxy to the admin API's `/console/vlogs/query`, which greps the log files of an im node. That query is fenced on the node — 24h lookback, ≤1000 lines, 2 concurrent queries — and routed by `target_id` (falling back to `user_id`), so every call must carry a user id even when filtering by session. `services/serverlogservice.go` mirrors those bounds and decodes each returned log line into a map.

### Database migrations

`dbcommons.Upgrade()` runs at startup. Migrations are numbered SQL files embedded from `commons/dbcommons/sqls/` (e.g. `20260720.sql`); the applied version is stored in the `global_confs` row keyed `jchatdb_version`, and every file with a higher number is executed in order. To add a schema change, drop a new `YYYYMMDD.sql` in that directory — no code change needed. Files are matched by their leading 8 digits, so a suffixed name still parses.

## OpenSpec

This repo uses a spec-driven workflow (`openspec/` at the root for backend/product specs, `webconsole/web/openspec/` for UI ones). Capability specs live in `openspec/specs/<capability>/spec.md` as `Requirement:` / `Scenario:` blocks; in-flight work lives under `changes/` and is moved to `changes/archive/<date>-<name>/` once shipped. Specs are written in Chinese. When changing a behavior that has a spec, update the spec alongside the code.

## Frontend notes

Vue 3 + Vite + CoreUI, no TypeScript. `src/services/*.js` wraps the backend endpoints, `src/views/` holds the pages, and `src/i18n` + `src/locales/{en-US,zh-CN}` back the bilingual UI — any user-facing string needs entries in both locale directories.
