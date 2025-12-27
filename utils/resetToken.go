package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"tallyRead.com/config"
)

func GenerateResetToken() (string, string) {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	return token, hashedToken
}

func GetHashedToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateResetURL(token string) string {
	return config.Smtp.ClientURL + "/reset-password?token=" + token
}
