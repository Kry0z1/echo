package handlers

import (
	"net/http"
	"slices"
	"strings"

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

func GetTasksForUser(ctx echo.Context) error {
	sortBy := ctx.QueryParam("sort")
	done := ctx.QueryParam("done")
	desc := ctx.QueryParam("desc")

	user := ContextUser(ctx)
	var tasks []database.TaskStored
	var err error

	if done == "true" {
		tasks, err = database.GetDoneTasksForUser(user.ID)
	} else {
		tasks, err = database.GetTasksForUser(user.ID)
	}

	if err != nil {
		return err
	}

	slices.SortFunc(tasks, func(a database.TaskStored, b database.TaskStored) int {
		var res int
		switch sortBy {
		case "priority":
			res = a.Priority - b.Priority
		case "due_to":
			res = a.DueTo - b.DueTo
		case "starts_at":
			res = a.DueTo - b.DueTo
		case "title":
			res = strings.Compare(a.Title, b.Title)
		default:
			res = a.ID - b.ID
		}
		if desc == "true" {
			res *= -1
		}
		return res
	})

	return ctx.JSON(http.StatusOK, tasks)
}
