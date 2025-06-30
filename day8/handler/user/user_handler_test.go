package userhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	usermodel "assignment8/models/user"
)

type mockUserService struct {
	CreateUserFunc func(user *usermodel.User) error
	GetUserFunc    func(id int) (*usermodel.User, error)
}

func (m *mockUserService) CreateUser(user *usermodel.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return nil
}
func (m *mockUserService) GetUser(id int) (*usermodel.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(id)
	}
	return nil, nil
}

func TestCreateUserHandler_Success(t *testing.T) {
	mockSvc := &mockUserService{
		CreateUserFunc: func(user *usermodel.User) error {
			user.ID = 1
			return nil
		},
	}
	h := NewUserHandler(mockSvc)
	user := usermodel.User{Name: "Aman"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateUserHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
	var got usermodel.User
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.Name != "Aman" {
		t.Errorf("expected name Aman, got %s", got.Name)
	}
}

func TestCreateUserHandler_InvalidBody(t *testing.T) {
	mockSvc := &mockUserService{}
	h := NewUserHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("notjson")))
	w := httptest.NewRecorder()

	h.CreateUserHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateUserHandler_ServiceError(t *testing.T) {
	mockSvc := &mockUserService{
		CreateUserFunc: func(user *usermodel.User) error {
			return errors.New("service error")
		},
	}
	h := NewUserHandler(mockSvc)
	user := usermodel.User{Name: "Aman"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateUserHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestGetUserHandler_Success(t *testing.T) {
	mockSvc := &mockUserService{
		GetUserFunc: func(id int) (*usermodel.User, error) {
			return &usermodel.User{ID: 2, Name: "Aman"}, nil
		},
	}
	h := NewUserHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
	req.SetPathValue("id", "2")
	w := httptest.NewRecorder()

	h.GetUserHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	var got usermodel.User
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.ID != 2 || got.Name != "Aman" {
		t.Errorf("unexpected user: %+v", got)
	}
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	mockSvc := &mockUserService{}
	h := NewUserHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.GetUserHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetUserHandler_NotFound(t *testing.T) {
	mockSvc := &mockUserService{
		GetUserFunc: func(id int) (*usermodel.User, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewUserHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	req.SetPathValue("id", "99")
	w := httptest.NewRecorder()

	h.GetUserHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}
