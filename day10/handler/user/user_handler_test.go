package userhandler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	usermodel "assignment8/models/user"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/http/response"
)

func TestCreateUserHandler(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)

	tests := []struct {
		name             string
		requestBody      string
		expectedUser     usermodel.User
		expectedResponse gofrResponse
		ifMock           bool
		mockError        error
	}{
		{
			name:        "Successful Create User",
			requestBody: `{"id":1,"name":"Tester"}`,
			expectedUser: usermodel.User{
				ID:   1,
				Name: "Tester",
			},
			expectedResponse: gofrResponse{
				result: usermodel.User{ID: 1, Name: "Tester"},
				err:    nil,
			},
			ifMock:    true,
			mockError: nil,
		},
		{
			name:        "Failed Binding",
			requestBody: `{"id":1,"name":}`,
			expectedUser: usermodel.User{
				ID:   0,
				Name: "",
			},
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("bro this is wrong."),
			},
			ifMock:    false,
			mockError: nil,
		},
		{
			name:        "Service Error",
			requestBody: `{"id":2,"name":"ErrorUser"}`,
			expectedUser: usermodel.User{
				ID:   2,
				Name: "ErrorUser",
			},
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("service error"),
			},
			ifMock:    true,
			mockError: errors.New("service error"),
		},
		{
			name:        "Empty Body",
			requestBody: ``,
			expectedUser: usermodel.User{
				ID:   0,
				Name: "",
			},
			expectedResponse: gofrResponse{
				result: nil,
				err:    errors.New("Body Is Empty Bro"),
			},
			ifMock:    false,
			mockError: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := NewMockUserServiceInterface(ctrl)
			handler := NewUserHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.requestBody))

			req.Header.Set("Content-Type", "application/json")

			request := gofrHttp.NewRequest(req)

			ctx := &gofr.Context{
				Context:   context.Background(),
				Request:   request,
				Container: mockContainer,
			}

			if tt.ifMock {
				mockService.EXPECT().CreateUser(ctx, &tt.expectedUser).Return(tt.mockError)
			}

			val, err := handler.CreateUserHandler(ctx)

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

func TestGetUserHandler(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)

	tests := []struct {
		name             string
		idParam          string
		expectedID       int
		mockUser         *usermodel.User
		mockError        error
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{
			name:       "Successful Get User",
			idParam:    "1",
			expectedID: 1,
			mockUser:   &usermodel.User{ID: 1, Name: "Tester"},
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: response.Raw{Data: &usermodel.User{ID: 1, Name: "Tester"}},
				err:    nil,
			},
			ifMock: true,
		},
		{
			name:       "Invalid ID Param",
			idParam:    "abc",
			expectedID: 0,
			mockUser:   nil,
			mockError:  nil,
			expectedResponse: gofrResponse{
				result: nil,
				err:    gofrHttp.ErrorInvalidParam{Params: []string{"Give Correct Input"}},
			},
			ifMock: false,
		},
		{
			name:       "Service Error",
			idParam:    "2",
			expectedID: 2,
			mockUser:   nil,
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
			mockService := NewMockUserServiceInterface(ctrl)
			handler := NewUserHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/user/"+tt.idParam, http.NoBody)

			// Gorilla Mux to the Rescue - Yaad Rakhio isko
			req = mux.SetURLVars(req, map[string]string{"id": tt.idParam})
			request := gofrHttp.NewRequest(req)

			ctx := &gofr.Context{
				Context:   context.Background(),
				Request:   request,
				Container: mockContainer,
			}

			if tt.ifMock {
				mockService.EXPECT().GetUser(ctx, tt.expectedID).Return(tt.mockUser, tt.mockError)
			}

			val, err := handler.GetUserHandler(ctx)

			response := gofrResponse{val, err}

			if tt.name == "Invalid ID Param" {
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.IsType(t, gofrHttp.ErrorInvalidParam{}, err)
				assert.Equal(t, tt.expectedResponse.err, err)
			} else {
				assert.Equal(t, tt.expectedResponse, response, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}
