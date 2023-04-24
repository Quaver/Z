package db

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/config"
	"testing"
)

func TestGetUserStats(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	statsKeys4, err := GetUserStats(1, common.ModeKeys4)

	if err != nil {
		t.Fatal(err)
	}

	if statsKeys4.TotalScore == 0 {
		t.Fatalf("expected a non-zero total score value for keys4")
	}

	statsKeys7, err := GetUserStats(1, common.ModeKeys7)

	if err != nil {
		t.Fatal(err)
	}

	if statsKeys7.TotalScore == 0 {
		t.Fatalf("expected a non-zero total score value for keys7")
	}

	CloseSQLConnection()
}
