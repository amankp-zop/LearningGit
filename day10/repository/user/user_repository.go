package userrepository

import (
	"database/sql"
	"errors"

	usermodel "assignment8/models/user"

	"gofr.dev/pkg/gofr"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (*UserRepository) CreateUser(c *gofr.Context, user *usermodel.User) error {
	db := c.SQL
	query := "INSERT INTO users (name) VALUES (?)"
	result, err := db.Exec(query, user.Name)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)

	return nil
}

func (*UserRepository) GetUser(c *gofr.Context, id int) (*usermodel.User, error) {
	db := c.SQL
	query := "SELECT id, name FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	var user usermodel.User
	err := row.Scan(&user.ID, &user.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
