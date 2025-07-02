package taskrepository

import (
	taskmodel "assignment8/models/task"

	"gofr.dev/pkg/gofr"
)

type TaskRepository struct {
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (*TaskRepository) CreateTask(c *gofr.Context, task *taskmodel.Task) error {
	db := c.SQL
	query := "INSERT INTO tasks (title, description, user_id, is_completed) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, task.Title, task.Description, task.UserID, false)

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

func (*TaskRepository) GetTasksByUserID(c *gofr.Context, userID int) ([]taskmodel.Task, error) {
	db := c.SQL
	query := "SELECT id, title, description, user_id, is_completed FROM tasks WHERE user_id = ?"
	rows, err := db.Query(query, userID)

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

func (*TaskRepository) UpdateTask(c *gofr.Context, task *taskmodel.Task) error {
	db := c.SQL
	query := "UPDATE tasks SET title=?, description=?, is_completed=? WHERE id=?"
	_, err := db.Exec(query, task.Title, task.Description, task.IsCompleted, task.ID)

	return err
}

func (*TaskRepository) MarkTaskComplete(c *gofr.Context, taskID int) error {
	db := c.SQL
	query := "UPDATE tasks SET is_completed=true WHERE id=?"
	_, err := db.Exec(query, taskID)

	return err
}

func (*TaskRepository) DeleteTask(c *gofr.Context, taskID int) error {
	db := c.SQL
	query := "DELETE FROM tasks WHERE id=?"
	_, err := db.Exec(query, taskID)

	return err
}
