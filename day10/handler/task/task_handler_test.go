package taskhandler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/http/response"

	taskmodel "assignment8/models/task"

	gofrHttp "gofr.dev/pkg/gofr/http"
)

func TestCreateTask(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)

	tests := []struct {
		name             string
		requestBody      string
		expectedTask     taskmodel.Task
		expectedResponse gofrResponse
		ifMock           bool
		mockError        error
	}{
		{
			name:        "Successful Create Task",
			requestBody: `{"id":1,"title":"Test Task","description":"desc","user_id":1,"is_completed":false}`,
			expectedTask: taskmodel.Task{
				ID:          1,
				Title:       "Test Task",
				Description: "desc",
				UserID:      1,
				IsCompleted: false,
			},
			expectedResponse: gofrResponse{
				result: response.Raw{Data: taskmodel.Task{ID: 1, Title: "Test Task", Description: "desc", UserID: 1, IsCompleted: false}},
				err:    nil,
			},
			ifMock:    true,
			mockError: nil,
		},
		{
			name:         "Failed Binding",
			requestBody:  `{"id":1,"title":}`,
			expectedTask: taskmodel.Task{},
			expectedResponse: gofrResponse{
				result: nil,
				// Accept any error, will check error type in assertion
				err: nil,
			},
			ifMock:    false,
			mockError: nil,
		},
		{
			name:        "Service Error",
			requestBody: `{"id":2,"title":"Error Task","description":"desc2","user_id":2,"is_completed":true}`,
			expectedTask: taskmodel.Task{
				ID:          2,
				Title:       "Error Task",
				Description: "desc2",
				UserID:      2,
				IsCompleted: true,
			},
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("service error"),
			},
			ifMock:    true,
			mockError: errors.New("service error"),
		},
		{
			name:         "Empty Body",
			requestBody:  ``,
			expectedTask: taskmodel.Task{},
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("Kuch nhi hai body me bro"),
			},
			ifMock:    false,
			mockError: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := NewMockTaskServicePort(ctrl)
			handler := NewTaskHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			request := gofrHttp.NewRequest(req)

			ctx := &gofr.Context{
				Context:   context.Background(),
				Request:   request,
				Container: mockContainer,
			}

			if tt.ifMock {
				mockService.EXPECT().CreateTask(ctx, &tt.expectedTask).Return(tt.mockError)
			}

			val, err := handler.CreateTask(ctx)

			if tt.name == "Failed Binding" {
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid character")
			} else if tt.name == "Empty Body" {
				assert.Nil(t, val)
				assert.Error(t, err)
			} else {
				response := gofrResponse{val, err}
				assert.Equal(t, tt.expectedResponse, response, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}

func TestGetUserTasks(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)

	tests := []struct {
		name             string
		idParam          string
		expectedID       int
		mockTasks        []taskmodel.Task
		mockError        error
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{
			name:       "Successful Get User Tasks",
			idParam:    "1",
			expectedID: 1,
			mockTasks:  []taskmodel.Task{{ID: 1, UserID: 1, Title: "Task1"}},
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: response.Raw{Data: []taskmodel.Task{{ID: 1, UserID: 1, Title: "Task1"}}},
				err:    nil,
			},
			ifMock: true,
		},
		{
			name:       "Invalid ID Param",
			idParam:    "abc",
			expectedID: 0,
			mockTasks:  nil,
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: nil,
				err:    &strconv.NumError{Func: "Atoi", Num: "abc", Err: errors.New("invalid syntax")},
			},
			ifMock: false,
		},
		{
			name:       "Service Error",
			idParam:    "2",
			expectedID: 2,
			mockTasks:  nil,
			mockError:  errors.New("service error"),
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("service error"),
			},
			ifMock: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := NewMockTaskServicePort(ctrl)
			handler := NewTaskHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/task/user/"+tt.idParam, http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tt.idParam})
			request := gofrHttp.NewRequest(req)

			ctx := &gofr.Context{
				Context:   context.Background(),
				Request:   request,
				Container: mockContainer,
			}

			if tt.ifMock {
				mockService.EXPECT().GetTasksForUser(ctx, tt.expectedID).Return(tt.mockTasks, tt.mockError)
			}

			val, err := handler.GetUserTasks(ctx)

			if tt.name == "Invalid ID Param" {
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid syntax")
			} else {
				response := gofrResponse{val, err}
				assert.Equal(t, tt.expectedResponse, response, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)

	tests := []struct {
		name             string
		requestBody      string
		expectedTask     taskmodel.Task
		expectedResponse gofrResponse
		ifMock           bool
		mockError        error
	}{
		{
			name:        "Successful Update Task",
			requestBody: `{"id":1,"title":"Updated Task","description":"desc","user_id":1,"is_completed":true}`,
			expectedTask: taskmodel.Task{
				ID:          1,
				Title:       "Updated Task",
				Description: "desc",
				UserID:      1,
				IsCompleted: true,
			},
			expectedResponse: gofrResponse{
				result: response.Raw{Data: taskmodel.Task{ID: 1, Title: "Updated Task", Description: "desc", UserID: 1, IsCompleted: true}},
				err:    nil,
			},
			ifMock:    true,
			mockError: nil,
		},
		{
			name:         "Failed Binding",
			requestBody:  `{"id":1,"title":}`,
			expectedTask: taskmodel.Task{},
			expectedResponse: gofrResponse{
				result: nil,
				// Accept any error, will check error type in assertion
				err: nil,
			},
			ifMock:    false,
			mockError: nil,
		},
		{
			name:        "Service Error",
			requestBody: `{"id":2,"title":"Error Task","description":"desc2","user_id":2,"is_completed":false}`,
			expectedTask: taskmodel.Task{
				ID:          2,
				Title:       "Error Task",
				Description: "desc2",
				UserID:      2,
				IsCompleted: false,
			},
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("service error"),
			},
			ifMock:    true,
			mockError: errors.New("service error"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := NewMockTaskServicePort(ctrl)
			handler := NewTaskHandler(mockService)

			req := httptest.NewRequest(http.MethodPut, "/task", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			request := gofrHttp.NewRequest(req)

			ctx := &gofr.Context{
				Context:   context.Background(),
				Request:   request,
				Container: mockContainer,
			}

			if tt.ifMock {
				mockService.EXPECT().UpdateTask(ctx, &tt.expectedTask).Return(tt.mockError)
			}

			val, err := handler.UpdateTask(ctx)

			if tt.name == "Failed Binding" {
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid character")
			} else {
				response := gofrResponse{val, err}
				assert.Equal(t, tt.expectedResponse, response, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)

	tests := []struct {
		name             string
		idParam          string
		expectedID       int
		mockError        error
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{
			name:       "Successful Delete Task",
			idParam:    "1",
			expectedID: 1,
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: response.Raw{Data: "Task deleted"},
				err:    nil,
			},
			ifMock: true,
		},
		{
			name:       "Invalid ID Param",
			idParam:    "abc",
			expectedID: 0,
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: nil,
				err:    &strconv.NumError{Func: "Atoi", Num: "abc", Err: errors.New("invalid syntax")},
			},
			ifMock: false,
		},
		{
			name:       "Service Error",
			idParam:    "2",
			expectedID: 2,
			mockError:  errors.New("service error"),
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("service error"),
			},
			ifMock: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := NewMockTaskServicePort(ctrl)
			handler := NewTaskHandler(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/task/"+tt.idParam, http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tt.idParam})
			request := gofrHttp.NewRequest(req)

			ctx := &gofr.Context{
				Context:   context.Background(),
				Request:   request,
				Container: mockContainer,
			}

			if tt.ifMock {
				mockService.EXPECT().DeleteTask(ctx, tt.expectedID).Return(tt.mockError)
			}

			val, err := handler.DeleteTask(ctx)

			if tt.name == "Invalid ID Param" {
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid syntax")
			} else {
				response := gofrResponse{val, err}
				assert.Equal(t, tt.expectedResponse, response, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}

func TestMarkTaskComplete(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)

	tests := []struct {
		name             string
		idParam          string
		expectedID       int
		mockError        error
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{
			name:       "Successful Mark Task Complete",
			idParam:    "1",
			expectedID: 1,
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: response.Raw{Data: "Task marked as complete"},
				err:    nil,
			},
			ifMock: true,
		},
		{
			name:       "Invalid ID Param",
			idParam:    "abc",
			expectedID: 0,
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: nil,
				err:    &strconv.NumError{Func: "Atoi", Num: "abc", Err: errors.New("invalid syntax")},
			},
			ifMock: false,
		},
		{
			name:       "Service Error",
			idParam:    "2",
			expectedID: 2,
			mockError:  errors.New("service error"),
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("service error"),
			},
			ifMock: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := NewMockTaskServicePort(ctrl)
			handler := NewTaskHandler(mockService)

			req := httptest.NewRequest(http.MethodPatch, "/task/"+tt.idParam+"/complete", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tt.idParam})
			request := gofrHttp.NewRequest(req)

			ctx := &gofr.Context{
				Context:   context.Background(),
				Request:   request,
				Container: mockContainer,
			}

			if tt.ifMock {
				mockService.EXPECT().MarkTaskComplete(ctx, tt.expectedID).Return(tt.mockError)
			}

			val, err := handler.MarkTaskComplete(ctx)

			if tt.name == "Invalid ID Param" {
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid syntax")
			} else {
				response := gofrResponse{val, err}
				assert.Equal(t, tt.expectedResponse, response, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}
