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
`fx.Module` with its own Repository → Service → Controller layers. The `core`
package provides the HTTP server and registers the routes of every feature
module supplied to it as an explicit `[]core.Module` list.

```text
src/
├── main.go                      # composes fx modules + lists feature modules
├── core/
│   ├── bootstrap.go             # core.Server factory + initServer (NestFactory-like)
│   ├── doc.go                   # package design & rationale
│   ├── router.go                # Controller/Module interfaces, registerRoutes loop
│   ├── router_test.go
│   └── server.go                # Fiber + Huma providers and server lifecycle
└── modules/
    └── healthModule/            # example feature module
        ├── health.module.go     # fx providers + Module (Controllers())
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
NestJS's `NestFactory.create(AppModule)` — passing the feature modules and the
explicit list of modules to register, then calls `Run`:

```go
func provideModules(health *healthModule.Controller) []core.Module {
    return []core.Module{
        healthModule.NewModule(health),
    }
}

func main() {
    app := core.Server(
        healthModule.HealthModule,
        fx.Provide(provideModules),
    )
    app.Run()
}
```

`core.Server` provides the Fiber + Huma server and registers a single
`initServer` invoke, which triggers `registerRoutes` and then ties the server to
the fx lifecycle.

### Route registration

A feature is expressed through two small contracts in `core`:

- `Controller` — anything with `RegisterRoutes(api huma.API)`.
- `Module` — bundles a feature's controllers: `Controllers() []Controller`.

`core` registers everything with a plain nested loop — for each module, for each
controller, call `RegisterRoutes`. The module list is explicit and owned by the
composition root, so registration order is visible and greppable (no reflection
or hidden value groups).

To add a new feature module:

1. Create `src/modules/<name>Module/` with a Repository, Service, and a
   Controller exposing `RegisterRoutes(api huma.API)`.
2. Declare its `fx.Module` (providers) and a `Module` type bundling its
   controllers:

   ```go
   var FooModule = fx.Module("FooModule",
       fx.Provide(NewFooRepository, NewFooService, NewFooController),
   )

   type Module struct{ controllers []core.Controller }

   func NewModule(foo *Controller) *Module {
       return &Module{controllers: []core.Controller{foo}}
   }

   func (m *Module) Controllers() []core.Controller { return m.controllers }
   ```

3. Pass its `fx.Module` to `core.Server(...)` and add its `Module` to
   `provideModules` in `main.go`.

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
