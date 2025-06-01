package util

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

const (
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*"
)

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandPassword(length int) string {
	return generateRandomPassword(length)
}

func generateRandomPassword(length int) string {
	// Create a character set by combining the desired characters
	characters := uppercaseLetters + lowercaseLetters + digits + specialChars

	// Create a byte slice to store the random bytes
	randomBytes := make([]byte, length)

	// Generate random bytes
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err) // Handle the error appropriately
	}

	// Create a password string by selecting random characters from the character set based on the random bytes
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		password[i] = characters[int(randomBytes[i])%len(characters)]
	}

	return string(password)
}
