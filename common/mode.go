package common

import (
	"errors"
	"fmt"
	"strings"
)

type Mode int32

const (
	ModeKeys4 Mode = iota + 1
	ModeKeys7

	// ModeKeys1 New game modes so they start counting from 3

	ModeKeys1 = iota + 1
	ModeKeys2
	ModeKeys3
	ModeKeys5
	ModeKeys6
	ModeKeys8
	ModeKeys9
	ModeKeys10
	ModeEnumMaxValue
)

// GetModeString Returns a string version of game mode
func GetModeString(mode Mode) (string, error) {
	switch mode {
	case ModeKeys1:
		return "keys1", nil
	case ModeKeys2:
		return "keys2", nil
	case ModeKeys3:
		return "keys3", nil
	case ModeKeys4:
		return "keys4", nil
	case ModeKeys5:
		return "keys5", nil
	case ModeKeys6:
		return "keys6", nil
	case ModeKeys7:
		return "keys7", nil
	case ModeKeys8:
		return "keys8", nil
	case ModeKeys9:
		return "keys9", nil
	case ModeKeys10:
		return "keys10", nil
	default:
		return "", fmt.Errorf("%v is not a valid mode", mode)
	}
}

func GetShorthandGameModeString(mode Mode) string {
	switch mode {
	case ModeKeys1:
		return "1K"
	case ModeKeys2:
		return "2K"
	case ModeKeys3:
		return "3K"
	case ModeKeys4:
		return "4K"
	case ModeKeys5:
		return "5K"
	case ModeKeys6:
		return "6K"
	case ModeKeys7:
		return "7K"
	case ModeKeys8:
		return "8K"
	case ModeKeys9:
		return "9K"
	case ModeKeys10:
		return "10K"
	default:
		return "not_implemented"
	}
}

func GetModeFromShortHand(str string) (Mode, error) {
	switch strings.ToLower(str) {
	case "1k":
		return ModeKeys1, nil
	case "2k":
		return ModeKeys2, nil
	case "3k":
		return ModeKeys3, nil
	case "4k":
		return ModeKeys4, nil
	case "5k":
		return ModeKeys5, nil
	case "6k":
		return ModeKeys6, nil
	case "7k":
		return ModeKeys7, nil
	case "8k":
		return ModeKeys8, nil
	case "9k":
		return ModeKeys9, nil
	case "10k":
		return ModeKeys10, nil
	default:
		return -1, errors.New("game mode not valid")
	}
}
