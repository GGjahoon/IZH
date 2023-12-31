package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	number   = "1234567890"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Int63n(int64(k))]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomUser() string {
	return RandomString(6)
}

func RandomNumber(n int) string {
	var randomNumber strings.Builder
	k := len(number)
	for i := 0; i < n; i++ {
		c := number[rand.Int63n(int64(k))]
		randomNumber.WriteByte(c)
	}
	return randomNumber.String()
}
func RandomMobile() string {
	return RandomNumber(11)
}

func RandomAvatar() string {
	return RandomNumber(20)
}
func RandoInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
