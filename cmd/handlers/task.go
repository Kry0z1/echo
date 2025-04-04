package handlers

import (
	"net/http"
	"slices"
	"strconv"
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

	result := make([]database.TaskOut, len(tasks))
	for i, task := range tasks {
		result[i] = task.TaskOut
	}

	return ctx.JSON(http.StatusOK, result)
}

func UpdateTask(ctx echo.Context) error {
	var task database.TaskStored
	user := ContextUser(ctx)

	err := ctx.Bind(&task)

	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid format")
		return nil
	}

	taskStored, err := database.GetTask(task.ID)

	if err != nil {
		ctx.String(http.StatusNotFound, "Task with such id does not exists")
		return nil
	}

	if taskStored.UserID != task.UserID {
		ctx.String(http.StatusUnauthorized, "Cannot update creator id")
		return nil
	}

	if taskStored.UserID != user.ID {
		ctx.String(http.StatusUnauthorized, "Cannot update tasks for other users")
		return nil
	}

	err = database.UpdateTask(&task)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, task.TaskOut)
}

func RemoveTask(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.QueryParam("id"))
	user := ContextUser(ctx)

	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad id")
		return nil
	}

	if id == 0 {
		ctx.String(http.StatusBadRequest, "Missing id")
		return nil
	}

	task, err := database.GetTask(id)

	if err != nil {
		ctx.String(http.StatusNotFound, "Task is not found")
		return nil
	}

	if task.UserID != user.ID {
		ctx.String(http.StatusUnauthorized, "Cannot remove tasks if other users")
		return nil
	}

	err = database.RemoveTask(id)
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, "")
}
