package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrMissingID = errors.New("missing ID")

const pathParamsLength = 3

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type taskManager struct {
	tasks     []Task
	getNextID func() int
}

func idGenerator() func() int {
	id := 0

	return func() int {
		id++

		return id
	}
}

func (tm *taskManager) addTask(description string) Task {
	id := tm.getNextID()
	newTask := Task{
		ID:          id,
		Description: description,
		Completed:   false,
	}
	tm.tasks = append(tm.tasks, newTask)

	return newTask
}

func (tm *taskManager) getTaskByID(id int) (*Task, bool) {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			return &tm.tasks[i], true
		}
	}

	return nil, false
}

func (tm *taskManager) deleteTask(id int) bool {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)

			return true
		}
	}

	return false
}

func (tm *taskManager) completeTask(id int) (message string, statusCode int) {
	task, found := tm.getTaskByID(id)

	if !found {
		return fmt.Sprintf("Task with ID %d not found.", id), http.StatusNotFound
	}

	if task.Completed {
		return fmt.Sprintf("Task %d is already completed.", id), http.StatusConflict
	}

	task.Completed = true

	return fmt.Sprintf("Marked task %d as completed.", id), http.StatusAccepted
}

func (tm *taskManager) handlePostTask(w http.ResponseWriter, r *http.Request) {
	var err error
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, `{"error":"Unable to read request body"}`, http.StatusBadRequest)

		return
	}

	var input struct {
		Description string `json:"description"`
	}

	if err = json.Unmarshal(body, &input); err != nil || strings.TrimSpace(input.Description) == "" {
		http.Error(w, `{"error":"Invalid or missing description"}`, http.StatusBadRequest)

		return
	}

	newTask := tm.addTask(input.Description)
	response, _ := json.Marshal(newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (tm *taskManager) handleGetTask(w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := extractID(r.URL.Path)

	if err != nil {
		http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)

		return
	}

	task, found := tm.getTaskByID(id)
	if !found {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)

		return
	}

	response, _ := json.Marshal(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (tm *taskManager) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)

		return
	}

	if deleted := tm.deleteTask(id); deleted {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
}

func (tm *taskManager) handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := extractID(r.URL.Path)

	if err != nil {
		http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	msg, status := tm.completeTask(id)
	resp := map[string]string{"message": msg}
	response, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func extractID(path string) (int, error) {
	parts := strings.Split(path, "/")

	if len(parts) < pathParamsLength {
		return 0, ErrMissingID
	}

	return strconv.Atoi(parts[2])
}

func main() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil, // Default ServeMux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	manager := &taskManager{
		getNextID: idGenerator(),
	}

	manager.addTask("Go to Gym")
	manager.addTask("Buy groceries")

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			manager.handlePostTask(w, r)
		} else {
			http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			manager.handleGetTask(w, r)
		case http.MethodDelete:
			manager.handleDeleteTask(w, r)
		case http.MethodPut:
			manager.handleCompleteTask(w, r)
		default:
			http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running at http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
