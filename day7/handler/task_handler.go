package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"taskmanager/models"
	"taskmanager/service"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskservice service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskservice,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
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

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user Id", http.StatusBadRequest)
		return
	}

	var tasks []models.Task
	tasks, err = h.taskService.GetTasksForUser(id)
	if err != nil {
		http.Error(w, "Invalid user Id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
