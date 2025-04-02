package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

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
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
