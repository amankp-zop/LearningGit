package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	dsn := "root:rootpassword@tcp(127.0.0.1:3306)/taskmanager"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to the Database and got error :%v", err)
		return
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to the Database and got error :%v", err)
		return
	}

	fmt.Println("Connected to MySQL successfully!")
}
