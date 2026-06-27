# gonest-practice

A [Fiber](https://gofiber.io/) web server in Go organized in a NestJS-style
modular architecture, using [Uber fx](https://uber-go.github.io/fx/) for
dependency injection, [Huma](https://huma.rocks/) for automatic OpenAPI
generation, and [GORM](https://gorm.io/) for PostgreSQL data access.

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
| GET    | `/employees`    | List employees (offset pagination)   | JSON              |
| GET    | `/docs`         | Swagger UI (interactive API docs)    | HTML              |
| GET    | `/openapi.json` | Auto-generated OpenAPI 3.1 spec      | JSON              |

`GET /employees` accepts `limit` (1–100, default 20) and `offset` (default 0)
query parameters and returns the page plus pagination metadata.

The OpenAPI spec and docs are generated automatically by Huma from the
registered operations — no separate spec file to maintain.

```bash
curl http://localhost:3000/health
# {"status":"ok"}

curl "http://localhost:3000/employees?limit=2"
# {"employees":[{"id":10001,"firstName":"Georgi", ...}], "limit":2, "offset":0, "total":300024}
```

## Architecture

The app follows a NestJS-style module layout. Each feature is a self-contained
`fx.Module` with its own Repository → Service → Controller layers. Each module
contributes its controller(s) to a `"controllers"` fx value group, and `core`
collects the whole group and registers every route — so a module is wired up
just by including it.

```text
.
├── docker-compose.yml           # PostgreSQL 17 + first-run import of the sample DB
├── db/
│   ├── download.sh              # fetches the Employees sample dump into db/dump/
│   └── init/                    # docker-entrypoint-initdb.d scripts (pg_restore)
└── src/
    ├── main.go                  # bootstraps via core.Server(...)
    ├── config/
    │   └── settings.go          # env-validated Settings (go-env-validator) + fx provider
    ├── core/
    │   ├── bootstrap.go         # core.Server factory + initServer (NestFactory-like)
    │   ├── doc.go               # package design & rationale
    │   ├── router.go            # Controller interface, AsController, registerRoutes
    │   └── server.go            # Fiber + Huma providers and server lifecycle
    └── modules/
        ├── gormModule/          # shared *gorm.DB (PostgreSQL) provider, reused by feature modules
        │   ├── gorm.module.go       # fx.Module exposing the connection
        │   └── gorm.provider.go     # NewGorm: opens GORM + lifecycle close
        ├── healthModule/        # example feature module
        │   ├── health.module.go     # fx providers + AsController(HealthController)
        │   ├── health.controller.go # registers Huma routes
        │   ├── health.service.go    # business logic
        │   ├── health.repository.go # data-access layer
        │   └── health.dto.go        # request/response types
        └── employeeModule/      # Employees feature (GORM-backed)
            ├── employee.module.go     # fx providers + AsController(EmployeeController)
            ├── employee.controller.go # GET /employees route
            ├── employee.service.go    # business logic + pagination clamping
            ├── employee.repository.go # GORM data-access + EmployeeReader interface
            ├── employee.entity.go     # GORM entities for the 6 employees tables
            └── employee.dto.go        # request/response types
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
        gormModule.GormModule, // shared infrastructure module
        healthModule.HealthModule,
        employeeModule.EmployeeModule,
    )
    app.Run()
}
```

`core.Server` provides the Fiber + Huma server and registers a single
`initServer` invoke, which collects every controller in the `"controllers"`
group, triggers `registerRoutes`, and then ties the server to the fx lifecycle.

### Shared modules

Infrastructure that several features need is its own `fx.Module`, included once
at the composition root. `gormModule.GormModule` provides a single `*gorm.DB`
(PostgreSQL) connection; because fx providers are application-wide, any other
module reuses it just by declaring `*gorm.DB` as a constructor parameter — no
re-wiring per module. This mirrors a NestJS global database module
(`TypeOrmModule.forRoot()`). `employeeModule`'s repository depends on `*gorm.DB`
and resolves the shared handle automatically.

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
