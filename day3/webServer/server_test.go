package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() *server {
	return newServer()
}

func TestGreetHandler(t *testing.T) {
	srv := newTestServer()

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

func TestGetUserHandler(t *testing.T) {
	srv := newTestServer()
	srv.userCache[1] = user{Name: "Aman"}

	req := httptest.NewRequest(http.MethodGet, "/users/1", http.NoBody)
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	var u user

	err := json.NewDecoder(rec.Body).Decode(&u)
	if err != nil || u.Name != "Aman" {
		t.Errorf("Unexpected body: %s", rec.Body.String())
	}
}

func TestDeleteUserHandler(t *testing.T) {
	srv := newTestServer()
	srv.userCache[1] = user{Name: "Aman"}

	req := httptest.NewRequest(http.MethodDelete, "/users/1", http.NoBody)

	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status %d and got status %d", http.StatusNoContent, rec.Code)
	}

	if _, ok := srv.userCache[1]; ok {
		t.Errorf("Expected user to be deleted")
	}
}

func TestPostUserHandle(t *testing.T) {
	srv := newTestServer()

	body := bytes.NewBufferString(`{"name":"AMAN"}`)

	req := httptest.NewRequest(http.MethodPost, "/users", body)
	rec := httptest.NewRecorder()

	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status %d and got status %d", http.StatusCreated, rec.Code)
	}

	if _, ok := srv.userCache[1]; !ok {
		t.Errorf("Expected user to be deleted")
	}
}

func TestUpdateHandler(t *testing.T) {
	srv := newTestServer()

	srv.userCache[1] = user{Name: "OldName"}

	newUser := `{"name":"UpdatedName"}`
	body := bytes.NewBufferString(newUser)

	req := httptest.NewRequest(http.MethodPut, "/users/1", body)
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	srv.updateHandler(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	updated := srv.userCache[1]
	if updated.Name != "UpdatedName" {
		t.Errorf("Expected user name to be UpdatedName, got %s", updated.Name)
	}
}
