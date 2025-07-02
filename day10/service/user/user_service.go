package userservice

import (
	"errors"

	usermodel "assignment8/models/user"

	"gofr.dev/pkg/gofr"
)

var (
	ErrUserNameRequired = errors.New("user name is required")
	ErrUserNotFound     = errors.New("user not found")
)

type UserService struct {
	repository UserRepositoryInterface
}

func NewUserService(repository UserRepositoryInterface) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateUser(c *gofr.Context, user *usermodel.User) error {
	if user.Name == "" {
		return ErrUserNameRequired
	}

	return s.repository.CreateUser(c, user)
}

func (s *UserService) GetUser(c *gofr.Context, id int) (*usermodel.User, error) {
	user, err := s.repository.GetUser(c, id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
