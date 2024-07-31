package main

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/handlers"
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/webhooks"
	"flag"
)

func main() {
	configPath := flag.String("config", "../../config.json", "path to config file")
	flag.Parse()

	if err := config.Load(*configPath); err != nil {
		panic(err)
	}

	db.InitializeSQL()
	db.InitializeRedis()
	handlers.AddRedisHandlers()
	webhooks.Initialize()
	chat.Initialize()
	multiplayer.InitializeChatBot()
	multiplayer.InitializeLobby()

	s := NewServer(config.Instance.Server.Port)
	s.Start()
}
