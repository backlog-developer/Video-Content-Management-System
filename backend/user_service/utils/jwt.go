package utils

import (
	"os"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))
