package userrepository

import (
	"database/sql"
	"errors"

	usermodel "assignment8/models/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *usermodel.User) error {
	query := "INSERT INTO users (name) VALUES (?)"
	result, err := r.db.Exec(query, user.Name)

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

func (r *UserRepository) GetUser(id int) (*usermodel.User, error) {
	query := "SELECT id, name FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

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
