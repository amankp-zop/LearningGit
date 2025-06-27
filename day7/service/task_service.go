package service

import (
	"errors"
	"taskmanager/models"
	"taskmanager/repository"
)

type TaskService interface {
	CreateTask(task *models.Task) error
	GetTasksForUser(userID int) ([]models.Task, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) CreateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("task title is required")
	}
	if task.UserId == 0 {
		return errors.New("user_id is required")
	}
	return s.taskRepo.CreateTask(task)
}

func (s *taskService) GetTasksForUser(userID int) ([]models.Task, error) {
	if userID == 0 {
		return nil, errors.New("user_id must be valid")
	}
	return s.taskRepo.GetTasksByUserID(userID)
}
