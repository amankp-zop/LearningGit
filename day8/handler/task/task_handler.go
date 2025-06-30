package taskhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	taskmodel "assignment8/models/task"
)

type TaskServicePort interface {
	CreateTask(task *taskmodel.Task) error
	GetTasksForUser(userID int) ([]taskmodel.Task, error)
	UpdateTask(task *taskmodel.Task) error
	MarkTaskComplete(taskID int) error
	DeleteTask(taskID int) error
}

type TaskHandler struct {
	taskService TaskServicePort
}

func NewTaskHandler(taskService TaskServicePort) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task taskmodel.Task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, "Invalid Task Body", http.StatusBadRequest)
		return
	}

	err = h.taskService.CreateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Unable to Find User", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	tasks, err := h.taskService.GetTasksForUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "Unable to Find User", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task taskmodel.Task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, "Invalid Task Body", http.StatusBadRequest)
		return
	}

	err = h.taskService.UpdateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Unable to Find User", http.StatusInternalServerError)
		return
	}
}

func (*TaskHandler) handleIDAction(w http.ResponseWriter, r *http.Request, action func(int) error, successMessage string) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := action(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": successMessage}); err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.handleIDAction(w, r, h.taskService.DeleteTask, "Task deleted")
}

func (h *TaskHandler) MarkTaskComplete(w http.ResponseWriter, r *http.Request) {
	h.handleIDAction(w, r, h.taskService.MarkTaskComplete, "Task marked as complete")
}
