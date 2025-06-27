package service

import (
	"errors"
	"taskmanager/models"
	"taskmanager/repository"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUser(id int) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *models.User) error {
	if user.Email == "" || user.Name == "" {
		return errors.New("name and emalil cannot be empty")
	}
	return s.userRepo.CreateUser(user)
}

func (s *userService) GetUser(id int) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
