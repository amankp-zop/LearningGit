package taskservice

import (
	"errors"

	taskmodel "assignment8/models/task"
)

var (
	ErrTaskTitleRequired = errors.New("task title is required")
	ErrUserIDRequired    = errors.New("user_id is required")
	ErrUserIDInvalid     = errors.New("user_id must be valid")
	ErrTaskIDRequired    = errors.New("task id is required")
)

type TaskService struct {
	repo TaskRepositoryPort
}

type TaskRepositoryPort interface {
	CreateTask(task *taskmodel.Task) error
	GetTasksByUserID(userID int) ([]taskmodel.Task, error)
	UpdateTask(task *taskmodel.Task) error
	MarkTaskComplete(taskID int) error
	DeleteTask(taskID int) error
}

func NewTaskService(repo TaskRepositoryPort) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *taskmodel.Task) error {
	if task.Title == "" {
		return ErrTaskTitleRequired
	}

	if task.UserID == 0 {
		return ErrUserIDRequired
	}

	return s.repo.CreateTask(task)
}

func (s *TaskService) GetTasksForUser(userID int) ([]taskmodel.Task, error) {
	if userID == 0 {
		return nil, ErrUserIDInvalid
	}

	return s.repo.GetTasksByUserID(userID)
}

func (s *TaskService) UpdateTask(task *taskmodel.Task) error {
	if task.ID == 0 {
		return ErrTaskIDRequired
	}

	return s.repo.UpdateTask(task)
}

func (s *TaskService) MarkTaskComplete(taskID int) error {
	if taskID == 0 {
		return ErrTaskIDRequired
	}

	return s.repo.MarkTaskComplete(taskID)
}

func (s *TaskService) DeleteTask(taskID int) error {
	if taskID == 0 {
		return ErrTaskIDRequired
	}

	return s.repo.DeleteTask(taskID)
}
