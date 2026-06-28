# gonest-practice

A Go HTTP server organized in a NestJS-style modular architecture, built on the
[gonest](https://github.com/0xfurai/gonest) framework — which provides the
dependency-injection container, module system, router, request pipeline,
lifecycle hooks, and automatic OpenAPI/Swagger generation — with
[GORM](https://gorm.io/) for PostgreSQL data access.

## Requirements

- Go 1.26+
- Docker (for the local PostgreSQL database)

## Getting started

Start the PostgreSQL database and import the sample data (see
[Database](#database)), then install dependencies and run the server:

```bash
docker compose up -d     # start PostgreSQL
# ...then import the Employees database (one time) — see "Database" below
go mod tidy
go run ./src
```

The server listens on `http://localhost:3000` by default. Override the port via
the `PORT` environment variable (see [Configuration](#configuration)):

```bash
PORT=4567 go run ./src   # listens on http://localhost:4567
```

## Database

The app reads from the
[Employees sample database](https://github.com/neondatabase/postgres-sample-dbs#employees-database)
(~300k employees across 6 related tables in an `employees` schema).
[`docker-compose.yml`](docker-compose.yml) runs PostgreSQL 17; import the sample
data once with `pg_restore`.

**1. Start PostgreSQL:**

```bash
docker compose up -d
```

**2. Download the dump** into `db/dump/`, which the container mounts at `/dump`.
Despite the `.sql.gz` name it is a PostgreSQL **custom-format** archive (~33 MB),
so it is restored with `pg_restore` (not `psql`):

```bash
curl -fSL -o db/dump/employees.sql.gz \
  https://raw.githubusercontent.com/neondatabase/postgres-sample-dbs/main/employees.sql.gz
```

**3. Restore it** into the `employees` database (uses the `pg_restore` bundled in
the container, so no host PostgreSQL tools are required):

```bash
docker compose exec postgres \
  pg_restore -U postgres -d employees --no-owner --no-privileges /dump/employees.sql.gz
```

Data persists in a named volume across restarts. To start over, run
`docker compose down -v`, then repeat from step 1.

The default database settings match the compose file, so no configuration is
needed for local development. Point the app at another database via the `DB_*`
variables (see [Configuration](#configuration)).

## Configuration

Runtime configuration is read from environment variables, validated, and
defaulted up front by
[go-env-validator](https://github.com/philiprehberger/go-env-validator). The
validated values live in a single `config.Settings` struct
([`src/config/settings.go`](src/config/settings.go)); invalid configuration
fails app startup with every offending variable reported at once.

| Variable      | Type   | Default      | Description                              |
| ------------- | ------ | ------------ | ---------------------------------------- |
| `PORT`        | int    | `3000`       | TCP port the HTTP server binds to        |
| `DB_HOST`     | string | `localhost`  | PostgreSQL server hostname               |
| `DB_PORT`     | int    | `5432`       | PostgreSQL server port                   |
| `DB_USER`     | string | `postgres`   | PostgreSQL user                          |
| `DB_PASSWORD` | string | `postgres`   | PostgreSQL password                      |
| `DB_NAME`     | string | `employees`  | PostgreSQL database name                 |
| `DB_SSLMODE`  | string | `disable`    | libpq `sslmode` (e.g. `require`)         |

`core.New` calls `config.Load` once at startup, then shares the validated
`Settings` through gonest's DI container as a global value provider — every
module resolves the same instance, and `main` derives the listen address from
`Settings.Port`. To add a new setting, add a field with an `env` struct tag to
`Settings`:

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
| GET    | `/employees`    | List employees (offset pagination)   | JSON              |
| GET    | `/swagger`      | Swagger UI (interactive API docs)    | HTML              |
| GET    | `/swagger/json` | Auto-generated OpenAPI 3.0 spec      | JSON              |

`GET /employees` accepts `limit` (1–100, default 20) and `offset` (default 0)
query parameters and returns the page plus pagination metadata.

The OpenAPI spec and Swagger UI are generated automatically by gonest's swagger
module from the route metadata each controller declares (`Summary`/`Tags`/
`Response`) — no separate spec file to maintain.

```bash
curl http://localhost:3000/health
# {"status":"ok"}

curl "http://localhost:3000/employees?limit=2"
# {"employees":[{"id":10001,"firstName":"Georgi", ...}], "limit":2, "offset":0, "total":300024}
```

## Architecture

The app follows a NestJS-style module layout. Each feature is a self-contained
gonest module (`gonest.NewModule`) with its own Repository → Service →
Controller layers. A module lists its controllers in its `Controllers`, and
gonest registers every controller's routes automatically when the module is
imported — so a feature is wired up just by including it at the composition root.

```text
.
├── docker-compose.yml           # PostgreSQL 17 for local development
├── db/
│   └── dump/                    # downloaded sample dump (git-ignored); restored via pg_restore
├── migrations/                  # golang-migrate SQL migrations (one per table)
└── src/
    ├── main.go                  # bootstraps via core.New(...) + ListenAndServeWithGracefulShutdown
    ├── config/
    │   └── settings.go          # env-validated Settings (go-env-validator)
    ├── core/
    │   ├── core.go              # composition root: core.New builds the gonest app (config + GORM + swagger + features)
    │   └── doc.go               # package design & rationale
    └── modules/
        ├── gormModule/          # shared *gorm.DB (PostgreSQL) provider, reused by feature modules
        │   ├── gorm.module.go       # global gonest module exporting the connection
        │   └── gorm.provider.go     # NewGorm + lifecycle close on shutdown
        ├── healthModule/        # example feature module
        │   ├── health.module.go     # gonest module: Controllers + providers
        │   ├── health.controller.go # registers routes via Register(gonest.Router)
        │   ├── health.service.go    # business logic
        │   ├── health.repository.go # data-access layer
        │   └── health.dto.go        # request/response types
        └── employeeModule/      # Employees feature (GORM-backed)
            ├── employee.module.go     # gonest module: Controllers + providers (gonest.Bind for the reader interface)
            ├── employee.controller.go # GET /employees route
            ├── employee.service.go    # business logic + pagination clamping
            ├── employee.repository.go # GORM data-access + EmployeeReader interface
            ├── employee.entity.go     # GORM entities for the 6 employees tables
            └── employee.dto.go        # request/response types
```

### Dependency injection

gonest resolves each module's `Repository → Service → Controller` graph by type,
so there is no hand-written intra-module wiring, and the HTTP server lifecycle
(start/stop) is managed by the framework.

`main.go` bootstraps the app with the `core.New(...)` factory — analogous to
NestJS's `NestFactory.create(AppModule)` — which composes the infrastructure and
feature modules under a single root module and returns the runnable
`Application` plus the validated `Settings`. Including a module is all that's
needed:

```go
func main() {
    app, settings, err := core.New(
        healthModule.HealthModule,
        employeeModule.EmployeeModule,
    )
    if err != nil {
        log.Fatal(err)
    }

    addr := fmt.Sprintf(":%d", settings.Port)
    if err := app.ListenAndServeWithGracefulShutdown(addr); err != nil {
        log.Fatal(err)
    }
}
```

`core.New` wires the shared infrastructure (the validated config and the GORM
connection) and the swagger module, imports the feature modules, and builds the
gonest `Application`. `main` then calls `ListenAndServeWithGracefulShutdown`,
which compiles the module tree, registers every controller's routes, starts the
HTTP server, and on SIGINT/SIGTERM runs the framework's shutdown hooks.

### Shared modules

Infrastructure that several features need is its own gonest module, imported once
at the composition root (inside `core.New`). `gormModule.GormModule` provides a
single `*gorm.DB` (PostgreSQL) connection and is marked `Global`, so any other
module reuses it just by declaring `*gorm.DB` as a constructor parameter — no
re-wiring per module. This mirrors a NestJS global database module
(`TypeOrmModule.forRoot()`). `employeeModule`'s repository depends on `*gorm.DB`
and resolves the shared handle automatically. The validated `*config.Settings`
is shared the same way, as a global value provider, and its pool is closed
cleanly on shutdown via the framework's `OnApplicationShutdown` hook.

### Route registration

A controller is anything implementing `gonest.Controller`
(`Register(r gonest.Router)`). A module lists its controllers in
`gonest.ModuleOptions.Controllers`, and gonest registers every collected
controller's routes automatically — no central list to maintain.

To add a new feature module:

1. Create `src/modules/<name>Module/` with a Repository, Service, and a
   Controller exposing `Register(r gonest.Router)`.
2. Declare its gonest module, listing the controller and providers:

   ```go
   var FooModule = gonest.NewModule(gonest.ModuleOptions{
       Controllers: []any{FooController},
       Providers: []any{
           FooRepository,
           FooService,
       },
   })
   ```

   Use `gonest.Bind[Interface](Constructor)` in `Providers` to bind a
   constructor to an interface, as `employeeModule` does so its service depends
   on the `EmployeeReader` abstraction rather than the concrete repository.
3. Pass its module to `core.New(...)` in `main.go`.

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

CI enforces a minimum total statement coverage of 90% using
[go-test-coverage](https://github.com/vladopajic/go-test-coverage), configured
by [`.testcoverage.yml`](.testcoverage.yml). Check it locally the same way CI
does:

```bash
go tool lefthook run test   # runs the tests and the coverage threshold check
```

## Build

The entry package lives in `src/`, so pass an explicit output path with `-o`
(building `./src` directly would try to emit a binary named `src`, colliding
with the directory):

```bash
go build -o bin/server ./src
./bin/server
```
</content>
</invoke>
