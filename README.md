# gonest-practice

A simple [Fiber](https://gofiber.io/) web server written in Go.

## Requirements

- Go 1.26+

## Getting started

Install dependencies and run the server:

```bash
go mod tidy
go run ./src
```

The server listens on `http://localhost:3000`.

## Project layout

```
.
├── go.mod
├── go.sum
└── src/
    └── main.go   # entry point (package main)
```

## Endpoints

| Method | Path      | Description          | Response          |
| ------ | --------- | -------------------- | ----------------- |
| GET    | `/health` | Health check         | `{"status":"ok"}` |

Example:

```bash
curl http://localhost:3000/health
# {"status":"ok"}
```

## Build

The entry package lives in `src/`, so pass an explicit output path with `-o`
(building `./src` directly would try to emit a binary named `src`, colliding
with the directory):

```bash
go build -o bin/server ./src
./bin/server
```
