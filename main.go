package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
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

		task := Task{Title: reqBody.Title, Done: false, ID: len(taskList) + 1}
		taskList = append(taskList, task)

		writeErr := writeJSON(w, http.StatusCreated, task)
		if writeErr != nil {
			return
		}
		return
	}

	if r.Method == http.MethodGet {
		writeErr := writeJSON(w, http.StatusOK, taskList)
		if writeErr != nil {
			return
		}
		return
	}

	w.Header().Set("Allow", "GET, POST")
	writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/tasks", createTaskHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Got error: ", err)
	}
}
