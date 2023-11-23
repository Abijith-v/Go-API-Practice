package misc

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func GenerateString() (string, error) {

	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	password := make([]byte, 10)
	maxIndex := big.NewInt(int64(len(chars)))
	for i := 0; i < 10; i++ {
		index, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "", err
		}
		password[i] = chars[index.Int64()]
	}

	return string(password), nil
}

func RemoveSlashes(jsonStr string) string {
	var sb strings.Builder
	for _, char := range jsonStr {
		if char != '\\' {
			sb.WriteRune(char)
		}
	}
	return sb.String()
}
