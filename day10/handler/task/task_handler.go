package taskhandler

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"

	taskmodel "assignment8/models/task"
)

type TaskHandler struct {
	taskService TaskServicePort
}

func NewTaskHandler(taskService TaskServicePort) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(ctx *gofr.Context) (any, error) {
	var task taskmodel.Task

	err := ctx.Bind(&task)
	if err != nil {
		return nil, err
	}

	err = h.taskService.CreateTask(ctx, &task)
	if err != nil {
		return nil, err
	}

	return response.Raw{Data: task}, nil
}

func (h *TaskHandler) GetUserTasks(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	tasks, err := h.taskService.GetTasksForUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.Raw{Data: tasks}, nil
}

func (h *TaskHandler) UpdateTask(ctx *gofr.Context) (any, error) {
	var task taskmodel.Task

	err := ctx.Bind(&task)
	if err != nil {
		return nil, err
	}

	err = h.taskService.UpdateTask(ctx, &task)
	if err != nil {
		return nil, err
	}

	return response.Raw{Data: task}, nil
}

func (h *TaskHandler) DeleteTask(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	err = h.taskService.DeleteTask(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.Raw{
		Data: "Task deleted",
	}, nil
}

func (h *TaskHandler) MarkTaskComplete(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	err = h.taskService.MarkTaskComplete(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.Raw{
		Data: "Task marked as complete",
	}, nil
}
