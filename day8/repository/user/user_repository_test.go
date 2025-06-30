package userrepository_test

import (
	"errors"
	"testing"

	usermodel "assignment8/models/user"
	userrepository "assignment8/repository/user"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := userrepository.NewUserRepository(db)
	user := &usermodel.User{Name: "Aman"}

	mock.ExpectExec("INSERT INTO users (name) VALUES (?)").WithArgs(user.Name).WillReturnResult(sqlmock.NewResult(5, 1))

	err = repo.CreateUser(user)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if user.ID != 5 {
		t.Errorf("expected user ID 5, got %d", user.ID)
	}
}

func TestCreateUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := userrepository.NewUserRepository(db)
	user := &usermodel.User{Name: "Aman"}

	mock.ExpectExec("INSERT INTO users (name) VALUES (?)").WithArgs(user.Name).WillReturnError(errors.New("insert error"))

	err = repo.CreateUser(user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := userrepository.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "Aman")
	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").WithArgs(2).WillReturnRows(rows)

	user, err := repo.GetUser(2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if user == nil || user.ID != 2 || user.Name != "Aman" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := userrepository.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").WithArgs(3).WillReturnRows(rows)

	user, err := repo.GetUser(3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}

func TestGetUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := userrepository.NewUserRepository(db)

	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").WithArgs(4).WillReturnError(errors.New("query error"))

	user, err := repo.GetUser(4)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Errorf("expected nil user, got %v", user)
	}
}
