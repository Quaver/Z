package db

import (
	"example.com/Quaver/Z/config"
	"testing"
)

func TestGetUserBySteamId(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	_, err := GetUserBySteamId("76561198201861833")

	if err != nil {
		t.Fatal(err)
	}

	CloseSQLConnection()
}

func TestUpdateUserLatestActivity(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	err := UpdateUserLatestActivity(1)

	if err != nil {
		t.Fatal(err)
	}

	CloseSQLConnection()
}