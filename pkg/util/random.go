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

func randomNumber(n int) string {
	var randomNumber strings.Builder
	k := len(number)
	for i := 0; i < n; i++ {
		c := number[rand.Int63n(int64(k))]
		randomNumber.WriteByte(c)
	}
	return randomNumber.String()
}
func RandomMobile() string {
	return randomNumber(11)
}

func RandomAvatar() string {
	return randomNumber(20)
}
