package taskhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	taskmodel "assignment8/models/task"
)

type mockTaskService struct {
	CreateTaskFunc       func(task *taskmodel.Task) error
	GetTasksForUserFunc  func(userID int) ([]taskmodel.Task, error)
	UpdateTaskFunc       func(task *taskmodel.Task) error
	MarkTaskCompleteFunc func(taskID int) error
	DeleteTaskFunc       func(taskID int) error
}

func (m *mockTaskService) CreateTask(task *taskmodel.Task) error {
	if m.CreateTaskFunc != nil {
		return m.CreateTaskFunc(task)
	}
	return nil
}
func (m *mockTaskService) GetTasksForUser(userID int) ([]taskmodel.Task, error) {
	if m.GetTasksForUserFunc != nil {
		return m.GetTasksForUserFunc(userID)
	}
	return nil, nil
}
func (m *mockTaskService) UpdateTask(task *taskmodel.Task) error {
	if m.UpdateTaskFunc != nil {
		return m.UpdateTaskFunc(task)
	}
	return nil
}
func (m *mockTaskService) MarkTaskComplete(taskID int) error {
	if m.MarkTaskCompleteFunc != nil {
		return m.MarkTaskCompleteFunc(taskID)
	}
	return nil
}
func (m *mockTaskService) DeleteTask(taskID int) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(taskID)
	}
	return nil
}

func TestCreateTask_Success(t *testing.T) {
	mockSvc := &mockTaskService{
		CreateTaskFunc: func(task *taskmodel.Task) error {
			task.ID = 1
			return nil
		},
	}
	h := NewTaskHandler(mockSvc)
	task := taskmodel.Task{Title: "Task", Description: "Desc", UserID: 1}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
	var got taskmodel.Task
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.Title != "Task" {
		t.Errorf("expected title Task, got %s", got.Title)
	}
}

func TestCreateTask_InvalidBody(t *testing.T) {
	mockSvc := &mockTaskService{}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte("notjson")))
	w := httptest.NewRecorder()

	h.CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateTask_ServiceError(t *testing.T) {
	mockSvc := &mockTaskService{
		CreateTaskFunc: func(task *taskmodel.Task) error {
			return errors.New("service error")
		},
	}
	h := NewTaskHandler(mockSvc)
	task := taskmodel.Task{Title: "Task", Description: "Desc", UserID: 1}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetUserTasks_Success(t *testing.T) {
	mockSvc := &mockTaskService{
		GetTasksForUserFunc: func(userID int) ([]taskmodel.Task, error) {
			return []taskmodel.Task{{ID: 1, Title: "Task", UserID: userID}}, nil
		},
	}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.GetUserTasks(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	var got []taskmodel.Task
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("unexpected tasks: %+v", got)
	}
}

func TestGetUserTasks_InvalidID(t *testing.T) {
	mockSvc := &mockTaskService{}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.GetUserTasks(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetUserTasks_ServiceError(t *testing.T) {
	mockSvc := &mockTaskService{
		GetTasksForUserFunc: func(userID int) ([]taskmodel.Task, error) {
			return nil, errors.New("service error")
		},
	}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.GetUserTasks(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	mockSvc := &mockTaskService{
		UpdateTaskFunc: func(task *taskmodel.Task) error {
			return nil
		},
	}
	h := NewTaskHandler(mockSvc)
	task := taskmodel.Task{ID: 1, Title: "Task", UserID: 1}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPut, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.UpdateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_InvalidBody(t *testing.T) {
	mockSvc := &mockTaskService{}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/tasks", bytes.NewReader([]byte("notjson")))
	w := httptest.NewRecorder()

	h.UpdateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_ServiceError(t *testing.T) {
	mockSvc := &mockTaskService{
		UpdateTaskFunc: func(task *taskmodel.Task) error {
			return errors.New("service error")
		},
	}
	h := NewTaskHandler(mockSvc)
	task := taskmodel.Task{ID: 1, Title: "Task", UserID: 1}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPut, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.UpdateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	mockSvc := &mockTaskService{
		DeleteTaskFunc: func(taskID int) error {
			return nil
		},
	}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_InvalidID(t *testing.T) {
	mockSvc := &mockTaskService{}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_ServiceError(t *testing.T) {
	mockSvc := &mockTaskService{
		DeleteTaskFunc: func(taskID int) error {
			return errors.New("service error")
		},
	}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestMarkTaskComplete_Success(t *testing.T) {
	mockSvc := &mockTaskService{
		MarkTaskCompleteFunc: func(taskID int) error {
			return nil
		},
	}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1/complete", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.MarkTaskComplete(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMarkTaskComplete_InvalidID(t *testing.T) {
	mockSvc := &mockTaskService{}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/tasks/abc/complete", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.MarkTaskComplete(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestMarkTaskComplete_ServiceError(t *testing.T) {
	mockSvc := &mockTaskService{
		MarkTaskCompleteFunc: func(taskID int) error {
			return errors.New("service error")
		},
	}
	h := NewTaskHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1/complete", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.MarkTaskComplete(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
