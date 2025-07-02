package taskservice

import (
	"errors"

	"gofr.dev/pkg/gofr"

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

func NewTaskService(repo TaskRepositoryPort) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(c *gofr.Context, task *taskmodel.Task) error {
	if task.Title == "" {
		return ErrTaskTitleRequired
	}

	if task.UserID == 0 {
		return ErrUserIDRequired
	}

	return s.repo.CreateTask(c, task)
}

func (s *TaskService) GetTasksForUser(c *gofr.Context, userID int) ([]taskmodel.Task, error) {
	if userID == 0 {
		return nil, ErrUserIDInvalid
	}

	return s.repo.GetTasksByUserID(c, userID)
}

func (s *TaskService) UpdateTask(c *gofr.Context, task *taskmodel.Task) error {
	if task.ID == 0 {
		return ErrTaskIDRequired
	}

	return s.repo.UpdateTask(c, task)
}

func (s *TaskService) MarkTaskComplete(c *gofr.Context, taskID int) error {
	if taskID == 0 {
		return ErrTaskIDRequired
	}

	return s.repo.MarkTaskComplete(c, taskID)
}

func (s *TaskService) DeleteTask(c *gofr.Context, taskID int) error {
	if taskID == 0 {
		return ErrTaskIDRequired
	}

	return s.repo.DeleteTask(c, taskID)
}
