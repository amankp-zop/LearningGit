package taskrepository_test

import (
	taskmodel "assignment8/models/task"
	taskrepository "assignment8/repository/task"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTaskByUserID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}

	defer db.Close()

	repo := taskrepository.NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "is_completed"}).AddRow(1, "Task 1", "Description 1", 2, false).AddRow(2, "Task 2", "Description 2", 2, false)

	mock.ExpectQuery("SELECT id, title, description, user_id, is_completed FROM tasks WHERE user_id = ?").WithArgs(2).WillReturnRows(rows)

	tasks, err := repo.GetTasksByUserID(2)
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
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := taskrepository.NewTaskRepository(db)
	task := &taskmodel.Task{
		Title:       "New Task",
		Description: "A new task",
		UserID:      1,
		IsCompleted: false,
	}

	mock.ExpectExec("INSERT INTO tasks (title, description, user_id, is_completed) VALUES (?, ?, ?, ?)").WithArgs(task.Title, task.Description, task.UserID, task.IsCompleted).WillReturnResult(sqlmock.NewResult(10, 1))

	err = repo.CreateTask(task)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if task.ID != 10 {
		t.Errorf("expected task ID 10, got %d", task.ID)
	}
}

func TestCreateTask_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := taskrepository.NewTaskRepository(db)
	task := &taskmodel.Task{
		Title:       "New Task",
		Description: "A new task",
		UserID:      1,
		IsCompleted: false,
	}

	mock.ExpectExec("INSERT INTO tasks (title, description, user_id, is_completed) VALUES (?, ?, ?, ?)").WithArgs(task.Title, task.Description, task.UserID, task.IsCompleted).WillReturnError(sqlmock.ErrCancelled)

	err = repo.CreateTask(task)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := taskrepository.NewTaskRepository(db)
	task := &taskmodel.Task{
		ID:          5,
		Title:       "Updated Task",
		Description: "Updated Description",
		UserID:      3,
		IsCompleted: false,
	}

	mock.ExpectExec("UPDATE tasks SET title=?, description=?, is_completed=? WHERE id=?").
		WithArgs(task.Title, task.Description, task.IsCompleted, task.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateTask(task)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := taskrepository.NewTaskRepository(db)

	mock.ExpectExec("DELETE FROM tasks WHERE id=?").
		WithArgs(7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteTask(7)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
