# TodoList

A small REST API for managing to-dos, written in Go.

Started in 2021 as a Go + gin "hello world" (`/ping` only); rebuilt in 2026 as a working to-do API.

## Tech stack

- Go 1.24, standard library only (`net/http` routing with method + path patterns, added in Go 1.22)
- In-memory store guarded by a `sync.RWMutex`
- Table-driven tests with `net/http/httptest`

## API

| Method | Path          | Body                              | Response          |
| ------ | ------------- | --------------------------------- | ----------------- |
| GET    | `/ping`       |                                   | `{"message":"pong"}` |
| GET    | `/todos`      |                                   | list of to-dos    |
| POST   | `/todos`      | `{"title": "buy beans"}`          | 201 + created to-do |
| GET    | `/todos/{id}` |                                   | to-do or 404      |
| PATCH  | `/todos/{id}` | `{"title": "...", "done": true}` (either field) | updated to-do |
| DELETE | `/todos/{id}` |                                   | 204               |

Errors are returned as `{"error": "..."}` with a 4xx status.

## Getting started

```bash
go run .            # http://localhost:8080 (set PORT to change)
go test ./...

curl -X POST localhost:8080/todos -d '{"title":"buy beans"}'
curl -X PATCH localhost:8080/todos/1 -d '{"done":true}'
curl localhost:8080/todos
```

Data lives in memory and is cleared when the server restarts.

## 2026 update

- Implemented the actual to-do API (create, list, get, update, delete) with validation
- Replaced gin with the standard library router — Go 1.22+ supports `GET /todos/{id}`-style patterns, so the project now has zero dependencies
- Added tests for every endpoint and error case
- Removed committed macOS binaries (`main`, `first_tutorial`) and a broken empty submodule entry
