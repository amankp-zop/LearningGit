package taskservice

import (
	"errors"
	"testing"

	taskmodel "assignment8/models/task"
)

type mockRepo struct{}

func (m *mockRepo) CreateTask(task *taskmodel.Task) error {
	if task.Title == "" {
		return ErrTaskTitleRequired
	}
	if task.UserID == 0 {
		return ErrUserIDRequired
	}
	task.ID = 1
	return nil
}
func (m *mockRepo) GetTasksByUserID(userID int) ([]taskmodel.Task, error) {
	if userID == 0 {
		return nil, ErrUserIDInvalid
	}
	if userID == 1 {
		return []taskmodel.Task{{ID: 1, Title: "Task", Description: "Desc", UserID: 1, IsCompleted: false}}, nil
	}
	return []taskmodel.Task{}, nil
}
func (m *mockRepo) UpdateTask(task *taskmodel.Task) error {
	if task.ID == 0 {
		return ErrTaskIDRequired
	}
	return nil
}
func (m *mockRepo) MarkTaskComplete(taskID int) error {
	if taskID == 0 {
		return ErrTaskIDRequired
	}
	return nil
}
func (m *mockRepo) DeleteTask(taskID int) error {
	if taskID == 0 {
		return ErrTaskIDRequired
	}
	return nil
}

func Test_CreateTask_Success(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	task := &taskmodel.Task{Title: "Task", Description: "Desc", UserID: 1}
	err := ts.CreateTask(task)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
}

func Test_CreateTask_EmptyTitle(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	task := &taskmodel.Task{Title: "", UserID: 1}

	err := ts.CreateTask(task)
	if !errors.Is(err, ErrTaskTitleRequired) {
		t.Errorf("expected ErrTaskTitleRequired, got %v", err)
	}
}

func Test_CreateTask_EmptyUserID(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	task := &taskmodel.Task{Title: "Task", UserID: 0}

	err := ts.CreateTask(task)
	if !errors.Is(err, ErrUserIDRequired) {
		t.Errorf("expected ErrUserIDRequired, got %v", err)
	}
}

func Test_GetTasksForUser_Success(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	tasks, err := ts.GetTasksForUser(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(tasks) != 1 || tasks[0].ID != 1 {
		t.Errorf("unexpected tasks: %+v", tasks)
	}
}

func Test_GetTasksForUser_InvalidUserID(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	_, err := ts.GetTasksForUser(0)

	if !errors.Is(err, ErrUserIDInvalid) {
		t.Errorf("expected ErrUserIDInvalid, got %v", err)
	}
}

func Test_UpdateTask_Success(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	task := &taskmodel.Task{ID: 1, Title: "Task", UserID: 1}

	err := ts.UpdateTask(task)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_UpdateTask_EmptyID(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	task := &taskmodel.Task{ID: 0, Title: "Task", UserID: 1}

	err := ts.UpdateTask(task)
	if !errors.Is(err, ErrTaskIDRequired) {
		t.Errorf("expected ErrTaskIDRequired, got %v", err)
	}
}

func Test_MarkTaskComplete_Success(t *testing.T) {
	ts := NewTaskService(&mockRepo{})
	err := ts.MarkTaskComplete(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_MarkTaskComplete_EmptyID(t *testing.T) {
	ts := NewTaskService(&mockRepo{})

	err := ts.MarkTaskComplete(0)
	if !errors.Is(err, ErrTaskIDRequired) {
		t.Errorf("expected ErrTaskIDRequired, got %v", err)
	}
}

func Test_DeleteTask_Success(t *testing.T) {
	ts := NewTaskService(&mockRepo{})

	err := ts.DeleteTask(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_DeleteTask_EmptyID(t *testing.T) {
	ts := NewTaskService(&mockRepo{})

	err := ts.DeleteTask(0)
	if !errors.Is(err, ErrTaskIDRequired) {
		t.Errorf("expected ErrTaskIDRequired, got %v", err)
	}
}
