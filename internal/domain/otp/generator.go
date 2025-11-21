package otp

import (
	"math/rand"
	"time"
)

func GenerateCode(length int, charset string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
