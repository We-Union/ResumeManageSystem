package utils

import (
	"math/rand"
	"time"
)

func GetToken(length int) (token string) {
	rand.Seed(time.Now().UnixNano())
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}

	return string(result)
}
