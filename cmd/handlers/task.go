package handlers

import (
	"net/http"

	"github.com/Kry0z1/echo/internal/database"
	"github.com/labstack/echo/v4"
)

func CreateTask(ctx echo.Context) error {
	var task database.TaskIn
	user := ContextUser(ctx)

	err := ctx.Bind(&task)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid format")
		return nil
	}

	if task.UserID == 0 {
		task.UserID = user.ID
	}

	if task.Title == "" {
		ctx.String(http.StatusBadRequest, "Title cannot be empty")
		return nil
	}

	if task.Priority < 0 || task.Priority > 3 {
		ctx.String(http.StatusBadRequest, "Priority must be an integer between 0 and 3")
		return nil
	}

	if user == nil {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return nil
	}

	if user.ID != task.UserID {
		ctx.String(http.StatusUnauthorized, "Cannot create tasks for other users")
		return nil
	}

	taskDB, err := database.CreateTask(task)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, taskDB.TaskOut)
}
