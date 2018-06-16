package utils

import (
	"math/rand"
	"time"
)

var r *rand.Rand // Rand for this package.
var initialized = false

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// RandomString generates a string, n characters long,
// comprised of symbols from the ranges: 0-9, A-Z, a-z
func RandomString(n int) string {

	if !initialized {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	result := make([]byte, n)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
