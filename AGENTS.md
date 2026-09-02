# AGENTS.md

Read this file with the Read tool at the start of every session, before any other work. Stay in this repository unless the user explicitly expands the scope.

## What this project is

`ratelim` (`github.com/ryoeuyo/ratelim`) is a single-instance HTTP rate limiter that sits in front of backend services as a reverse proxy.

Clients talk only to this process. It decides whether a request is allowed, then forwards allowed traffic to the matching upstream. Routing is config-driven: one limiter instance, many backends.

Example flow from `REQUIREMENTS.md`:

```
Client -> /api/users          -> rate limiter -> service.v1 /api/users
Client -> /api/orders/create  -> rate limiter -> service.v2 /api/orders/create
```

Go version: see `go.mod` (currently 1.26). Entry point: `cmd/main.go`.

## Architecture

```
cmd/main.go                 process entry: logger, HTTP server, signal-based shutdown
internal/server             server lifecycle (listen, graceful shutdown)
internal/handler            HTTP entry: reverse proxy (rate limiting will live here)
internal/model              domain types (rate-limit key today)
```

Request path today:

1. `main` builds `http.Server` with `handler.ProxyHandler` and `server.Runner`.
2. `Runner.Run` listens until SIGINT/SIGTERM, then `Shutdown` with a 10s timeout.
3. `ProxyHandler.ServeHTTP` builds `httputil.NewSingleHostReverseProxy` and forwards the request.

Intended path (not fully built):

1. Match request path to an upstream from config.
2. Build a rate-limit key (see `model.NewKey`).
3. Allow or reject (429 when over limit).
4. Reverse-proxy allowed requests to the chosen host, preserving path/query unless config says otherwise.

## Current code (WIP)

This is early scaffolding, not the finished limiter.

- `ProxyHandler` always targets `http://localhost:8080`. The process also listens on `localhost:8080`, so this is a self-proxy stub, not real routing.
- Rate limiting is not implemented. `model.NewKey(ip, uri)` exists for a per-client, per-URI key (`ip:uri`) and is unused.
- Routing config does not exist yet. `REQUIREMENTS.md` is the source of truth for that behavior.
- The reverse proxy is created per request. Prefer a long-lived proxy (or one per upstream) once routing exists.
- `NewSingleHostReverseProxy` rewrites `req.URL` but not `Host`; the handler sets `r.Host = target.Host` on purpose.

## Conventions

- Keep the limiter a single process. Do not add a cluster/coordination layer unless asked.
- Routing belongs in config, not hardcoded hosts in the handler.
- `internal/server` owns process lifetime only. Do not put proxy or limiter logic there.
- `internal/handler` is the HTTP edge: routing, limiting, then proxy.
- `internal/model` stays free of HTTP and I/O.
- Use `log/slog`. Match existing graceful-shutdown and error-handling style.
- Prefer standard library (`net/http`, `httputil`) unless a dependency is clearly needed.
- Do not commit `.idea/` or secrets. Do not commit unless the user asks.

## Working in this repo

- If `REQUIREMENTS.md` and the code disagree, treat requirements as the goal and the code as current state. Call out the gap when it matters.
- Small, focused changes. Do not invent extra services, metrics stacks, or config formats beyond what the task needs.
