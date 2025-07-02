// @title           Task Manager API
// @version         1.0
// @description     This is a sample server for managing tasks.
// @termsOfService  http://example.com/terms/

// @contact.name   Aman
// @contact.email  aman@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @schemes http
package main

import (
	"fmt"
	"net/http"

	taskhandler "assignment8/handler/task"
	userhandler "assignment8/handler/user"
	"assignment8/migrations"
	taskrepository "assignment8/repository/task"
	userrepository "assignment8/repository/user"
	taskservice "assignment8/service/task"
	userservice "assignment8/service/user"

	_ "assignment8/docs" // ✨ Import generated docs (adjust if needed)

	httpSwagger "github.com/swaggo/http-swagger"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	userRepo := userrepository.NewUserRepository()
	userSvc := userservice.NewUserService(userRepo)
	userHlr := userhandler.NewUserHandler(userSvc)

	taskRepo := taskrepository.NewTaskRepository()
	taskService := taskservice.NewTaskService(taskRepo)
	taskHandler := taskhandler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	app.GET("/users/{id}", userHlr.GetUserHandler)
	app.POST("/users", userHlr.CreateUserHandler)
	app.POST("/tasks", taskHandler.CreateTask)
	app.GET("/tasks/{id}", taskHandler.GetUserTasks)
	app.PUT("/tasks", taskHandler.UpdateTask)
	app.PUT("/tasks/{id}/complete", taskHandler.MarkTaskComplete)
	app.DELETE("/tasks/{id}", taskHandler.DeleteTask)

	fmt.Println("Server running at http://localhost:8080")
	app.Run()
}
