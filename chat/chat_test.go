package chat

import (
	"example.com/Quaver/Z/config"
	"testing"
)

func TestInitialize(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	Initialize()

	if len(channels) != len(config.Instance.ChatChannels) {
		t.Fatal("Expected more than zero initialized chat channels")
	}
}
