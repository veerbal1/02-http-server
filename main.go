package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var taskList = []Task{}
var taskListMu sync.Mutex
var requestCounter = 0
var requestCounterMu sync.Mutex

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func requestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCounterMu.Lock()
		requestCounter++
		requestID := "req-" + strconv.Itoa(requestCounter)
		requestCounterMu.Unlock()

		w.Header().Set("X-Request-ID", requestID)
		next(w, r)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		requestID := w.Header().Get("X-Request-ID")
		fmt.Printf("request_id=%s method=%s path=%s duration=%s\n", requestID, r.Method, r.URL.Path, duration)
	}
}

func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := w.Header().Get("X-Request-ID")
				fmt.Printf("request_id=%s method=%s path=%s panic=%v\n", requestID, r.Method, r.URL.Path, err)
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
			}
		}()
		next(w, r)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var reqBody CreateTaskRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			writeErr := writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "invalid JSON body",
			})
			if writeErr != nil {
				return
			}
			return
		}

		if len(strings.TrimSpace(reqBody.Title)) == 0 {
			writeErr := writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "title is required"})
			if writeErr != nil {
				return
			}
			return
		}

		taskListMu.Lock()
		task := Task{Title: reqBody.Title, Done: false, ID: len(taskList) + 1}
		taskList = append(taskList, task)
		taskListMu.Unlock()

		writeErr := writeJSON(w, http.StatusCreated, task)
		if writeErr != nil {
			return
		}
		return
	}

	if r.Method == http.MethodGet {
		taskListMu.Lock()
		tasks := make([]Task, 0, len(taskList))
		tasks = append(tasks, taskList...)
		taskListMu.Unlock()

		writeErr := writeJSON(w, http.StatusOK, tasks)
		if writeErr != nil {
			return
		}
		return
	}

	w.Header().Set("Allow", "GET, POST")
	writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
}

func main() {
	http.HandleFunc("/health", requestIDMiddleware(recoveryMiddleware(loggingMiddleware(healthHandler))))
	http.HandleFunc("/tasks", requestIDMiddleware(recoveryMiddleware(loggingMiddleware(createTaskHandler))))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Got error: ", err)
	}
}
