package repository

import (
	"database/sql"
	"errors"
	"taskmanager/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	result, err := r.db.Exec(query, user.Name, user.Email)
	if err != nil {
		return err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(insertedID)
	return nil
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	query := "SELECT id, name, email FROM users WHERE id = ?"
	result := r.db.QueryRow(query, id)

	var user models.User
	err := result.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
