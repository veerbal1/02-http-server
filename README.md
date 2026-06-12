# Stage 02 - HTTP Server

## Goal

Build a small in-memory task API in Go and understand the HTTP request lifecycle.

## Why This Matters For Go Backend/Platform Jobs

Most backend services are request/response systems. They receive requests, validate input, call business logic, return useful responses, log what happened, and shut down safely.

This stage turns Go fundamentals into the first real backend shape:

```text
request -> middleware -> handler -> response
```

## Concepts Practiced

- HTTP methods, paths, headers, and bodies
- HTTP status codes
- `net/http`
- handlers and handler functions
- routing
- JSON request decoding
- JSON response encoding
- validation and error responses
- request ID middleware
- logging middleware
- panic recovery middleware
- graceful shutdown
- handler tests with `httptest`

## Build Tasks

1. Create `task-http-api`. - Done
2. Add `GET /health`. - Done
3. Add in-memory `POST /tasks`. - Done
4. Add in-memory `GET /tasks`. - Done
5. Add JSON validation and useful error responses. - Done for current create/list slice
6. Add request ID middleware.
7. Add logging middleware.
8. Add panic recovery middleware.
9. Add graceful shutdown.
10. Add handler tests.

## Current Progress

Working now:

- Go module initialized as `task-http-api`.
- Server starts on port `8080`.
- `GET /health` returns `ok`.
- `POST /tasks` accepts JSON with a non-empty `title`.
- `POST /tasks` returns `201 Created` with the created task as JSON.
- Created tasks are stored in memory in `taskList`.
- `GET /tasks` returns the in-memory task list as JSON.
- Malformed JSON returns `400 Bad Request`.
- Empty or whitespace-only titles return `400 Bad Request`.
- Unsupported methods on `/tasks` return `405 Method Not Allowed`.

Current request flow:

```text
client request -> route -> handler -> method check -> decode JSON -> validate -> create/list -> response
```

Important limitation:

- Tasks are stored only in memory. Restarting the server clears the list.

## Test Tasks

- Verify `/health` returns success.
- Test task creation.
- Test task listing.
- Test invalid JSON.
- Test important middleware behavior where practical.

## Done Checklist

- [x] Server starts locally.
- [x] `/health` works.
- [x] Tasks can be created in memory.
- [x] Tasks can be listed from memory.
- [x] Invalid JSON returns a useful error.
- [ ] Request IDs appear in responses or logs.
- [ ] Logs show useful request information.
- [ ] Panic recovery returns a controlled response.
- [ ] Graceful shutdown is implemented.
- [ ] Handler tests pass.
- [x] README explains current request flow and endpoints.

## Endpoints

### `GET /health`

Checks whether the server is running.

Expected response:

```text
ok
```

### `POST /tasks`

Creates a task in memory.

Example:

```bash
curl -i -X POST http://localhost:8080/tasks -d '{"title":"learn Go HTTP"}'
```

Current response:

```json
{
  "ID": 1,
  "Title": "learn Go HTTP",
  "Done": false
}
```

### `GET /tasks`

Lists the in-memory tasks as JSON.

Example:

```bash
curl -i http://localhost:8080/tasks
```

Example response:

```json
[
  {
    "ID": 1,
    "Title": "learn Go HTTP",
    "Done": false
  }
]
```

## How To Run

```bash
go run .
```

Then test from another terminal:

```bash
curl -i http://localhost:8080/health
curl -i -X POST http://localhost:8080/tasks -d '{"title":"learn Go HTTP"}'
curl -i http://localhost:8080/tasks
```

## Current Manual Checks

```bash
go test ./...
```

Current result:

```text
?    task-http-api    [no test files]
```

## What I Learned

- A Go backend server stays running and waits for requests.
- A handler receives a request and writes a response.
- `w.Write([]byte("ok"))` writes raw bytes to the response.
- `json.NewEncoder(w).Encode(value)` converts a Go value to JSON and writes it directly to the response.
- HTTP response order matters: set headers first, write status second, write body third.
- `201 Created` is the right success status when `POST /tasks` creates a new task.
- `POST /tasks` and `GET /tasks` can share the same path but do different work based on the HTTP method.
- In-memory storage disappears when the server restarts.

## Public Post Ideas

- Built my first Go HTTP server with a `/health` endpoint.
- Added JSON task endpoints and learned the request -> handler -> response flow.
- Added middleware for request IDs, logging, and panic recovery.
- Tested Go HTTP handlers with `httptest`.
