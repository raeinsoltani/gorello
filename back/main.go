package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/raeinsoltani/gorello/back/db"
	"github.com/raeinsoltani/gorello/back/handlers"
	customMiddleware "github.com/raeinsoltani/gorello/back/middleware"
	"github.com/raeinsoltani/gorello/back/repository/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	time.Sleep(3 * time.Second)
	db.Init()
	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	c := jaegertracing.New(e, nil)
	defer func(c io.Closer) {
		err := c.Close()
		if err != nil {
			log.Printf("Failed to close tracer: %v\n", err)
		}
	}(c)

	e.Use(middleware.Logger())
	// e.Use(middleware.CSRF())
	e.Use(middleware.CORS())
	e.Use(echoprometheus.NewMiddleware("gorello"))

	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/", func(c echo.Context) error {
		log.Println("Handling the request for the root route")
		return c.String(http.StatusOK, "Hello, World!")
	})

	userRepo := gorm.NewUserRepo(db.DB)
	workspaceRepo := gorm.NewWorkspaceRepo(db.DB)
	userWorkspaceRoleRepo := gorm.NewUserWorkspaceRoleRepo(db.DB)
	taskRepo := gorm.NewTaskRepo(db.DB)

	userHandler := handlers.NewUserHandler(userRepo)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceRepo, userWorkspaceRoleRepo, userRepo)

	e.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	auth := e.Group("/auth")
	users := e.Group("/users")
	workspaces := e.Group("/workspaces")
	tasks := e.Group("/workspaces/:workspaceId/tasks")

	// User auth Handlers
	auth.POST("/signup", userHandler.Register)
	auth.POST("/login", userHandler.Login)

	// Users Handlers
	users.Use(customMiddleware.JWTAuthentication)
	users.GET("/", userHandler.GetUsers)
	users.GET("/:username", userHandler.GetUser)
	users.PUT("/:username", userHandler.UpdateUser)
	users.DELETE("/:username", userHandler.DeleteUser)
	users.GET("/search", userHandler.SearchUsers)

	// Workspaces Handlers
	workspaces.Use(customMiddleware.JWTAuthentication)
	workspaces.GET("/", workspaceHandler.GetWorkspaces)
	workspaces.POST("/", workspaceHandler.CreateWorkspace)
	workspaces.GET("/:workspaceId", workspaceHandler.GetWorkspaceDescription)
	workspaces.PUT("/:workspaceId", workspaceHandler.UpdateWorkspace)
	workspaces.DELETE("/:workspaceId", workspaceHandler.DeleteWorkspace)

	// Task Handlers
	tasks.Use(customMiddleware.JWTAuthentication)
	tasks.GET("/", handlers.NewTaskHandler(taskRepo).GetTasks)
	tasks.POST("/", handlers.NewTaskHandler(taskRepo).CreateTask)
	tasks.GET("/:taskId", handlers.NewTaskHandler(taskRepo).GetTask)
	tasks.PUT("/:taskId", handlers.NewTaskHandler(taskRepo).UpdateTask)
	tasks.DELETE("/:taskId", handlers.NewTaskHandler(taskRepo).DeleteTask)

	log.Println("Starting Echo server on port 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
