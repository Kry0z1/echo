package auth

import (
	"encoding/hex"
	"log"
	"os"
)

var (
	secretKey []byte
	algo      string
	err       error
)

func Init() {
	sk := os.Getenv("SECRET_KEY")

	secretKey, err = hex.DecodeString(sk)

	if err != nil {
		log.Fatal("failed to decode secret key: %w", err)
	}

	algo = os.Getenv("JWT_ALGO")
	if algo == "" {
		log.Fatal("missing JWT_ALGO in .env")
	}
}
