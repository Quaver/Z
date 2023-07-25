package spectator

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/sessions"
	"testing"
)

func TestAddSpectator(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	db.InitializeRedis()
	chat.Initialize()

	user1 := sessions.NewUser(nil, &db.User{Id: 1, SteamId: "1", Username: "User #1"})
	InitializeUser(user1)

	user2 := sessions.NewUser(nil, &db.User{Id: 2, SteamId: "2", Username: "User #2"})
	InitializeUser(user2)

	GetUser(user1).AddSpectator(GetUser(user2))

	if len(GetUser(user1).GetSpectators()) != 1 {
		t.Fatal("Expected user 1 to have 1 spectator")
	}

	if len(GetUser(user2).GetSpectating()) != 1 {
		t.Fatal("Expected user 2 to be spectating 1 person")
	}

	channel := chat.GetSpectatorChannel(user1.Info.Id)

	if channel == nil {
		t.Fatal("Spectator chat channel is nil")
	}

	if len(channel.Participants) != 2 {
		t.Fatal("Expected 2 participants in the channel")
	}

	UninitializeUser(user1)
	UninitializeUser(user2)
}

func TestRemoveSpectator(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	db.InitializeRedis()
	chat.Initialize()

	user1 := sessions.NewUser(nil, &db.User{Id: 1, SteamId: "1", Username: "User #1"})
	InitializeUser(user1)

	user2 := sessions.NewUser(nil, &db.User{Id: 2, SteamId: "2", Username: "User #2"})
	InitializeUser(user2)

	GetUser(user1).AddSpectator(GetUser(user2))
	GetUser(user1).RemoveSpectator(GetUser(user2))

	channel := chat.GetSpectatorChannel(user1.Info.Id)

	if channel != nil {
		t.Fatal("Spectator chat channel is still active")
	}

	UninitializeUser(user1)
	UninitializeUser(user2)
}

func TestStopSpectatingAll(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	db.InitializeRedis()
	chat.Initialize()

	user1 := sessions.NewUser(nil, &db.User{Id: 1, SteamId: "1", Username: "User #1"})
	InitializeUser(user1)

	user2 := sessions.NewUser(nil, &db.User{Id: 2, SteamId: "2", Username: "User #2"})
	InitializeUser(user2)

	GetUser(user1).AddSpectator(GetUser(user2))
	GetUser(user2).StopSpectatingAll()

	channel := chat.GetSpectatorChannel(user1.Info.Id)

	if channel != nil {
		t.Fatal("Spectator chat channel is still active")
	}

	UninitializeUser(user1)
	UninitializeUser(user2)
}
