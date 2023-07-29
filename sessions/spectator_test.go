package sessions

import (
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"testing"
)

func TestAddSpectator(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	db.InitializeRedis()

	user1 := NewUser(nil, &db.User{Id: 1, SteamId: "1", Username: "User #1"})
	user2 := NewUser(nil, &db.User{Id: 2, SteamId: "2", Username: "User #2"})

	user1.AddSpectator(user2)

	if len(user1.GetSpectators()) != 1 {
		t.Fatal("Expected user 1 to have 1 spectator")
	}

	if len(user2.GetSpectating()) != 1 {
		t.Fatal("Expected user 2 to be spectating 1 person")
	}
}

func TestRemoveSpectator(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	db.InitializeRedis()

	user1 := NewUser(nil, &db.User{Id: 1, SteamId: "1", Username: "User #1"})
	user2 := NewUser(nil, &db.User{Id: 2, SteamId: "2", Username: "User #2"})

	user1.AddSpectator(user2)
	user1.RemoveSpectator(user2)

	if len(user1.GetSpectators()) != 0 {
		t.Fatal("Expected user 1 to have 0 spectators")
	}
}

func TestStopSpectatingAll(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	db.InitializeRedis()

	user1 := NewUser(nil, &db.User{Id: 1, SteamId: "1", Username: "User #1"})
	user2 := NewUser(nil, &db.User{Id: 2, SteamId: "2", Username: "User #2"})

	user1.AddSpectator(user2)
	user2.StopSpectatingAll()

	if len(user1.GetSpectators()) != 0 {
		t.Fatal("expected user 1 to have 0 spectators")
	}

	if len(user2.GetSpectating()) != 0 {
		t.Fatal("expected user 2 to be spectating 0 people")
	}
}
