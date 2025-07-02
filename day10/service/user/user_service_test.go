package userservice

import (
	"errors"
	"testing"

	usermodel "assignment8/models/user"

	gomock "go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
)

func Test_GetUser(t *testing.T) {
	ctx := &gofr.Context{}
	controller := gomock.NewController(t)
	mockStore := NewMockUserRepositoryInterface(controller)
	svc := NewUserService(mockStore)

	dummyUser := usermodel.User{
		ID:   101,
		Name: "Aman",
	}
	mockStore.EXPECT().GetUser(ctx, 101).Return(&dummyUser, nil)

	user, err := svc.GetUser(ctx, 101)
	if err != nil {
		t.Error("somthign happe")
	}

	if user.Name != "Aman" {
		t.Error("Expected aman")
	}
}

func Test_CreateUser_Success(t *testing.T) {
	ctx := &gofr.Context{}
	tests := []struct {
		name        string
		inputUser   *usermodel.User
		mockSetup   func(m *MockUserRepositoryInterface)
		expectedErr error
	}{
		{
			name:      "Success - Valid User",
			inputUser: &usermodel.User{ID: 1, Name: "Aman"},
			mockSetup: func(m *MockUserRepositoryInterface) {
				m.EXPECT().CreateUser(ctx, &usermodel.User{ID: 1, Name: "Aman"}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "Failure - Empty Name",
			inputUser:   &usermodel.User{ID: 2, Name: ""},
			mockSetup:   func(*MockUserRepositoryInterface) {}, // no call expected
			expectedErr: ErrUserNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockUserRepositoryInterface(ctrl)
			tt.mockSetup(mockRepo)

			svc := NewUserService(mockRepo)
			err := svc.CreateUser(ctx, tt.inputUser)

			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func Test_CreateUser_EmptyName(t *testing.T) {
	ctx := &gofr.Context{}
	tests := []struct {
		name         string
		inputID      int
		mockSetup    func(m *MockUserRepositoryInterface)
		expectedUser *usermodel.User
		expectedErr  error
	}{
		{
			name:    "Success - User Found",
			inputID: 101,
			mockSetup: func(m *MockUserRepositoryInterface) {
				m.EXPECT().GetUser(ctx, 101).Return(&usermodel.User{ID: 101, Name: "Aman"}, nil)
			},
			expectedUser: &usermodel.User{ID: 101, Name: "Aman"},
			expectedErr:  nil,
		},
		{
			name:    "Failure - User Not Found",
			inputID: 999,
			mockSetup: func(m *MockUserRepositoryInterface) {
				m.EXPECT().GetUser(ctx, 999).Return(nil, nil)
			},
			expectedUser: nil,
			expectedErr:  ErrUserNotFound,
		},
		{
			name:    "Failure - Repository Error",
			inputID: 123,
			mockSetup: func(m *MockUserRepositoryInterface) {
				m.EXPECT().GetUser(ctx, 123).Return(nil, errors.New("db error"))
			},
			expectedUser: nil,
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockUserRepositoryInterface(ctrl)
			tc.mockSetup(mockRepo)

			svc := NewUserService(mockRepo)

			result, err := svc.GetUser(ctx, tc.inputID)

			if result != nil && tc.expectedUser != nil {
				if result.ID != tc.expectedUser.ID || result.Name != tc.expectedUser.Name {
					t.Errorf("expected user %+v, got %+v", tc.expectedUser, result)
				}
			} else if result != tc.expectedUser {
				t.Errorf("expected user %+v, got %+v", tc.expectedUser, result)
			}

			if (err == nil) != (tc.expectedErr == nil) || (err != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
