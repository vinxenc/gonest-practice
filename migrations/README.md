# Database migrations

Versioned SQL migrations for the `employees` schema, in
[golang-migrate](https://github.com/golang-migrate/migrate) format. Each
migration is a pair of files:

- `NNNNNN_<name>.up.sql` — applies the change
- `NNNNNN_<name>.down.sql` — reverts it

Migrations are applied in ascending version order and reverted in descending
order. One migration per table, preceded by a bootstrap that creates the schema
and the gender enum every table depends on:

| Version  | Migration                  | Creates                                  |
| -------- | -------------------------- | ---------------------------------------- |
| `000001` | create_employees_schema    | `employees` schema + `employee_gender`   |
| `000002` | create_employee            | `employees.employee`                     |
| `000003` | create_department          | `employees.department`                   |
| `000004` | create_department_employee | `employees.department_employee` (+ FKs)  |
| `000005` | create_department_manager  | `employees.department_manager` (+ FKs)   |
| `000006` | create_salary              | `employees.salary` (+ FK)                |
| `000007` | create_title               | `employees.title` (+ FK)                 |

Order matters: `department_employee`/`department_manager` reference both
`employee` and `department`, so those tables come first; `salary` and `title`
reference `employee`.

## Running

`DATABASE_URL` points at the database (defaults below match
[`docker-compose.yml`](../docker-compose.yml)):

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/employees?sslmode=disable"
```

Using the standalone CLI (`brew install golang-migrate`):

```bash
migrate -path migrations -database "$DATABASE_URL" up        # apply all
migrate -path migrations -database "$DATABASE_URL" down 1    # revert the latest
migrate -path migrations -database "$DATABASE_URL" version   # current version
```

Or without installing anything (the `postgres` build tag pulls in the driver):

```bash
go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
  -path migrations -database "$DATABASE_URL" up
```

## Relationship to the sample dump

These migrations recreate the **same** schema that `docker compose up` imports
from the neon sample dump (see the project README). They are an alternative way
to build the schema from scratch — for a fresh/empty database, or as the basis
for evolving it going forward. Don't run them against a database that already
has the dump imported; the objects would already exist.

## Adding a migration

Create the next numbered pair, e.g.:

```bash
migrate create -ext sql -dir migrations -seq add_employee_email
# -> migrations/000008_add_employee_email.up.sql
# -> migrations/000008_add_employee_email.down.sql
```

Keep every `up` reversible by its matching `down`.
