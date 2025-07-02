package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createTableTask = `CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    user_id INT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE
);
`

func create_task_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableTask)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
