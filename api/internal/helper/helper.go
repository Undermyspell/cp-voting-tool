package helper

import (
	"math/rand"
	"time"
)

func GetRandomId(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ12456789"

	randId := make([]byte, length)
	for i := 0; i < length; i++ {
		randId[i] = charset[r.Intn(len(charset))]
	}

	return string(randId)
}
