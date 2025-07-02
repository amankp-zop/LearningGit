package userservice

import (
	usermodel "assignment8/models/user"

	"gofr.dev/pkg/gofr"
)

type UserRepositoryInterface interface {
	CreateUser(c *gofr.Context, user *usermodel.User) error
	GetUser(c *gofr.Context, id int) (*usermodel.User, error)
}
