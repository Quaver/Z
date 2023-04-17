package main

import "example.com/Quaver/Z/config"

func main() {
	err := config.Load()

	if err != nil {
		panic(err)
	}

	s := NewServer(config.Instance.Server.Port)
	s.Start()
}
