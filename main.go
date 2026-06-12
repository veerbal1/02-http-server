package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

var taskList = []Task{}

type CreateTaskRequest struct {
	Title string
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var reqBody CreateTaskRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(strings.TrimSpace(reqBody.Title)) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		task := Task{Title: reqBody.Title, Done: false, ID: len(taskList) + 1}
		taskList = append(taskList, task)

		w.Write([]byte("create task"))
		return
	}

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(taskList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/tasks", createTaskHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Got error: ", err)
	}
}
