package main

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/webhooks"
)

func main() {
	err := config.Load("./config.json")

	if err != nil {
		panic(err)
	}

	db.InitializeSQL()
	db.InitializeRedis()
	webhooks.Initialize()
	chat.Initialize()
	multiplayer.InitializeBot()
	multiplayer.InitializeLobby()

	s := NewServer(config.Instance.Server.Port)
	s.Start()
}
