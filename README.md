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

The server listens on `http://localhost:3000` by default. Override the port via
the `PORT` environment variable (see [Configuration](#configuration)):

```bash
PORT=4567 go run ./src   # listens on http://localhost:4567
```

## Configuration

Runtime configuration is read from environment variables, validated, and
defaulted up front by
[go-env-validator](https://github.com/philiprehberger/go-env-validator). The
validated values live in a single `config.Settings` struct
([`src/config/settings.go`](src/config/settings.go)); invalid configuration
fails app startup with every offending variable reported at once.

| Variable | Type | Default | Description                       |
| -------- | ---- | ------- | --------------------------------- |
| `PORT`   | int  | `3000`  | TCP port the HTTP server binds to |

`config.Load` is registered as an fx provider, so `Settings` is constructed (and
validated) as part of the dependency graph — `core` injects it and the server
binds to `Settings.Port`. To add a new setting, add a field with an `env` struct
tag to `Settings`:

```go
type Settings struct {
    Port int    `env:"PORT,default=3000"`
    Env  string `env:"APP_ENV,required,choices=development|staging|production"`
}
```

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
`fx.Module` with its own Repository → Service → Controller layers. Each module
contributes its controller(s) to a `"controllers"` fx value group, and `core`
collects the whole group and registers every route — so a module is wired up
just by including it.

```text
src/
├── main.go                      # bootstraps via core.Server(...)
├── config/
│   └── settings.go              # env-validated Settings (go-env-validator) + fx provider
├── core/
│   ├── bootstrap.go             # core.Server factory + initServer (NestFactory-like)
│   ├── doc.go                   # package design & rationale
│   ├── router.go                # Controller interface, AsController, registerRoutes
│   ├── router_test.go
│   └── server.go                # Fiber + Huma providers and server lifecycle
└── modules/
    └── healthModule/            # example feature module
        ├── health.module.go     # fx providers + AsController(HealthController)
        ├── health.controller.go # registers Huma routes
        ├── health.service.go    # business logic
        ├── health.repository.go # data-access layer
        └── health.dto.go        # request/response types
```

### Dependency injection with fx

fx resolves each module's `Repository → Service → Controller` graph by type, so
there is no hand-written intra-module wiring, and the server lifecycle
(start/stop) is managed by fx hooks.

`main.go` bootstraps the app with the `core.Server(...)` factory — analogous to
NestJS's `NestFactory.create(AppModule)` — passing the feature modules, then
calls `Run`. Including a module is all that's needed:

```go
func main() {
    app := core.Server(
        healthModule.HealthModule,
    )
    app.Run()
}
```

`core.Server` provides the Fiber + Huma server and registers a single
`initServer` invoke, which collects every controller in the `"controllers"`
group, triggers `registerRoutes`, and then ties the server to the fx lifecycle.

### Route registration

A controller is anything implementing `core.Controller`
(`RegisterRoutes(api huma.API)`). A module contributes its controller to the
`"controllers"` group with `core.AsController(...)`, and `core` registers every
collected controller with a plain loop — no central list to maintain.

To add a new feature module:

1. Create `src/modules/<name>Module/` with a Repository, Service, and a
   Controller exposing `RegisterRoutes(api huma.API)`.
2. Declare its `fx.Module`, wrapping the controller constructor with
   `core.AsController(...)`:

   ```go
   var FooModule = fx.Module("FooModule",
       fx.Provide(
           NewFooRepository,
           NewFooService,
           core.AsController(NewFooController),
       ),
   )
   ```

3. Pass its `fx.Module` to `core.Server(...)` in `main.go`.

## Development

Git hooks are managed by [lefthook](https://lefthook.dev/) and
linting/formatting by [golangci-lint](https://golangci-lint.run/). Both are
pinned as Go tool dependencies in `go.mod` (the `tool` directive), so there is
nothing to install separately — run them through `go tool`.

Install the git hooks:

```bash
go tool lefthook install
```

Hooks:

- **pre-commit** — `go tool golangci-lint fmt`, `go vet`, and
  `go tool golangci-lint run --fix` (auto-fixes are restaged).
- **pre-push** — `go test ./...`.

Run them manually:

```bash
go test ./...
go tool golangci-lint run ./...
```

## Build

The entry package lives in `src/`, so pass an explicit output path with `-o`
(building `./src` directly would try to emit a binary named `src`, colliding
with the directory):

```bash
go build -o bin/server ./src
./bin/server
```
