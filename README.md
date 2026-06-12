# Stage 02 - HTTP Server

## Status

Implementation complete. Final completion review is pending the explain-back checkpoint.

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
6. Add request ID middleware. - Done
7. Add logging middleware. - Done
8. Add panic recovery middleware. - Done
9. Add graceful shutdown. - Done
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
- Malformed JSON returns `400 Bad Request` with a JSON error body.
- Empty or whitespace-only titles return `400 Bad Request` with a JSON error body.
- Unsupported methods on `/tasks` return `405 Method Not Allowed` with a JSON error body and `Allow: GET, POST`.
- `GET /health` and `/tasks` responses include `X-Request-ID`.
- Request IDs are generated as `req-<number>` and protected with a mutex.
- Requests are logged with request ID, method, path, and duration.
- Panic recovery returns `500 Internal Server Error` with a safe JSON error body and logs panic details.
- Graceful shutdown listens for `Ctrl+C`/`SIGTERM` and shuts down with a 5 second timeout.
- In-memory task creation and listing are protected with a mutex and list snapshot.

Completion review:

- Automated tests pass.
- Manual curl checklist passed.
- README documents current behavior.
- Remaining checkpoint: explain the request flow, middleware chain, testing approach, in-memory limitation, and graceful shutdown in your own words.

Current request flow:

```text
client request -> request ID middleware -> recovery middleware -> logging middleware -> route handler -> method check -> decode JSON -> validate -> create/list -> response
```

Important limitation:

- Tasks are stored only in memory. Restarting the server clears the list.

## Test Tasks

- Verify `/health` returns success.
- Test `/health` with `httptest`. - Done
- Test task creation. - Done
- Test task listing. - Done
- Test invalid JSON. - Done
- Test important middleware behavior where practical.

## Done Checklist

- [x] Server starts locally.
- [x] `/health` works.
- [x] Tasks can be created in memory.
- [x] Tasks can be listed from memory.
- [x] Invalid JSON returns a useful error.
- [x] Request IDs appear in responses.
- [x] Logs show useful request information.
- [x] Panic recovery returns a controlled response.
- [x] Graceful shutdown is implemented.
- [x] Handler tests pass.
- [x] README explains current request flow and endpoints.
- [ ] User can explain the project flow in their own words.

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
  "id": 1,
  "title": "learn Go HTTP",
  "done": false
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
    "id": 1,
    "title": "learn Go HTTP",
    "done": false
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
curl -i -X POST http://localhost:8080/tasks -d '{bad json}'
curl -i -X DELETE http://localhost:8080/tasks
```

Graceful shutdown check:

```bash
go run .
# press Ctrl+C
```

Expected:

```text
shutting down
```

## Current Manual Checks

```bash
go test ./...
```

Current result:

```text
ok   task-http-api
```

Automated tests:

- `TestHealthHandler` checks `GET /health` returns `200 OK` and body `ok`.
- `TestCreateTaskHandlerBadJSON` checks malformed JSON returns `400`.
- `TestCreateTaskHandlerEmptyTitle` checks empty title returns `400`.
- `TestCreateTaskHandlerWhitespaceTitle` checks whitespace-only title returns `400`.
- `TestCreateTaskHandlerCreatesTask` checks successful task creation returns `201` and task JSON.
- `TestListTasksHandlerEmpty` checks an empty in-memory list returns `[]`.
- `TestListTasksHandlerWithTasks` checks populated in-memory tasks are returned.
- `TestTasksHandlerMethodNotAllowed` checks unsupported methods return `405`, `Allow: GET, POST`, and JSON error.
- `TestRequestIDMiddlewareAddsHeader` checks middleware adds `X-Request-ID` and still calls the wrapped handler.
- `TestRequestIDMiddlewareGeneratesDifferentIDs` checks sequential requests receive `req-1` and `req-2`.
- `TestRecoveryMiddlewareHandlesPanic` checks a panic becomes `500 Internal Server Error` with a safe JSON error.

Manual curl checklist passed:

| Test | Expected Result | Status |
| --- | --- | --- |
| `GET /health` | `200 OK` with `ok` | Passed |
| `GET /tasks` | `200 OK` with JSON task list | Passed |
| malformed JSON | `400 Bad Request` with JSON error | Passed |
| empty title | `400 Bad Request` with JSON error | Passed |
| `POST /tasks` | `201 Created` with created task JSON | Passed |
| list after create | task list contains created tasks | Passed |
| unsupported method | `405 Method Not Allowed`, `Allow: GET, POST`, JSON error | Passed |

## What I Learned

- A Go backend server stays running and waits for requests.
- A handler receives a request and writes a response.
- `w.Write([]byte("ok"))` writes raw bytes to the response.
- `json.NewEncoder(w).Encode(value)` converts a Go value to JSON and writes it directly to the response.
- HTTP response order matters: set headers first, write status second, write body third.
- `201 Created` is the right success status when `POST /tasks` creates a new task.
- JSON error responses are more useful than bare status codes because clients can read the reason.
- JSON tags let Go keep exported field names like `Title` while the public API uses names like `title`.
- Small helpers like `writeJSON` are useful once the same response-writing pattern appears multiple times.
- `POST /tasks` and `GET /tasks` can share the same path but do different work based on the HTTP method.
- `httptest` can call handlers directly without starting the real server on port `8080`.
- Middleware wraps handlers so shared behavior like request IDs can be applied without repeating code in every handler.
- A mutex protects shared data so one request updates a shared counter or task list at a time.
- For `GET /tasks`, copying the slice under lock gives the response a stable snapshot before JSON encoding.
- Logging middleware measures duration around the handler and logs `request_id`, method, path, and duration.
- Recovery middleware uses `defer` and `recover` to catch panics and return a controlled `500` response.
- Middleware order matters: request ID runs first, recovery wraps logging and handlers, and logging wraps the route handler.
- Graceful shutdown uses an explicit `http.Server`, listens for stop signals, and calls `Shutdown` with a timeout context.
- In-memory storage disappears when the server restarts.

## Completion Explain-Back

Before moving to the next stage, explain:

1. What happens when a request hits `/tasks`.
2. What the middleware layers do.
3. Why `httptest` is useful.
4. Why in-memory storage disappears after restart.
5. What graceful shutdown does.

## Public Post Ideas

- Built my first Go HTTP server with a `/health` endpoint.
- Added JSON task endpoints and learned the request -> handler -> response flow.
- Added middleware for request IDs, logging, and panic recovery.
- Tested Go HTTP handlers with `httptest`.
