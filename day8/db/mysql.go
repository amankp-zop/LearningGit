// db/mysql.go
package db

import (
	"database/sql"
	"fmt"

	// MySQL driver is registered via its init() function.
	// This blank import is necessary for `sql.Open("mysql", ...)` to work.
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	dsn := "root:rootpassword@tcp(127.0.0.1:3306)/taskmanager"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	fmt.Println("Connected to MySQL successfully!")

	return db, nil
}
