package utils

import (
	"crypto/sha256"
	"fmt"
)

func SaltAndHashPwd (username string, password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(username + password)))
}
