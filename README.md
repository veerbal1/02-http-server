# Golang HTTP API

A small in-memory task API built with Go's standard library. Zero external dependencies.

## Features

- REST-like task creation and listing (`POST /tasks`, `GET /tasks`)
- Health check endpoint (`GET /health`)
- JSON request/response with input validation
- Request ID middleware (`X-Request-ID` header on every response)
- Request logging middleware (method, path, duration, request ID)
- Panic recovery middleware (catches panics, returns safe 500 responses)
- Graceful shutdown on `SIGINT` / `SIGTERM` with 5-second drain timeout
- Thread-safe in-memory storage using `sync.Mutex`

## Quick Start

```bash
go run .
```

Server starts on `http://localhost:8080`.

## Endpoints

### `GET /health`

Returns the server status.

```bash
curl -i http://localhost:8080/health
```

```
HTTP/1.1 200 OK
X-Request-Id: req-1

ok
```

### `POST /tasks`

Creates a task. Requires a JSON body with a non-empty `title`.

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"learn Go HTTP"}'
```

```
HTTP/1.1 201 Created
Content-Type: application/json
X-Request-Id: req-2

{"id":1,"title":"learn Go HTTP","done":false}
```

### `GET /tasks`

Lists all tasks.

```bash
curl -i http://localhost:8080/tasks
```

```
HTTP/1.1 200 OK
Content-Type: application/json
X-Request-Id: req-3

[{"id":1,"title":"learn Go HTTP","done":false}]
```

### Error Responses

All errors return a JSON body with an `error` field and the appropriate HTTP status code.

| Scenario | Status | Body |
|---|---|---|
| Malformed JSON | `400 Bad Request` | `{"error":"invalid JSON body"}` |
| Empty or whitespace-only title | `400 Bad Request` | `{"error":"title is required"}` |
| Unsupported HTTP method | `405 Method Not Allowed` | `{"error":"method not allowed"}` + `Allow: GET, POST` header |
| Internal panic | `500 Internal Server Error` | `{"error":"internal server error"}` |

## Request Flow

```
client request
  -> request ID middleware
    -> panic recovery middleware
      -> logging middleware
        -> route handler
          -> method check
          -> decode/validate JSON
          -> business logic
          -> response
```

## Running Tests

```bash
go test ./...
```

## Graceful Shutdown

Press `Ctrl+C` or send `SIGTERM`. The server stops accepting new connections and waits up to 5 seconds for in-flight requests to complete before exiting.

```text
shutting down
```

## Limitations

- **In-memory storage only** — tasks are lost when the server restarts.
- **No persistent IDs** — task IDs reset to 1 on restart.
- **No `DELETE` or `PATCH`** — only create and list are implemented.
