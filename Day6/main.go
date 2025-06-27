package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	Name string `json:"name"`
}

type server struct {
	mux *http.ServeMux
	db  *sql.DB
}

func newServer() *server {
	dsn := "root:rootpassword@tcp(localhost:3306)/mydb" // adjust as needed
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	s := &server{
		mux: http.NewServeMux(),
		db:  db,
	}
	s.routes()
	return s
}

func (s *server) routes() {
	s.mux.HandleFunc("/", s.greetHandler)
	s.mux.HandleFunc("POST /users", s.postHandler)
	s.mux.HandleFunc("GET /users/{id}", s.getUserHandler)
	s.mux.HandleFunc("DELETE /users/{id}", s.deleteUserHandler)
	s.mux.HandleFunc("PUT /users/{id}", s.updateHandler)
}

func main() {
	s := newServer()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      s.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server starts at 8080")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}

func (*server) greetHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := map[string]string{
		"Message": "Hello World using JSON in GOLANG",
	}

	msg, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(msg)
}

func (s *server) postHandler(w http.ResponseWriter, r *http.Request) {
	var usr user
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil || usr.Name == "" {
		http.Error(w, "Invalid user payload", http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("INSERT INTO users (name) VALUES (?)", usr.Name)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var usr user
	err = s.db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&usr.Name)
	if err == sql.ErrNoRows {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	res, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) updateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var usr user
	err = json.NewDecoder(r.Body).Decode(&usr)
	if err != nil || usr.Name == "" {
		http.Error(w, "Invalid user payload", http.StatusBadRequest)
		return
	}

	res, err := s.db.Exec("UPDATE users SET name = ? WHERE id = ?", usr.Name, id)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
