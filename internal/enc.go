package internal

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

func GetPasswordHash(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(result), err
}
