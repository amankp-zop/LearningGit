package taskrepository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"

	taskmodel "assignment8/models/task"
	taskrepository "assignment8/repository/task"
)

func TestGetTaskByUserID(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := taskrepository.NewTaskRepository()

	rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "is_completed"}).
		AddRow(1, "Task 1", "Description 1", 2, false).
		AddRow(2, "Task 2", "Description 2", 2, false)

	mock.SQL.ExpectQuery("SELECT id, title, description, user_id, is_completed FROM tasks WHERE user_id = ?").
		WithArgs(2).
		WillReturnRows(rows)

	tasks, err := repo.GetTasksByUserID(ctx, 2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != 1 || tasks[0].Title != "Task 1" || tasks[0].Description != "Description 1" || tasks[0].UserID != 2 {
		t.Errorf("unexpected first task: %+v", tasks[0])
	}

	if tasks[1].ID != 2 || tasks[1].Title != "Task 2" || tasks[1].Description != "Description 2" || tasks[1].UserID != 2 {
		t.Errorf("unexpected second task: %+v", tasks[1])
	}
}

func TestCreateTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := taskrepository.NewTaskRepository()
	task := &taskmodel.Task{
		Title:       "New Task",
		Description: "A new task",
		UserID:      1,
		IsCompleted: false,
	}

	mock.SQL.ExpectExec("INSERT INTO tasks (title, description, user_id, is_completed) VALUES (?, ?, ?, ?)").
		WithArgs(task.Title, task.Description, task.UserID, task.IsCompleted).
		WillReturnResult(sqlmock.NewResult(10, 1))

	err := repo.CreateTask(ctx, task)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if task.ID != 10 {
		t.Errorf("expected task ID 10, got %d", task.ID)
	}
}

func TestCreateTask_Error(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := taskrepository.NewTaskRepository()
	task := &taskmodel.Task{
		Title:       "New Task",
		Description: "A new task",
		UserID:      1,
		IsCompleted: false,
	}

	mock.SQL.ExpectExec("INSERT INTO tasks (title, description, user_id, is_completed) VALUES (?, ?, ?, ?)").
		WithArgs(task.Title, task.Description, task.UserID, task.IsCompleted).
		WillReturnError(sqlmock.ErrCancelled)

	err := repo.CreateTask(ctx, task)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := taskrepository.NewTaskRepository()
	task := &taskmodel.Task{
		ID:          5,
		Title:       "Updated Task",
		Description: "Updated Description",
		UserID:      3,
		IsCompleted: false,
	}

	mock.SQL.ExpectExec("UPDATE tasks SET title=?, description=?, is_completed=? WHERE id=?").
		WithArgs(task.Title, task.Description, task.IsCompleted, task.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestDeleteTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := taskrepository.NewTaskRepository()

	mock.SQL.ExpectExec("DELETE FROM tasks WHERE id=?").
		WithArgs(7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteTask(ctx, 7)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
