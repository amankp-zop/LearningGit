package userrepository_test

import (
	"errors"
	"testing"

	usermodel "assignment8/models/user"
	userrepository "assignment8/repository/user"

	"github.com/DATA-DOG/go-sqlmock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

func TestCreateUser_Success(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := userrepository.NewUserRepository()
	user := &usermodel.User{Name: "Aman"}

	mock.SQL.ExpectExec("INSERT INTO users (name) VALUES (?)").
		WithArgs(user.Name).
		WillReturnResult(sqlmock.NewResult(5, 1))

	err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if user.ID != 5 {
		t.Errorf("expected user ID 5, got %d", user.ID)
	}
}

func TestCreateUser_Error(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := userrepository.NewUserRepository()
	user := &usermodel.User{Name: "Aman"}

	mock.SQL.ExpectExec("INSERT INTO users (name) VALUES (?)").
		WithArgs(user.Name).
		WillReturnError(errors.New("insert error"))

	err := repo.CreateUser(ctx, user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUser_Success(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := userrepository.NewUserRepository()

	rows := mock.SQL.NewRows([]string{"id", "name"}).AddRow(2, "Aman")
	mock.SQL.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(2).
		WillReturnRows(rows)

	user, err := repo.GetUser(ctx, 2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if user == nil || user.ID != 2 || user.Name != "Aman" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := userrepository.NewUserRepository()

	rows := mock.SQL.NewRows([]string{"id", "name"})
	mock.SQL.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(3).
		WillReturnRows(rows)

	user, err := repo.GetUser(ctx, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}

func TestGetUser_Error(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer}

	repo := userrepository.NewUserRepository()

	mock.SQL.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(4).
		WillReturnError(errors.New("query error"))

	user, err := repo.GetUser(ctx, 4)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if user != nil {
		t.Errorf("expected nil user, got %v", user)
	}
}
