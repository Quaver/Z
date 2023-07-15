package db

import (
	"example.com/Quaver/Z/config"
	"fmt"
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

func TestGetUserRelationship(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	relationship, err := GetUserRelationship(1, 2)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(relationship)

	CloseSQLConnection()
}

func TestAddFriend(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	err := AddFriend(1, 2)

	if err != nil {
		t.Fatal(err)
	}

	err = RemoveFriend(1, 2)

	if err != nil {
		t.Fatal(err)
	}

	CloseSQLConnection()
}

