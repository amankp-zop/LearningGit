package userservice

import (
	"errors"

	usermodel "assignment8/models/user"
)

var (
	ErrUserNameRequired = errors.New("user name is required")
	ErrUserNotFound     = errors.New("user not found")
)

type UserRepositoryInterface interface {
	CreateUser(user *usermodel.User) error
	GetUser(id int) (*usermodel.User, error)
}

type UserService struct {
	repository UserRepositoryInterface
}

func NewUserService(repository UserRepositoryInterface) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateUser(user *usermodel.User) error {
	if user.Name == "" {
		return ErrUserNameRequired
	}

	return s.repository.CreateUser(user)
}

func (s *UserService) GetUser(id int) (*usermodel.User, error) {
	user, err := s.repository.GetUser(id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
