package main

import (
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
)

func main() {
	err := config.Load("./config.json")

	if err != nil {
		panic(err)
	}

	db.InitializeSQL()
	db.InitializeRedis()

	s := NewServer(config.Instance.Server.Port)
	s.Start()
}
