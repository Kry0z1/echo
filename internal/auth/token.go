package auth

import (
	"errors"
	"time"

	"github.com/Kry0z1/echo/internal"
	db "github.com/Kry0z1/echo/internal/database"
	"github.com/golang-jwt/jwt"
)

func AuthUser(username, password string) (db.UserStored, bool) {
	user, err := db.GetUserByUsername(username)

	return user, err != nil && !internal.VerifyPassword(password, user.HashedPassword)
}

func CreateAuthToken(username string, exp time.Duration) (string, error) {
	if exp <= 0 {
		return "", errors.New("invalid expiration duration")
	}

	t := jwt.New(jwt.GetSigningMethod(algo))

	t.Claims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(exp).Unix(),
		Subject:   username,
	}

	return t.SignedString(secretKey)
}

func CheckAuthToken(token string) (db.UserStored, error) {
	credentialsError := errors.New("invalid credentials")

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return db.UserStored{}, err
	}

	claims := t.Claims.(jwt.MapClaims)

	username := claims["sub"].(string)
	if username == "" {
		return db.UserStored{}, credentialsError
	}

	user, err := db.GetUserByUsername(username)
	if err != nil {
		return db.UserStored{}, credentialsError
	}
	return user, nil
}
