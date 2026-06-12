package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	healthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "ok" {
		t.Errorf("expected body %q, got %q", "ok", rr.Body.String())
	}
}

func TestCreateTaskHandlerBadJSON(t *testing.T) {
	body := strings.NewReader("{bad json}")
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rr := httptest.NewRecorder()

	createTaskHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	expected := `{"error":"invalid JSON body"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestCreateTaskHandlerEmptyTitle(t *testing.T) {
	body := strings.NewReader(`{"title":""}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rr := httptest.NewRecorder()

	createTaskHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	expected := `{"error":"title is required"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestCreateTaskHandlerWhitespaceTitle(t *testing.T) {
	body := strings.NewReader(`{"title":"   "}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rr := httptest.NewRecorder()

	createTaskHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	expected := `{"error":"title is required"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestCreateTaskHandlerCreatesTask(t *testing.T) {
	taskList = []Task{}

	body := strings.NewReader(`{"title":"learn Go HTTP"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rr := httptest.NewRecorder()

	createTaskHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var task Task
	err := json.NewDecoder(rr.Body).Decode(&task)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected id 1, got %d", task.ID)
	}
	if task.Title != "learn Go HTTP" {
		t.Errorf("expected title %q, got %q", "learn Go HTTP", task.Title)
	}
	if task.Done != false {
		t.Errorf("expected done false, got %v", task.Done)
	}
}

func TestListTasksHandlerEmpty(t *testing.T) {
	taskList = []Task{}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	createTaskHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expected := `[]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestListTasksHandlerWithTasks(t *testing.T) {
	taskList = []Task{
		{ID: 1, Title: "task one", Done: false},
		{ID: 2, Title: "task two", Done: true},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	createTaskHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var tasks []Task
	err := json.NewDecoder(rr.Body).Decode(&tasks)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "task one" {
		t.Errorf("expected first task title %q, got %q", "task one", tasks[0].Title)
	}
	if tasks[1].Title != "task two" {
		t.Errorf("expected second task title %q, got %q", "task two", tasks[1].Title)
	}
}

func TestTasksHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
	rr := httptest.NewRecorder()

	createTaskHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	allow := rr.Header().Get("Allow")
	if allow != "GET, POST" {
		t.Errorf("expected Allow header %q, got %q", "GET, POST", allow)
	}

	expected := `{"error":"method not allowed"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestRequestIDMiddlewareAddsHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	wrapped := requestIDMiddleware(healthHandler)
	wrapped(rr, req)

	reqID := rr.Header().Get("X-Request-ID")
	if reqID == "" {
		t.Error("expected X-Request-ID header to be set, got empty")
	}

	if rr.Body.String() != "ok" {
		t.Errorf("expected body %q, got %q", "ok", rr.Body.String())
	}
}
