package main

import (
	"testing"
)

func newTestTaskManager() *TaskManager {
	id := 0
	return &TaskManager{
		getNextID: func() int {
			id++
			return id
		},
	}
}

func TestAddTask(t *testing.T) {
	tm := newTestTaskManager()
	tm.AddTask("Test task 1")
	if len(tm.tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tm.tasks))
	}
	task := tm.tasks[0]
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	if task.Description != "Test task 1" {
		t.Errorf("expected description Test task 1, got '%s'", task.Description)
	}
	if task.Completed {
		t.Errorf("expected Completed false, got true")
	}
}

func TestCompleteTask(t *testing.T) {
	tm := newTestTaskManager()
	tm.AddTask("Test task 2")
	if tm.tasks[0].Completed {
		t.Errorf("expected task to be incomplete initially")
	}
	tm.CompleteTask(1)
	if !tm.tasks[0].Completed {
		t.Errorf("expected task to be marked as completed")
	}
}
