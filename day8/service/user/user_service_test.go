package userservice

import (
	"testing"

	usermodel "assignment8/models/user"
)

type mockStore struct{}

func (m *mockStore) CreateUser(user *usermodel.User) error {
	if user.Name == "" {
		return ErrUserNameRequired
	}
	user.ID = 1

	return nil
}
func (m *mockStore) GetUser(id int) (*usermodel.User, error) {
	var task *usermodel.User

	if id == 1 {
		task = &usermodel.User{
			ID:   1,
			Name: "Aman",
		}
	}
	return task, nil
}

func Test_GetUser(t *testing.T) {
	usr := NewUserService(&mockStore{})
	res, _ := usr.repository.GetUser(1)

	if res.ID != 1 {
		t.Errorf("shsh")
	}
}

func Test_CreateUser_Success(t *testing.T) {
	usr := NewUserService(&mockStore{})
	user := &usermodel.User{Name: "Aman"}

	err := usr.CreateUser(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_CreateUser_EmptyName(t *testing.T) {
	usr := NewUserService(&mockStore{})
	user := &usermodel.User{Name: ""}

	err := usr.CreateUser(user)

	if err != ErrUserNameRequired {
		t.Errorf("expected ErrUserNameRequired, got %v", err)
	}
}

func Test_GetUser_NotFound(t *testing.T) {
	usr := NewUserService(&mockStore{})

	res, err := usr.GetUser(999)

	if res != nil {
		t.Errorf("expected nil user, got %v", res)
	}

	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
