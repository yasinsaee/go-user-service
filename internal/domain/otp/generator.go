package otp

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateNumericCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	format := fmt.Sprintf("%%0%dd", n)
	return fmt.Sprintf(format, rand.Intn(pow10(n)))
}

func pow10(n int) int {
	p := 1
	for i := 0; i < n; i++ {
		p *= 10
	}
	return p
}
