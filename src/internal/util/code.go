package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateCode(n int) (string, error) {
	code := make([]byte, n)
	if _, err := rand.Read(code); err != nil {
		return "", err
	}
	return hex.EncodeToString(code), nil
}
