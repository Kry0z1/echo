package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Kry0z1/echo/cmd/handlers"
	"github.com/Kry0z1/echo/internal/auth"
	"github.com/Kry0z1/echo/internal/database"
)

func init() {
	godotenv.Load()

	database.Init()
	auth.Init()
}

func main() {
	e := echo.New()
	e.Debug = true

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogError:  true,
		// Skipper: func(c echo.Context) bool {
		// 	return c.Response().Status < 500
		// },
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				fmt.Printf("%v %v %v\n", v.Status, v.URI, v.Latency.Abs().String())
				return nil
			}
			fmt.Printf("%v %v %v %v\n", v.Status, v.URI, v.Latency.Abs().String(), v.Error)
			return nil
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}, handlers.CheckToken)
	e.POST("/login", handlers.LoginForAccessToken)
	e.POST("/signup", handlers.CreateUser)

	e.POST("/task/create", handlers.CreateTask, handlers.CheckToken)
	e.POST("/task/update", handlers.UpdateTask, handlers.CheckToken)
	e.GET("/tasks", handlers.GetTasksForUser, handlers.CheckToken)

	e.Logger.Fatal(e.Start(":21000"))
}
