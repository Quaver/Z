package db

import (
	"example.com/Quaver/Z/config"
	"testing"
)

func TestGetUserBySteamId(t *testing.T) {
	_ = config.Load("../config.json")
	InitializeSQL()

	_, err := GetUserBySteamId("76561198201861833")

	if err != nil {
		t.Fatal(err)
	}

	CloseSQLConnection()
}
