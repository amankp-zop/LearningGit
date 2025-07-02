package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createTableUser = `CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);`

func create_user_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableUser)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
