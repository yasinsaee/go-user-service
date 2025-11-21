package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	UpperCaseLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerCaseLetter = "abcdefghijklmnopqrstuvwxyz"
	Number          = "1234567890"
	Characters      = "!@$%^&*-+="
)

func Rand(length int, available ...string) (res string) {
	randIN := ""
	if len(available) < 1 {
		randIN = Number
	}
	for _, c := range available {
		randIN += c
	}
	randList := strings.Split(randIN, "")
	rand.Seed(time.Now().UTC().Unix())
	for i := 0; i < length; i++ {
		index := rand.Intn(len(randList) - 1)
		res += randList[index]
	}
	return res
}
func RandSeedUnixNano(length int, available ...string) (res string) {
	randIN := ""
	if len(available) < 1 {
		randIN = Number
	}
	for _, c := range available {
		randIN += c
	}
	randList := strings.Split(randIN, "")

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < length; i++ {
		index := rand.Intn(len(randList) - 1)
		res += randList[index]
	}
	return res
}

func RandSeedUnixNanoInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min+1) + min
}
