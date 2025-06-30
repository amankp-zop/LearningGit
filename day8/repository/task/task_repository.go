package taskrepository

import (
	"database/sql"

	taskmodel "assignment8/models/task"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *taskmodel.Task) error {
	query := "INSERT INTO tasks (title, description, user_id, is_completed) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, task.Title, task.Description, task.UserID, false)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	task.ID = int(id)
	task.IsCompleted = false

	return nil
}

func (r *TaskRepository) GetTasksByUserID(userID int) ([]taskmodel.Task, error) {
	query := "SELECT id, title, description, user_id, is_completed FROM tasks WHERE user_id = ?"
	rows, err := r.db.Query(query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []taskmodel.Task

	for rows.Next() {
		var task taskmodel.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.UserID, &task.IsCompleted)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(task *taskmodel.Task) error {
	query := "UPDATE tasks SET title=?, description=?, is_completed=? WHERE id=?"
	_, err := r.db.Exec(query, task.Title, task.Description, task.IsCompleted, task.ID)

	return err
}

func (r *TaskRepository) MarkTaskComplete(taskID int) error {
	query := "UPDATE tasks SET is_completed=true WHERE id=?"
	_, err := r.db.Exec(query, taskID)

	return err
}

func (r *TaskRepository) DeleteTask(taskID int) error {
	query := "DELETE FROM tasks WHERE id=?"
	_, err := r.db.Exec(query, taskID)

	return err
}
