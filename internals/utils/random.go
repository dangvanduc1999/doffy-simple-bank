package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphaNumeric)
	for i := 0; i < length; i++ {
		sb.WriteByte(alphaNumeric[rand.Intn(k)])
	}
	return sb.String()
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR"}
	return currencies[rand.Intn(len(currencies))]
}
