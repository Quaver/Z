package db

import (
	"example.com/Quaver/Z/config"
	"testing"
)

func TestInsertLoginIpAddress(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	err := InsertLoginIpAddress(1, "192.168.1.1")

	if err != nil {
		t.Fatal(err)
	}

	CloseSQLConnection()
}
