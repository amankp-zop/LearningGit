package repository

import (
	"database/sql"
	"taskmanager/models"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetTasksByUserID(id int) ([]models.Task, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	query := "INSERT INTO tasks (title, description, user_id) VALUES (?, ?, ?)"
	data, err := r.db.Exec(query, task.Title, task.Description, task.UserId)
	if err != nil {
		return err
	}

	id, err := data.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = int(id)
	return nil
}

func (r *taskRepository) GetTasksByUserID(id int) ([]models.Task, error) {
	query := "SELECT id,title,description,user_id FROM tasks where user_id=?"

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.UserId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
