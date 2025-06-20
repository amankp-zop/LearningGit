package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Task represents a single to-do item.
type Task struct {
	ID          int
	Description string
	Completed   bool
}

// TaskManager manages a list of tasks and ID generation.
type TaskManager struct {
	tasks     []Task
	getNextID func() int
}

// IdGenerator returns a closure that generates incrementing IDs.
func IdGenerator() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

// AddTask adds a new task with the given description.
func (tm *TaskManager) AddTask(description string) {
	id := tm.getNextID()
	task := Task{
		ID:          id,
		Description: description,
		Completed:   false,
	}
	tm.tasks = append(tm.tasks, task)
	fmt.Printf("Task Added: %d - %s\n", task.ID, task.Description)
}

// GetSpecificTask retrieves the description of a task by ID.
func (tm *TaskManager) GetSpecificTask(id int) string {
	for _, task := range tm.tasks {
		if task.ID == id {
			return task.Description
		}
	}
	return fmt.Sprintf("Task with ID %d not found.", id)
}

// DeleteTask removes a task by ID.
func (tm *TaskManager) DeleteTask(id int) {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			fmt.Printf("Deleted task %d: %s\n", id, task.Description)
			return
		}
	}
	fmt.Printf("Task with ID %d not found.\n", id)
}

func (tm *TaskManager) CompleteTask(id int) string {
	for i, task := range tm.tasks {
		if task.ID == id {
			if task.Completed {
				return fmt.Sprintf("Task %d is already completed.", id)
			}
			tm.tasks[i].Completed = true
			return fmt.Sprintf("Marked task %d as completed.", id)
		}
	}
	return fmt.Sprintf("Task with ID %d not found.", id)
}

// HTTP handler for POST /tasks
func (tm *TaskManager) handlePostTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	description := strings.TrimSpace(string(body))
	if description == "" {
		http.Error(w, "Task description is empty", http.StatusBadRequest)
		return
	}
	tm.AddTask(description)
	fmt.Fprintf(w, "Task created: %s", description)
}

// HTTP handler for DELETE /tasks/{id}
func (tm *TaskManager) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Missing task ID in URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	tm.DeleteTask(id)
	fmt.Fprint(w, "Task deleted successfully")
}

func (tm *TaskManager) handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	result := tm.CompleteTask(id)
	if strings.Contains(result, "not found") {
		http.Error(w, result, http.StatusNotFound)
	} else {
		fmt.Fprintln(w, result)
	}
}

// HTTP handler for GET /tasks/{id}
func (tm *TaskManager) handleGetTask(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Missing task ID in URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	taskData := tm.GetSpecificTask(id)
	fmt.Fprint(w, taskData)
}

func main() {
	manager := &TaskManager{
		getNextID: IdGenerator(),
	}

	// Example initial tasks
	manager.AddTask("Go to Gym")
	manager.AddTask("Buy groceries")

	// Routes
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			manager.handlePostTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/complete") && r.Method == http.MethodPut {
			manager.handleCompleteTask(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			manager.handleGetTask(w, r)
		case http.MethodDelete:
			manager.handleDeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}

}
