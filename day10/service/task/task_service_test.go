package taskservice

import (
	"errors"
	"testing"

	taskmodel "assignment8/models/task"

	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
)

func Test_CreateTask(t *testing.T) {
	ctx := &gofr.Context{}
	tests := []struct {
		name        string
		inputTask   *taskmodel.Task
		expectedErr error
		mockMethod  func(*MockTaskRepositoryPort)
	}{
		{
			name: "CreateTask-Success",
			inputTask: &taskmodel.Task{
				ID:          1,
				Title:       "task1",
				Description: "Desc1",
				IsCompleted: false,
				UserID:      1,
			},
			expectedErr: nil,
			mockMethod: func(m *MockTaskRepositoryPort) {
				m.EXPECT().CreateTask(ctx, &taskmodel.Task{
					ID:          1,
					Title:       "task1",
					Description: "Desc1",
					UserID:      1,
					IsCompleted: false,
				})
			},
		},
		{
			name: "CreateTask-EmptyTitle",
			inputTask: &taskmodel.Task{
				ID:          1,
				Title:       "",
				Description: "Descrip",
				IsCompleted: false,
				UserID:      1,
			},
			expectedErr: ErrTaskTitleRequired,
			mockMethod:  func(*MockTaskRepositoryPort) {},
		},
		{
			name: "CreateTask-InvalidUserID",
			inputTask: &taskmodel.Task{
				ID:          1,
				Title:       "Title1",
				Description: "Desc1",
				IsCompleted: false,
				UserID:      0,
			},
			expectedErr: ErrUserIDRequired,
			mockMethod:  func(*MockTaskRepositoryPort) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockStore := NewMockTaskRepositoryPort(controller)
			tt.mockMethod(mockStore)

			svc := NewTaskService(mockStore)
			err := svc.CreateTask(ctx, tt.inputTask)

			if !errors.Is(err, tt.expectedErr) {
				t.Error("EXpected Error")
			}
		})
	}
}

func Test_GetTasksForUser(t *testing.T) {
	ctx := &gofr.Context{}
	tests := []struct {
		name        string
		userID      int
		mockSetup   func(m *MockTaskRepositoryPort)
		expectedLen int
		expectedErr error
	}{
		{
			name:   "ValidUserID",
			userID: 1,
			mockSetup: func(m *MockTaskRepositoryPort) {
				m.EXPECT().GetTasksByUserID(ctx, 1).Return([]taskmodel.Task{
					{ID: 1, Title: "Task", Description: "Desc", UserID: 1},
				}, nil)
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name:   "InvalidUserID",
			userID: 0,
			mockSetup: func(*MockTaskRepositoryPort) {

			},
			expectedLen: 0,
			expectedErr: ErrUserIDInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockTaskRepositoryPort(ctrl)
			tt.mockSetup(mockRepo)

			svc := NewTaskService(mockRepo)

			tasks, err := svc.GetTasksForUser(ctx, tt.userID)

			if len(tasks) != tt.expectedLen {
				t.Errorf("expected task length %d, got %d", tt.expectedLen, len(tasks))
			}

			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func Test_UpdateTask(t *testing.T) {
	ctx := &gofr.Context{}
	tests := []struct {
		name        string
		input       *taskmodel.Task
		mockSetup   func(m *MockTaskRepositoryPort)
		expectedErr error
	}{
		{
			name: "Valid Task",
			input: &taskmodel.Task{
				ID:     1,
				Title:  "Task",
				UserID: 1,
			},
			mockSetup: func(m *MockTaskRepositoryPort) {
				m.EXPECT().UpdateTask(ctx, &taskmodel.Task{ID: 1, Title: "Task", UserID: 1}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Empty ID",
			input: &taskmodel.Task{
				ID:     0,
				Title:  "Task",
				UserID: 1,
			},
			mockSetup:   func(*MockTaskRepositoryPort) {},
			expectedErr: ErrTaskIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockTaskRepositoryPort(ctrl)
			tt.mockSetup(mockRepo)

			svc := NewTaskService(mockRepo)
			err := svc.UpdateTask(ctx, tt.input)

			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func Test_MarkTaskComplete(t *testing.T) {
	ctx := &gofr.Context{}
	tests := []struct {
		name        string
		taskID      int
		mockSetup   func(m *MockTaskRepositoryPort)
		expectedErr error
	}{
		{
			name:   "Valid Task ID",
			taskID: 1,
			mockSetup: func(m *MockTaskRepositoryPort) {
				m.EXPECT().MarkTaskComplete(ctx, 1).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "Invalid Task ID",
			taskID:      0,
			mockSetup:   func(*MockTaskRepositoryPort) {},
			expectedErr: ErrTaskIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockTaskRepositoryPort(ctrl)
			tt.mockSetup(mockRepo)

			svc := NewTaskService(mockRepo)
			err := svc.MarkTaskComplete(ctx, tt.taskID)

			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func Test_DeleteTask(t *testing.T) {
	ctx := &gofr.Context{}
	tests := []struct {
		name        string
		taskID      int
		mockSetup   func(m *MockTaskRepositoryPort)
		expectedErr error
	}{
		{
			name:   "Valid Task ID",
			taskID: 1,
			mockSetup: func(m *MockTaskRepositoryPort) {
				m.EXPECT().DeleteTask(ctx, 1).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "Invalid Task ID",
			taskID:      0,
			mockSetup:   func(*MockTaskRepositoryPort) {},
			expectedErr: ErrTaskIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockTaskRepositoryPort(ctrl)
			tt.mockSetup(mockRepo)

			svc := NewTaskService(mockRepo)
			err := svc.DeleteTask(ctx, tt.taskID)

			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
