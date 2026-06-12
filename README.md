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

1. Create `task-http-api`.
2. Add `GET /health`.
3. Add in-memory `POST /tasks`.
4. Add in-memory `GET /tasks`.
5. Add JSON validation and useful error responses.
6. Add request ID middleware.
7. Add logging middleware.
8. Add panic recovery middleware.
9. Add graceful shutdown.
10. Add handler tests.

## Test Tasks

- Verify `/health` returns success.
- Test task creation.
- Test task listing.
- Test invalid JSON.
- Test important middleware behavior where practical.

## Done Checklist

- [ ] Server starts locally.
- [ ] `/health` works.
- [ ] Tasks can be created in memory.
- [ ] Tasks can be listed from memory.
- [ ] Invalid JSON returns a useful error.
- [ ] Request IDs appear in responses or logs.
- [ ] Logs show useful request information.
- [ ] Panic recovery returns a controlled response.
- [ ] Graceful shutdown is implemented.
- [ ] Handler tests pass.
- [ ] README explains request flow and endpoints.

## What I Learned

Fill this in as the stage progresses.

## Public Post Ideas

- Built my first Go HTTP server with a `/health` endpoint.
- Added JSON task endpoints and learned the request -> handler -> response flow.
- Added middleware for request IDs, logging, and panic recovery.
- Tested Go HTTP handlers with `httptest`.
