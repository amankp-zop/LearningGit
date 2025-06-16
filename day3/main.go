package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Task represents a single to-do item.
type Task struct {
	ID          int    // Unique identifier for the task
	Description string // Description of the task
	Completed   bool   // Status of the task (completed or not)
}

// TaskManager manages a list of tasks and ID generation.
type TaskManager struct {
	tasks     []Task     // Slice to store tasks
	getNextID func() int // Function to generate unique IDs
}

// idGenerator returns a closure that generates incrementing IDs.
func idGenerator() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

// AddTask adds a new task with the given description.
func (tm *TaskManager) AddTask(description string) {
	id := tm.getNextID()
	task := Task{
		ID:          id,
		Description: description,
		Completed:   false,
	}
	tm.tasks = append(tm.tasks, task)
	fmt.Printf("Task Added: %d - %s\n", task.ID, task.Description)
}

// ListTasks prints all pending (not completed) tasks.
func (tm *TaskManager) ListTasks() {
	fmt.Println("\nPending Tasks:")
	found := false
	for _, task := range tm.tasks {
		if !task.Completed {
			fmt.Printf("%d: %s\n", task.ID, task.Description)
			found = true
		}
	}
	if !found {
		fmt.Println("No pending tasks.")
	}
}

// CompleteTask marks the task with the given ID as completed.
func (tm *TaskManager) CompleteTask(id int) {
	for i, task := range tm.tasks {
		if task.ID == id {
			if task.Completed {
				fmt.Printf("Task %d is already completed.\n", id)
				return
			}
			tm.tasks[i].Completed = true
			fmt.Printf("Marked task %d as completed.\n", id)
			return
		}
	}
	fmt.Printf("Task with ID %d not found.\n", id)
}

// DeleteTask removes the task with the given ID from the list.
func (tm *TaskManager) DeleteTask(id int) {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			fmt.Printf("Deleted task %d: %s\n", id, task.Description)
			return
		}
	}
	fmt.Printf("Task with ID %d not found.\n", id)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	manager := TaskManager{
		getNextID: idGenerator(),
	}

	// Main loop for the menu-driven interface
	for {
		fmt.Println("\n📝💪🏻 Task Tracker Menu")
		fmt.Println("1. Add Task")
		fmt.Println("2. Mark Task as Complete")
		fmt.Println("3. Delete Task")
		fmt.Println("4. List Pending Tasks")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			// Add a new task
			fmt.Print("Enter task description: ")
			desc, _ := reader.ReadString('\n')
			manager.AddTask(strings.TrimSpace(desc))

		case "2":
			// Mark a task as complete
			fmt.Print("Enter task ID to mark as complete: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("Invalid ID.")
				continue
			}
			manager.CompleteTask(id)

		case "3":
			// Delete a task
			fmt.Print("Enter task ID to delete: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("Invalid ID.")
				continue
			}
			manager.DeleteTask(id)

		case "4":
			// List all pending tasks
			manager.ListTasks()

		case "5":
			// Exit the program
			fmt.Println("Exiting. Goodbye! 👋")
			return

		default:
			// Handle invalid menu options
			fmt.Println("Invalid option. Please choose 1--5.")
		}
	}
}
