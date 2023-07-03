package utils

import (
	"fmt"
	"math"
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

// TruncateString Truncates a string to a given max length
func TruncateString(str string, maxLength int) string {
	return str[:int(math.Min(float64(len(str)), float64(maxLength)))]
}
