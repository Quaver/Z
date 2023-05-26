package chat

import (
	"example.com/Quaver/Z/common"
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

func TestGetAvailableChannels(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	Initialize()

	chansIncludingAdmin := GetAvailableChannels(common.UserGroupDeveloper | common.UserGroupNormal)

	if len(chansIncludingAdmin) <= 1 {
		t.Fatal("expected at least one admin channel")
	}

	chansWithoutAdmin := GetAvailableChannels(common.UserGroupNormal)

	if len(chansWithoutAdmin) != 1 {
		t.Fatal("expected only one channel without admin")
	}
}
