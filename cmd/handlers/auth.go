package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Kry0z1/echo/internal/auth"
	"github.com/Kry0z1/echo/internal/database"
	"github.com/labstack/echo/v4"
)

type Token struct {
	AccessToken string
	TokenType   string
}

type TokenIn struct {
	token string
}

func ContextUser(ctx echo.Context) *database.UserOut {
	user := ctx.Get("contextUser")
	if user == nil {
		return nil
	}
	return user.(*database.UserOut)
}

func LoginForAccessToken(ctx echo.Context) error {
	var user database.UserIn

	if err := ctx.Bind(&user); err != nil {
		ctx.String(http.StatusBadRequest, "Failed to parse json")
		return nil
	}

	if _, ok := auth.AuthUser(user.Username, user.Password); !ok {
		ctx.String(http.StatusUnauthorized, "Wrong username of password")
		return nil
	}

	token, err := auth.CreateAuthToken(user.Username, time.Hour*72)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, Token{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

func CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("WWW-Authenticate", "Bearer")

		authHeaders := ctx.Request().Header["Authorization"]
		if len(authHeaders) == 0 {
			ctx.String(http.StatusUnauthorized, "Missing Authorization header")
			return nil
		}

		var token string

		for _, header := range authHeaders {
			splitted := strings.Split(header, " ")
			if len(splitted) != 2 || splitted[0] != "Bearer" {
				continue
			}
			token = splitted[1]
			break
		}

		if token == "" {
			ctx.String(http.StatusUnauthorized, "Missing Authorization header")
			return nil
		}

		user, err := auth.CheckAuthToken(token)
		if errors.Is(err, auth.ErrInvalidToken) {
			ctx.String(http.StatusUnauthorized, "Invalid token")
			return nil
		}
		if err != nil {
			return err
		}

		ctx.Set("contextUser", &user.UserOut)
		return next(ctx)
	}
}
