package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func newTestServer(t *testing.T) *server {
	dsn := "root:rootpassword@tcp(localhost:3306)/mydb_test"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean test DB: %v", err)
	}

	s := &server{
		mux: http.NewServeMux(),
		db:  db,
	}
	s.routes()

	t.Cleanup(func() {
		db.Close()
	})

	return s
}

func insertUser(t *testing.T, db *sql.DB, name string) int {
	res, err := db.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func TestGreetHandler(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	expected := `{"Message":"Hello World using JSON in GOLANG"}`
	if strings.TrimSpace(rec.Body.String()) != expected {
		t.Errorf("Unexpected body: %s", rec.Body.String())
	}
}

func TestPostUserHandler(t *testing.T) {
	srv := newTestServer(t)

	body := bytes.NewBufferString(`{"name":"AMAN"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var count int
	err := srv.db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", "AMAN").Scan(&count)
	if err != nil {
		t.Fatalf("DB check failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

func TestGetUserHandler(t *testing.T) {
	srv := newTestServer(t)
	id := insertUser(t, srv.db, "Aman")

	req := httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(id), http.NoBody)
	req.SetPathValue("id", strconv.Itoa(id))
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, rec.Code)
	}

	var u user
	err := json.NewDecoder(rec.Body).Decode(&u)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}
	if u.Name != "Aman" {
		t.Errorf("Expected Aman, got %s", u.Name)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	srv := newTestServer(t)
	id := insertUser(t, srv.db, "Aman")

	req := httptest.NewRequest(http.MethodDelete, "/users/"+strconv.Itoa(id), http.NoBody)
	req.SetPathValue("id", strconv.Itoa(id))
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected %d, got %d", http.StatusNoContent, rec.Code)
	}

	var count int
	err := srv.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check DB: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 users, got %d", count)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	srv := newTestServer(t)
	id := insertUser(t, srv.db, "OldName")

	body := bytes.NewBufferString(`{"name":"NewName"}`)
	req := httptest.NewRequest(http.MethodPut, "/users/"+strconv.Itoa(id), body)
	req.SetPathValue("id", strconv.Itoa(id))
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Expected %d, got %d", http.StatusAccepted, rec.Code)
	}

	var name string
	err := srv.db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		t.Fatalf("DB check failed: %v", err)
	}
	if name != "NewName" {
		t.Errorf("Expected name to be NewName, got %s", name)
	}
}

func TestInvalidIDReturns400(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", http.NoBody)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestMissingUserReturns404(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/users/9999", http.NoBody)
	req.SetPathValue("id", "9999")
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, rec.Code)
	}
}
