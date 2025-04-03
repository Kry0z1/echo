package handlers

import (
	"net/http"

	"github.com/Kry0z1/echo/internal/database"
	"github.com/labstack/echo/v4"
)

func CreateUser(ctx echo.Context) error {
	var user database.UserIn
	if err := ctx.Bind(&user); err != nil {
		ctx.String(http.StatusBadRequest, "Failed to parse json")
		return nil
	}

	if user.Username == "" {
		ctx.String(http.StatusBadRequest, "Missing username")
		return nil
	}

	if user.Password == "" {
		ctx.String(http.StatusBadRequest, "Missing password")
		return nil
	}

	if _, err := database.GetUserByUsername(user.Username); err == nil {
		ctx.String(http.StatusBadRequest, "User with such username already exists")
		return nil
	}

	userDB, err := database.CreateUserWithUsernameAndPassword(user)

	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, userDB.UserOut)

	return nil
}
