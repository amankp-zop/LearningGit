package taskservice

import (
	"gofr.dev/pkg/gofr"

	taskmodel "assignment8/models/task"
)

type TaskRepositoryPort interface {
	CreateTask(c *gofr.Context, task *taskmodel.Task) error
	GetTasksByUserID(c *gofr.Context, userID int) ([]taskmodel.Task, error)
	UpdateTask(c *gofr.Context, task *taskmodel.Task) error
	MarkTaskComplete(c *gofr.Context, taskID int) error
	DeleteTask(c *gofr.Context, taskID int) error
}
