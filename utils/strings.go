package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var randomSeeded = false

// GenerateRandomString Generates a random string of a given length
func GenerateRandomString(length int) string {
	if !randomSeeded {
		rand.Seed(time.Now().UnixNano())
		randomSeeded = true
	}

	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
