package utils

import (
	"fmt"
	"github.com/TwiN/go-away"
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

// BoolToEnabledString Converts a bool into an "enabled" or "disabled" string
func BoolToEnabledString(b bool) string {
	if b {
		return "enabled"
	}

	return "disabled"
}

// BoolToOnOffString Converts a bool into an "on" or "off" string
func BoolToOnOffString(b bool) string {
	if b {
		return "on"
	}

	return "off"
}

func CensorString(s string) string {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	return goaway.Censor(s)
}
