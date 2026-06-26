# gonest-practice

A [Fiber](https://gofiber.io/) web server in Go organized in a NestJS-style
modular architecture, using [Uber fx](https://uber-go.github.io/fx/) for
dependency injection and [Huma](https://huma.rocks/) for automatic OpenAPI
generation.

## Requirements

- Go 1.26+

## Getting started

Install dependencies and run the server:

```bash
go mod tidy
go run ./src
```

The server listens on `http://localhost:3000`.

## Endpoints

| Method | Path            | Description                          | Response          |
| ------ | --------------- | ------------------------------------ | ----------------- |
| GET    | `/health`       | Health check                         | `{"status":"ok"}` |
| GET    | `/docs`         | Swagger UI (interactive API docs)    | HTML              |
| GET    | `/openapi.json` | Auto-generated OpenAPI 3.1 spec      | JSON              |

The OpenAPI spec and docs are generated automatically by Huma from the
registered operations — no separate spec file to maintain.

```bash
curl http://localhost:3000/health
# {"status":"ok"}
```

## Architecture

The app follows a NestJS-style module layout. Each feature is a self-contained
`fx.Module` with its own Repository → Service → Controller layers, and the
`core` module wires up the HTTP server and registers every module's routes
automatically.

```text
src/
├── main.go                      # composes fx modules and runs the app
├── core/
│   ├── core.module.go           # fx.Module: provides Fiber + Huma, invokes once
│   ├── router.go                # Route interface, AsRoute helper, registerRoutes
│   ├── router_test.go
│   └── server.go                # Fiber + Huma providers and server lifecycle
└── modules/
    └── healthModule/            # example feature module
        ├── health.module.go     # fx.Module("HealthModule", fx.Provide(...))
        ├── health.controller.go # registers Huma routes
        ├── health.service.go    # business logic
        ├── health.repository.go # data-access layer
        └── health.dto.go        # request/response types
```

### Dependency injection with fx

`main.go` only composes modules — there is no manual wiring:

```go
fx.New(
    core.Module,
    healthModule.HealthModule,
).Run()
```

fx resolves the `Repository → Service → Controller` graph by type, and the
server lifecycle (start/stop) is managed by fx hooks.

### Automatic route registration

Controllers join a `"routes"` fx value group via `core.AsRoute(...)`, and the
`core` module consumes the whole group to register every route. Adding a module
requires **no changes** to central wiring.

To add a new feature module:

1. Create `src/modules/<name>Module/` with a Repository, Service, and a
   Controller exposing `RegisterRoutes(api huma.API)`.
2. Declare its `fx.Module`, wrapping the controller constructor with
   `core.AsRoute(...)`:

   ```go
   var FooModule = fx.Module("FooModule",
       fx.Provide(
           NewFooRepository,
           NewFooService,
           core.AsRoute(NewFooController),
       ),
   )
   ```

3. Add the module to `fx.New(...)` in `main.go`.

## Development

Git hooks are managed by [lefthook](https://lefthook.dev/) and
linting/formatting by [golangci-lint](https://golangci-lint.run/).

Install the tools and the hooks:

```bash
brew install lefthook golangci-lint
lefthook install
```

Hooks:

- **pre-commit** — `golangci-lint fmt`, `go vet`, and `golangci-lint run --fix`
  (auto-fixes are restaged).
- **pre-push** — `go test ./...`.

Run them manually:

```bash
go test ./...
golangci-lint run ./...
```

## Build

The entry package lives in `src/`, so pass an explicit output path with `-o`
(building `./src` directly would try to emit a binary named `src`, colliding
with the directory):

```bash
go build -o bin/server ./src
./bin/server
```
