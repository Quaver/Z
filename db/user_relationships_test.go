package db

import (
	"example.com/Quaver/Z/config"
	"testing"
)

func TestUserFriendsList(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	friends, err := GetUserFriendsList(1)

	if err != nil {
		t.Fatal(err)
	}

	if len(friends) == 0 {
		t.Fatal("expected more than zero friends")
	}

	CloseSQLConnection()
}
