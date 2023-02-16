package helper

import (
	"math/rand"
	"time"
)

func GetRandomId() string {
	length := 15
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyz12456789"

	randId := make([]byte, length)
	for i := 0; i < length; i++ {
		randId[i] = charset[rand.Intn(len(charset))]
	}

	return string(randId)
}
