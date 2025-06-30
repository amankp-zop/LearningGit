package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"assignment8/db"
	taskhandler "assignment8/handler/task"
	userhandler "assignment8/handler/user"
	taskrepository "assignment8/repository/task"
	userrepository "assignment8/repository/user"
	taskservice "assignment8/service/task"
	userservice "assignment8/service/user"
)

func main() {
	database, err := db.ConnectDB()

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	userRepo := userrepository.NewUserRepository(database)
	userSvc := userservice.NewUserService(userRepo)
	userHlr := userhandler.NewUserHandler(userSvc)

	taskRepo := taskrepository.NewTaskRepository(database)
	taskService := taskservice.NewTaskService(taskRepo)
	taskHandler := taskhandler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", userHlr.GetUserHandler)
	mux.HandleFunc("POST /users", userHlr.CreateUserHandler)

	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetUserTasks)
	mux.HandleFunc("PUT /tasks", taskHandler.UpdateTask)
	mux.HandleFunc("PUT /tasks/{id}/complete", taskHandler.MarkTaskComplete)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	fmt.Println("Server running at http://localhost:8080")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
